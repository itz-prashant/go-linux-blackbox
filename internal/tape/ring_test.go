package tape

import (
	"fmt"
	"testing"
	"time"
)

func TestRingBufferOverWrite(t *testing.T) {
	rb := NewRingBuffer(3)

	rb.Push(Record{Payload: "Event 1", Timestamp: time.Now(), Count: 1})
	rb.Push(Record{Payload: "Event 2", Timestamp: time.Now(), Count: 1})
	rb.Push(Record{Payload: "Event 3", Timestamp: time.Now(), Count: 1})

	fmt.Println("\n==============================================")
	fmt.Println("--> [PEHLE 3 RECORDS] Tape Kaisi Dikh Rahi Hai:")

	for i, rec := range rb.GetAll() {
		fmt.Printf("    Index %d: %s\n", i, rec.Payload)
	}

	rb.Push(Record{Payload: "Event 4", Timestamp: time.Now(), Count: 1})

	fmt.Println("\n--> [4TH RECORD AANE KE BAAD] Overwrite Ka Jaadu:")
	for i, rec := range rb.GetAll() {
		fmt.Printf("    Index %d: %s\n", i, rec.Payload)
	}
	fmt.Println("==============================================")
	// Check: Size 3 hi rehna chahiye, 4 nahi hona chahiye!
	if rb.Size() != 3 {
		t.Fatalf("Expected size 3, got %d", rb.Size())
	}
}
