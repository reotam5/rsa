package main

import (
	"fmt"

	"reotamai/assign4/internal/aescbc"
	"reotamai/assign4/internal/rsa"
)

func main() {
	fmt.Println("1. Generating RSA Key Pair...")
	keyPair, _ := rsa.GenerateKeyPair(2048)
	pubKey := keyPair.PublicKey
	priKey := keyPair.PrivateKey

	fmt.Printf("  p: %s\n", trimOutput(priKey.P.String()))
	fmt.Printf("  q: %s\n", trimOutput(priKey.Q.String()))
	fmt.Println("  RSA Public Key:")
	fmt.Printf("    n: %s\n", trimOutput(pubKey.N.String()))
	fmt.Printf("    e: %s\n", pubKey.E)
	fmt.Println()
	fmt.Println("  RSA Private Key:")
	fmt.Printf("    n: %s\n", trimOutput(priKey.N.String()))
	fmt.Printf("    d: %s\n", trimOutput(priKey.D.String()))
	fmt.Println()

	fmt.Println("2. Generating AES-256 Key...")
	aesKey, _ := aescbc.GenerateKey()
	fmt.Printf("  Original AES Key (hex): %s\n", trimOutput(fmt.Sprintf("%x", aesKey)))
	fmt.Println()

	fmt.Println("3. Encrypting AES Key with RSA Public Key...")
	encryptedAESKey, _ := rsa.Encrypt(&keyPair.PublicKey, aesKey)
	fmt.Printf("  Encrypted AES Key (hex): %s\n", trimOutput(fmt.Sprintf("%x", encryptedAESKey)))
	fmt.Println()

	fmt.Println("4. Decrypting AES Key with RSA Private Key...")
	decryptedAESKey, _ := rsa.Decrypt(&keyPair.PrivateKey, encryptedAESKey)
	fmt.Printf("  Decrypted AES Key (hex): %s\n", trimOutput(fmt.Sprintf("%x", decryptedAESKey)))
	fmt.Println()

	fmt.Println("5. Verification...")
	if len(decryptedAESKey) == len(aesKey) {
		match := true
		for i := range aesKey {
			if decryptedAESKey[i] != aesKey[i] {
				match = false
				break
			}
		}
		if match {
			fmt.Println("  SUCCESS: Decrypted AES key matches the original key")
		} else {
			fmt.Println("  FAILED: Decrypted AES key does not match the original key")
		}
	} else {
		fmt.Println("  FAILED: Decrypted AES key length does not match the original key")
	}
}

func trimOutput(s string) string {
	return s[:10] + "..." + s[len(s)-10:]
}
