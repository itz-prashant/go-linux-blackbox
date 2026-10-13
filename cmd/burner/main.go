package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("--> [Burner] allocating 1GB RAM ")

	var chunks [][]byte

	for i := 0; i < 1000; i++ {
		block := make([]byte, 1024*1024)

		for j := 0; j < len(block); j += 4096 {
			block[j] = 1
		}
		chunks = append(chunks, block)
	}
	fmt.Println("--> [BURNER] 1GB RAM FULL! Ab 15 seconds ke liye CPU aur RAM par load...")

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		_ = 1 + 1
	}
	fmt.Printf("--> [BURNER] Complete! Releasing %d MB RAM.\n", len(chunks))
}
