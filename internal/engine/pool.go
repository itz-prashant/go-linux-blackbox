package engine

import (
	"bytes"
	"sync"
)

type BufferPool struct {
	pool sync.Pool
}

func NewBufferPool () *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func()any {
				// Agar pool khali ho, tabhi yeh function chalta hai
				// Pre-allocate: 1KB (1024 bytes) capacity ka buffer banata hai
				return bytes.NewBuffer(make([]byte,0,1024))
			},
		},
	}
}

// Get: Pool se ek saaf recycled buffer nikaalta hai (Zero Allocation)
func (p *BufferPool) Get() *bytes.Buffer{
	return p.pool.Get().(*bytes.Buffer)
}

// Put: Kaam khatam hone ke baad buffer ko dho kar (Reset) wapas pool mein daalta hai

func (p *BufferPool) Put(buf *bytes.Buffer){
	buf.Reset() // SABSE CRITICAL: Cursor ko 0 par laata hai taaki data leak na ho!
	p.pool.Put(buf)
}