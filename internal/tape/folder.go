package tape

import (
	"sync"
	"time"
)

// EventFolder identical duplicate errors ko 5-second window mein
// pichka kar (fold karke) sirf 1 entry bana deta hai.
type EventFolder struct {
	mu         sync.Mutex    // 4 workers ke concurrent access ke liye lock
	ring       *RingBuffer   // Asli Circular Ring Buffer
	window     time.Duration // Duplicate folding window (e.g. 5 seconds)
	lastHash   uint64        // Pichle error ka fingerprint
	lastRecord Record        // Current fold ho raha record
	hasActive  bool          // Kya abhi koi record pending hai?
}

// NewEventFolder ek naya folding engine banata hai
func NewEventFolder(ring *RingBuffer, window time.Duration) *EventFolder {
	return &EventFolder{
		ring:   ring,
		window: window,
	}
}

// Feed naye error log ko evaluate karta hai: fold kare ya tape par push kare
func (f *EventFolder) Feed(payload string, now time.Time) {
	// 1. Log ka fast 5-nanosecond fingerprint nikalo:
	hash := HashString(payload)

	f.mu.Lock()
	defer f.mu.Unlock()

	// 2. Agar koi record pehle se fold ho raha hai:
	if f.hasActive {
		// Check: Kya same error hai aur 5s window ke andar hai?
		if f.lastHash == hash && now.Sub(f.lastRecord.Timestamp) <= f.window {
			f.lastRecord.Count++ // Duplicate found! Sirf count badhao!
			return
		}

		// Naya error aa gaya ya 5s beet gaye -> Purana folded record Tape mein push kar do:
		f.ring.Push(f.lastRecord)
	}

	// 3. Naye error ko current active record banao:
	f.lastHash = hash
	f.lastRecord = Record{
		Timestamp: now,
		Payload:   payload,
		Count:     1,
	}
	f.hasActive = true
}

// Flush program band hone par pending folded record ko safely Tape mein push karta hai
func (f *EventFolder) Flush() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.hasActive {
		f.ring.Push(f.lastRecord)
		f.hasActive = false
	}
}
