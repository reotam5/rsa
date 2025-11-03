package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"reotamai/assign4/internal/aescbc"
	"reotamai/assign4/internal/rsa"
)

func main() {
	inputImagePath := flag.String("input", "", "Path to the input image file")
	flag.Parse()

	if *inputImagePath == "" {
		fmt.Println("Usage: go run cmd/task4/image-encryption/main.go -input <image_path>")
		fmt.Println("Example: go run cmd/task4/image-encryption/main.go -input data/image.jpg")
		os.Exit(1)
	}

	if _, err := os.Stat(*inputImagePath); os.IsNotExist(err) {
		log.Fatalf("Error: Input image file '%s' does not exist", *inputImagePath)
	}

	fmt.Println("1. Loading original image...")
	originalImageData, _ := os.ReadFile(*inputImagePath)
	fmt.Printf("   Size: %d bytes\n", len(originalImageData))
	fmt.Println()

	fmt.Println("2. Generating RSA key pair for secure key exchange...")
	keyPair, _ := rsa.GenerateKeyPair(2048)
	fmt.Println("   RSA key pair generated successfully")
	fmt.Println()

	fmt.Println("3. Generating AES-256 key using BBS algorithm...")
	aesKey, _ := aescbc.GenerateKey()
	fmt.Printf("   AES Key (hex): %x\n", aesKey)
	fmt.Println()

	fmt.Println("4. Encrypting AES key with RSA public key...")
	encryptedAESKey, _ := rsa.Encrypt(&keyPair.PublicKey, aesKey)
	fmt.Printf("   Encrypted AES Key (hex): %x\n", encryptedAESKey)
	fmt.Println()

	fmt.Println("5. Decrypting AES key with RSA private key...")
	decryptedAESKey, _ := rsa.Decrypt(&keyPair.PrivateKey, encryptedAESKey)
	fmt.Printf("   Decrypted AES Key (hex): %x\n", decryptedAESKey)
	fmt.Println()

	fmt.Println("6. Encrypting image data using AES-256-CBC...")
	ciphertext, iv, _ := aescbc.Encrypt(decryptedAESKey, originalImageData)
	fmt.Printf("   IV (hex): %x\n", iv)
	fmt.Printf("   Encrypted image size: %d bytes\n", len(ciphertext))
	fmt.Println()

	encryptedFilePath := getOutputPath(*inputImagePath, "_encrypted.bin")
	fmt.Println("7. Saving encrypted image...")
	err := os.WriteFile(encryptedFilePath, ciphertext, 0644)
	if err != nil {
		log.Fatalf("Failed to save encrypted image: %v", err)
	}
	fmt.Printf("   Encrypted image saved to: %s\n", encryptedFilePath)

	// saving iv as well
	ivFilePath := getOutputPath(*inputImagePath, "_iv.hex")
	err = os.WriteFile(ivFilePath, []byte(hex.EncodeToString(iv)), 0644)
	if err != nil {
		log.Fatalf("Failed to save IV: %v", err)
	}
	fmt.Printf("   IV saved to: %s\n", ivFilePath)
	fmt.Println()

	fmt.Println("8. Loading encrypted image...")
	loadedCiphertext, _ := os.ReadFile(encryptedFilePath)
	loadedIVHex, _ := os.ReadFile(ivFilePath)
	loadedIV, err := hex.DecodeString(string(loadedIVHex))
	fmt.Printf("   Loaded encrypted image: %d bytes\n", len(loadedCiphertext))
	fmt.Printf("   Loaded IV: %x\n", loadedIV)
	fmt.Println()

	fmt.Println("9. Decrypting image data using AES-256-CBC...")
	decryptedImageData, _ := aescbc.Decrypt(decryptedAESKey, loadedIV, loadedCiphertext)
	fmt.Printf("   Decrypted image size: %d bytes\n", len(decryptedImageData))
	fmt.Println()

	fmt.Println("10. Saving decrypted image...")
	decryptedFilePath := getOutputPath(*inputImagePath, "_decrypted"+filepath.Ext(*inputImagePath))
	os.WriteFile(decryptedFilePath, decryptedImageData, 0644)
	if err != nil {
		log.Fatalf("Failed to save decrypted image: %v", err)
	}
	fmt.Printf("   Decrypted image saved to: %s\n", decryptedFilePath)
	fmt.Println()

	fmt.Println("11. Verification...")
	if len(decryptedImageData) != len(originalImageData) {
		fmt.Println("   FAILED: Decrypted image size does not match original")
		os.Exit(1)
	}

	match := true
	for i := range originalImageData {
		if decryptedImageData[i] != originalImageData[i] {
			match = false
			break
		}
	}

	if match {
		fmt.Println("   SUCCESS: Decrypted image data matches original image data")
	} else {
		fmt.Println("   FAILED: Decrypted image data does not match original image data")
	}
}

func getOutputPath(inputPath, suffix string) string {
	dir := filepath.Dir(inputPath)
	base := filepath.Base(inputPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	return filepath.Join(dir, name+suffix)
}
