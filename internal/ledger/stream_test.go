package ledger

import (
	"fmt"
	"testing"
	"time"
)

func TestStreamAckAndWatchdogRecovery(t *testing.T) {
	// 1. Ek naya stream ledger banaya:
	stream := NewStream()

	fmt.Println("\n==================================================================")
	fmt.Println("--> [SCENARIO 1] Worker A Normal Delivery (Successful ACK)")

	// Log 1: Worker A ne uthaya aur kaam khatam karke sign (ACK) kar diya:
	msg1 := stream.Add("LOG #1: User logged in successfully")
	fmt.Printf("    Message #%d dispatched to Worker A\n", msg1.ID)

	acked := stream.Ack(msg1.ID)
	fmt.Printf("    Worker A signed (ACK) Message #%d: %v\n", msg1.ID, acked)

	// -------------------------------------------------------------------------
	fmt.Println("\n--> [SCENARIO 2] Worker B Crashes Mid-Way! (Dead Worker Trap)")

	// Log 2: Worker B ne uthaya, par sign karne se pehle CRASH ho gaya:
	msg2 := stream.Add("LOG #2: CRITICAL: Payment gateway crashed!")
	fmt.Printf("    Message #%d dispatched to Worker B\n", msg2.ID)
	fmt.Println("    [BOOM] Worker B crashed and NEVER called Ack()!")

	// Abhi register mein 1 message PENDING hona chahiye:
	if stream.PendingCount() != 1 {
		t.Fatalf("Expected 1 pending message, got %d", stream.PendingCount())
	}
	fmt.Printf("    Current In-Flight Pending Messages: %d\n", stream.PendingCount())

	// -------------------------------------------------------------------------
	fmt.Println("\n--> [SCENARIO 3] Watchdog Inspector Jaagta Hai (Rescue Operation)")

	// Thoda waqt beeta (150ms) taaki timeout exceed ho jaye:
	time.Sleep(150 * time.Millisecond)

	// Watchdog 100ms se purane un-acked messages dhoondhta hai:
	orphaned := stream.GetPending(100 * time.Millisecond)

	fmt.Printf("    Watchdog found %d orphaned messages from dead workers!\n", len(orphaned))
	if len(orphaned) != 1 {
		t.Fatalf("Expected 1 orphaned message to be recovered, got %d", len(orphaned))
	}

	rescuedMsg := orphaned[0]
	fmt.Printf("    [RESCUED] Rescued Message #%d: \"%s\"\n", rescuedMsg.ID, rescuedMsg.Payload)

	// Worker C aakar us rescued message ko process karke sign (ACK) karta hai:
	stream.Ack(rescuedMsg.ID)
	fmt.Printf("    Worker C successfully processed and ACKed Message #%d!\n", rescuedMsg.ID)

	// Final check: Ab pending count 0 hona chahiye!
	if stream.PendingCount() != 0 {
		t.Fatalf("Expected 0 pending messages, got %d", stream.PendingCount())
	}
	fmt.Printf("\n--> [SUCCESS] Final Pending Count: %d (Zero Data Loss Guaranteed!)\n", stream.PendingCount())
	fmt.Println("==================================================================")
}