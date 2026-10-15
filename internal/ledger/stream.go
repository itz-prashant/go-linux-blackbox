package ledger

import (
	"sync"
	"time"
)

type StreamMessage struct {
	ID          uint64
	Payload       string
	DispatchedAt time.Time
	Acked        bool
}

type Stream struct {
	mu       sync.RWMutex
	messages map[uint64]*StreamMessage
	lastId   uint64
}

func NewStream() *Stream {
	return &Stream{
		messages: make(map[uint64]*StreamMessage),
		lastId:   0,
	}
}

func (s *Stream) Add(payload string) *StreamMessage {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastId++
	msg := &StreamMessage{
		ID:           s.lastId,
		Payload:       payload,
		DispatchedAt: time.Now(),
		Acked:        false,
	}
	s.messages[s.lastId] = msg
	return msg
}

func (s *Stream) Ack(id uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg, exists := s.messages[id]

	if !exists {
		return false
	}
	msg.Acked = true
	return true
}

func (s *Stream) GetPending(timeout time.Duration) []*StreamMessage {
	s.mu.Lock()
	defer s.mu.Unlock()

	var orphaned []*StreamMessage
	now := time.Now()

	for _, msg := range s.messages {
		if !msg.Acked && now.Sub(msg.DispatchedAt) > timeout {
			msg.DispatchedAt = now
			orphaned = append(orphaned, msg)
		}
	}
	return orphaned
}

func (s *Stream) PendingCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := 0
	for _, msg := range s.messages {
		if !msg.Acked {
			count++
		}
	}
	return count
}
