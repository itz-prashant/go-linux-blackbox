package tape

import (
	"fmt"
	"testing"
	"time"
)

func TestDuplicateFolding(t *testing.T) {
	// 1. 10 slots ki tape banayi:
	ring := NewRingBuffer(10)

	// 2. 5-second ka Folder initialize kiya:
	folder := NewEventFolder(ring, 5*time.Second)

	now := time.Now()
	duplicateError := "FATAL: Database connection timeout on port 5432"

	fmt.Println("\n========================================================")
	fmt.Println("--> [TEST] 1,000 Duplicate Errors Bhej Rahe Hain...")

	// 3. Ek loop mein 1,000 baar identical error feed karo:
	for i := 0; i < 1000; i++ {
		// Har log 1 millisecond ke gap par aaya (5s window ke andar):
		folder.Feed(duplicateError, now.Add(time.Duration(i)*time.Millisecond))
	}

	// 4. Flush karo (jo pending hai use tape mein daalo):
	folder.Flush()

	// 5. Tape se records nikaal kar dekho:
	records := ring.GetAll()

	fmt.Printf("--> [RESULT] Tape Mein Total Records: %d\n", len(records))
	if len(records) > 0 {
		fmt.Printf("    Payload: \"%s\"\n", records[0].Payload)
		fmt.Printf("    Repeated Count: %d Times!\n", records[0].Count)
	}
	fmt.Println("========================================================")

	// Verifications:
	// 1. Tape par 1000 entries nahi honi chahiye, sirf 1 honi chahiye!
	if len(records) != 1 {
		t.Fatalf("Expected 1 folded record, got %d", len(records))
	}

	// 2. Us single entry ka count 1000 hona chahiye!
	if records[0].Count != 1000 {
		t.Fatalf("Expected count 1000, got %d", records[0].Count)
	}
}