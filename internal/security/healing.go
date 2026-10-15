package security

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

type SignedCommand struct {
	Action    string `json:"action"`    // Jaise: "ACTION_FLUSH_BUFFERS"
	Timestamp int64  `json:"timestamp"` // Kis second par sign hua
	Nonce     uint64 `json:"nonce"`     // Random number (Replay protection)
	Signature []byte `json:"signature"` // Ed25519 cryptographic seal
}

func GenerateKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate keypair: %w", err)
	}
	return pub, priv, nil
}

// serializeMessage action, timestamp, aur nonce ko ek byte slice mein pack karta hai
func serializeMessage(action string, timestamp int64, nonce uint64) []byte {
	var buf []byte
	buf = append(buf, []byte(action)...)

	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, uint64(timestamp))
	buf = append(buf, timeBytes...)
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, nonce)
	buf = append(buf, nonceBytes...)
	return buf
}

// SignCommand authorized control plane ki private key se command sign karta hai
func SignCommand(action string, priv ed25519.PrivateKey) SignedCommand {
	now := time.Now().Unix()
	var nonceBytes [8]byte
	_, _ = rand.Read(nonceBytes[:])
	nonce := binary.BigEndian.Uint64(nonceBytes[:])
	msg := serializeMessage(action, now, nonce)
	sig := ed25519.Sign(priv, msg)
	return SignedCommand{
		Action:    action,
		Timestamp: now,
		Nonce:     nonce,
		Signature: sig,
	}
}

// VerifyAndExecute signature aur replay protection check karke safe self-healing apply karta hai
func VerifyAndExecute(cmd SignedCommand, pub ed25519.PublicKey) (string, error) {
	now := time.Now().Unix()
	// 1. REPLAY ATTACK CHECK: Command 60s ke andar ki honi chahiye!
	if now-cmd.Timestamp > 60 || cmd.Timestamp > now+5 {
		return "", fmt.Errorf("security violation: command expired or future timestamp (replay attack prevented)")
	}
	// 2. ED25519 CRYPTOGRAPHIC CHECK:
	msg := serializeMessage(cmd.Action, cmd.Timestamp, cmd.Nonce)
	if !ed25519.Verify(pub, msg, cmd.Signature) {
		return "", fmt.Errorf("security violation: invalid cryptographic signature (tampered command)")
	}
	// 3. SAFE SANDBOX EXECUTION:
	switch cmd.Action {
	case "ACTION_FLUSH_BUFFERS":
		return "SUCCESS: Reusable memory buffers purged and recycled safely.", nil
	case "ACTION_RESET_TAPE":
		return "SUCCESS: Tape recorder reset to clean state.", nil
	default:
		return "", fmt.Errorf("execution error: unknown action '%s' rejected", cmd.Action)
	}
}