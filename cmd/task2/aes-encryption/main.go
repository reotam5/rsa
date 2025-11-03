package main

import (
	"fmt"
	"log"

	"reotamai/assign4/internal/aescbc"
)

func main() {
	fmt.Println("Generating 256-bit AES key using BBS algorithm...")
	aesKey, err := aescbc.GenerateKey()
	if err != nil {
		log.Fatal("Failed to generate AES key:", err)
	}
	fmt.Printf("   Key (hex): %x\n\n", aesKey)

	sampleMessage := "My name is Reo Tamai. I am a BCIT student."
	fmt.Printf("Sample message to encrypt:\n")
	fmt.Printf("   %s\n\n", sampleMessage)

	fmt.Println("Encrypting message using AES-256-CBC...")
	ciphertext, iv, err := aescbc.Encrypt(aesKey, []byte(sampleMessage))
	if err != nil {
		log.Fatal("Failed to encrypt message:", err)
	}
	fmt.Printf("   IV (hex): %x\n", iv)
	fmt.Printf("   Ciphertext (hex): %x\n\n", ciphertext)

	fmt.Println("Decrypting message...")
	plaintext, err := aescbc.Decrypt(aesKey, iv, ciphertext)
	if err != nil {
		log.Fatal("Failed to decrypt message:", err)
	}
	fmt.Printf("   Decrypted message: %s\n\n", string(plaintext))

	fmt.Println("Verification:")
	if string(plaintext) == sampleMessage {
		fmt.Println("   SUCCESS: plaintexts match after encryption and decryption")
	} else {
		fmt.Println("   FAILED: Decrypted message does not match original message")
	}
}
