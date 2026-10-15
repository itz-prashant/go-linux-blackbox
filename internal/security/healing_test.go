package security

import (
	"fmt"
	"testing"
	"time"
	"crypto/ed25519"
)

func TestCryptographicSelfHealing(t *testing.T) {
	// 1. Ek authorized Ed25519 Key Pair generate kiya:
	pubKey, privKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	fmt.Println("\n==================================================================")
	fmt.Println("--> [SCENARIO 1] Authorized Control Plane Sends Signed Fix")

	// Authorized command sign ki gayi:
	validCmd := SignCommand("ACTION_FLUSH_BUFFERS", privKey)
	res, err := VerifyAndExecute(validCmd, pubKey)
	if err != nil {
		t.Fatalf("Expected success for valid signature, got: %v", err)
	}
	fmt.Printf("    Result: %s\n", res)

	// -------------------------------------------------------------------------
	fmt.Println("\n--> [SCENARIO 2] Hacker Tampering Attack (Altered Command)")

	// Hacker ne action badal kar "ACTION_HACK" kar diya:
	tamperedCmd := validCmd
	tamperedCmd.Action = "ACTION_HACK" // Signature ab invalid ho jayegi!

	_, err = VerifyAndExecute(tamperedCmd, pubKey)
	if err == nil {
		t.Fatalf("Expected signature failure for tampered command!")
	}
	fmt.Printf("    [SECURITY SHIELD] %v\n", err)

	// -------------------------------------------------------------------------
	fmt.Println("\n--> [SCENARIO 3] Replay Attack (Old Recorded Command Re-sent)")

	// Purani command jo 2 minute pehle (120 seconds ago) sign hui thi:
	expiredCmd := validCmd
	expiredCmd.Timestamp = time.Now().Unix() - 120 // 2 minutes in the past!
	// Update signature for this timestamp to show signature is valid but time expired:
	expiredMsg := serializeMessage(expiredCmd.Action, expiredCmd.Timestamp, expiredCmd.Nonce)
	expiredCmd.Signature = ed25519Sign(privKey, expiredMsg)

	_, err = VerifyAndExecute(expiredCmd, pubKey)
	if err == nil {
		t.Fatalf("Expected replay attack prevention!")
	}
	fmt.Printf("    [SECURITY SHIELD] %v\n", err)
	fmt.Println("==================================================================")
}

// Helper to test replay attack with expired valid signature
func ed25519Sign(priv []byte, msg []byte) []byte {
	return ed25519SignRaw(priv, msg)
}

func ed25519SignRaw(priv ed25519.PrivateKey, msg []byte) []byte {
	return ed25519.Sign(priv, msg)
}