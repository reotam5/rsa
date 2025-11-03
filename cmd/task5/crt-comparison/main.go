package main

import (
	"fmt"
	"time"

	"reotamai/assign4/internal/aescbc"
	"reotamai/assign4/internal/rsa"
)

func main() {
	fmt.Println("1. Encrypting AES Key with RSA Public Key...")
	keyPair, _ := rsa.GenerateKeyPair(2048)
	aesKey, _ := aescbc.GenerateKey()
	encryptedAESKey, _ := rsa.Encrypt(&keyPair.PublicKey, aesKey)
	fmt.Printf("  Encrypted AES Key (hex): %x\n", encryptedAESKey)
	fmt.Println()

	fmt.Println("2. Decrypting AES key using standard RSA Decryption...")
	startStandard := time.Now()
	decryptedAESKeyStandard, _ := rsa.Decrypt(&keyPair.PrivateKey, encryptedAESKey)
	durationStandard := time.Since(startStandard)
	fmt.Printf("  Time: %v\n", durationStandard)
	fmt.Printf("  Decrypted AES Key (hex): %x\n", decryptedAESKeyStandard)
	fmt.Println()

	fmt.Println("3. Decrypting AES key using CRT optimized RSA Decryption...")
	startCRT := time.Now()
	decryptedAESKeyCRT := rsa.DecryptCRT(&keyPair.PrivateKey, encryptedAESKey)
	durationCRT := time.Since(startCRT)
	fmt.Printf("  Time: %v\n", durationCRT)
	fmt.Printf("  Decrypted AES Key (hex): %x\n", decryptedAESKeyCRT)
	fmt.Println()

	speedup := float64(durationStandard) / float64(durationCRT)
	fmt.Printf("4. Performance Comparison:\n")
	fmt.Printf("  Standard Decryption: %v\n", durationStandard)
	fmt.Printf("  CRT-Optimized Decryption: %v\n", durationCRT)
	fmt.Printf("  Speedup: %.2fx faster with CRT\n", speedup)
	fmt.Println()
}
