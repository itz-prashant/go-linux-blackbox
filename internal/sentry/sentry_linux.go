package sentry

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"
	"unsafe"
)

type Sentry struct {
	Path      string
	Inode     uint64
	Offset    int64
	InotifyFd int
	watchFd   int
}

func NewSentry(path string) (*Sentry, error) {
	inode, err := GetInode(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get inode: %w", err)
	}

	fd, err := syscall.InotifyInit()
	if err != nil {
		return nil, fmt.Errorf("failed to init inotify: %w", err)
	}

	flags := uint32(syscall.IN_MODIFY | syscall.IN_MOVE_SELF | syscall.IN_DELETE_SELF)

	watchFd, err := syscall.InotifyAddWatch(fd, path, flags)
	if err != nil {
		syscall.Close(fd)
		return nil, fmt.Errorf("failed to add watch: %w", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		syscall.Close(fd)
		return nil, err
	}

	// Pehla offset file ke aakhiri byte par set kiya (sirf naye logs sunenge)
	offset := info.Size()

	return &Sentry{
		Path:      path,
		Inode:     inode,
		Offset:    offset,
		InotifyFd: fd,
		watchFd:   watchFd,
	}, nil
}

func (s *Sentry) Close() error {
	return syscall.Close(s.InotifyFd)
}

// readNewLine file ke offset se lekar nayi lines read karta hai aur offset update karta hai
func (s *Sentry) readNewLine() ([]string, error) {
	file, err := os.Open(s.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Cursor ko purane offset par le gaye:
	_, err = file.Seek(s.Offset, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("failed to seek: %w", err)
	}

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Drift-proof fix: Naya offset seedha file ke size se liya!
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}
	s.Offset = info.Size()

	return lines, nil
}

// Watch background mein Linux Kernel ke inotify events sunta hai aur channel par events bhejta hai
func (s *Sentry) Watch(ctx context.Context) (<-chan LogEvent, <-chan error) {
	eventsChan := make(chan LogEvent, 100)
	errChan := make(chan error, 1)

	// Background Worker (Goroutine) shuru kiya:
	go func() {
		defer close(eventsChan)
		defer close(errChan)

		// 4KB ka raw byte buffer jisme Linux Kernel C-struct likhega
		var buf [4096]byte

		for {
			// Check: Agar program band hone wala hai toh loop exit kar do
			select {
			case <-ctx.Done():
				return
			default:
			}

			// ASLI LOW-LEVEL MAGIC:
			// Yeh line program ko SLEEP kara deti hai (0% CPU).
			// Kernel tabhi jagata hai jab file mein kuch likha jaye!
			n, err := syscall.Read(s.InotifyFd, buf[:])
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					errChan <- fmt.Errorf("inotify read error: %w", err)
					return
				}
			}

			// Agar read bytes event ke size se chhota hai toh ignore karo
			if n < syscall.SizeofInotifyEvent {
				continue
			}

			// C-Struct ko Go struct mein cast kiya bina memory copy kiye (Zero-Copy):
			rawEvent := (*syscall.InotifyEvent)(unsafe.Pointer(&buf[0]))
			mask := rawEvent.Mask

			// CASE 1: File Modify Hui (Naya Error Likh Gaya!)
			if mask&syscall.IN_MODIFY != 0 {
				lines, err := s.readNewLine()
				if err != nil {
					errChan <- err
					continue
				}

				// Har nayi line ko channel pipe mein daal do:
				for _, line := range lines {
					eventsChan <- LogEvent{
						Line:      line,
						Inode:     s.Inode,
						Timestamp: time.Now().UTC(),
					}
				}
			}

			// CASE 2: Logrotate Hua (File Rename ya Delete hui!)
			if mask&(syscall.IN_MOVE_SELF|syscall.IN_DELETE_SELF) != 0 {
				fmt.Printf("[SENTRY] Logrotate detected! File on Inode #%d was moved/deleted.\n", s.Inode)
			}
		}
	}()

	return eventsChan, errChan
}
