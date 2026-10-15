package tape

import "time"

// Record tape recorder ke ek slot ka data hai
type Record struct {
	Timestamp time.Time // Kis waqt record hua
	Payload   string    // Error message / line
	Count     uint32    // Kitni baar repeat hua
}

// RingBuffer ek fixed-size circular memory buffer hai jo kabhi RAM nahi badhata
type RingBuffer struct {
	slots    []Record // Fixed array
	capacity int      // Kul slots (e.g. 100)
	head     int      // Current write cursor (0 to capacity-1)
	count    int      // Abhi kitne slots filled hain
}

// NewRingBuffer ek naya fixed-capacity ring buffer banata hai
func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		slots:    make([]Record, capacity),
		capacity: capacity,
		head:     0,
		count:    0,
	}
}

// Push naye record ko tape par likhta hai aur purane data ko overwrite karta hai
func (r *RingBuffer) Push(rec Record) {
	// 1. Current slot par data likho:
	r.slots[r.head] = rec

	// 2. Head cursor ko aage badhao (Circular wrap-around):
	r.head = (r.head + 1) % r.capacity

	// 3. Filled count update karo (Max = capacity):
	if r.count < r.capacity {
		r.count++
	}
}

// GetAll tape ke saare records ko oldest se newest order mein return karta hai
func (r *RingBuffer) GetAll() []Record {
	result := make([]Record, r.count)

	// Agar buffer full hai, toh sabse purana data r.head par hota hai
	start := 0
	if r.count == r.capacity {
		start = r.head
	}

	for i := 0; i < r.count; i++ {
		index := (start + i) % r.capacity
		result[i] = r.slots[index]
	}

	return result
}

// Size current bhare hue records ka count return karta hai
func (r *RingBuffer) Size() int {
	return r.count
}