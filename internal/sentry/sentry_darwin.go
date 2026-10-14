package sentry

import (
	"context"
	"fmt"
)

type Sentry struct {
	Path   string
	Inode  uint64
	Offset int64
}

func NewSentry(path string) (*Sentry, error) {
	return nil, fmt.Errorf("inotify sentry is only supported on Linux kernel (please run inside Docker container)")
}

func (s *Sentry) Close() error {
	return fmt.Errorf("file info only supported on Linux kernel (please run inside Docker container)")
}

func (s *Sentry) Watch(ctx context.Context) (<-chan LogEvent, <-chan error) {
	eventsChan := make(chan LogEvent)
	errChan := make(chan error)
	close(eventsChan)
	close(errChan)
	return eventsChan, errChan
}
