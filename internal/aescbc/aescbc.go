package aescbc

import (
	"crypto/aes"
	"fmt"

	"reotamai/assign4/internal/bbs"
)

const (
	// 32-byte (256-bit) key
	AESKeySize = 32
	// aes block size is 16 bytes
	AESBlockSize = 16
)

// generate random key using bbs
func GenerateKey() ([]byte, error) {
	bbsGen, err := bbs.CreateBBSStruct()
	if err != nil {
		return nil, fmt.Errorf("failed to create BBS generator: %w", err)
	}
	key := bbsGen.GenerateBytes(AESKeySize)
	return key, nil
}

func GenerateIV() ([]byte, error) {
	bbsGen, err := bbs.CreateBBSStruct()
	if err != nil {
		return nil, fmt.Errorf("failed to create BBS generator: %w", err)
	}
	iv := bbsGen.GenerateBytes(AESBlockSize)
	return iv, nil
}

// CBC encryption: given a key and a plaintext, return ciphertext, random iv used
func Encrypt(key []byte, plaintext []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// initial vector
	iv, err := GenerateIV()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create initial vector")
	}

	// pad the plaintext to be a multiple of block size
	// even if it is multiple of 16 already, we add padding so that we can remove padding when decrypting
	padding := AESBlockSize - (len(plaintext) % AESBlockSize)
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding) // pad with padding size number
	}
	paddedPlaintext := append(plaintext, padText...)

	ciphertext := make([]byte, len(paddedPlaintext))
	prevBlock := iv

	// loop of CBC implementation
	for i := 0; i < len(paddedPlaintext); i += AESBlockSize {
		currentBlock := paddedPlaintext[i : i+AESBlockSize]

		// XOR current block with previous ciphertext block (or IV for first block)
		xoredBlock := make([]byte, AESBlockSize)
		for j := 0; j < AESBlockSize; j++ {
			xoredBlock[j] = currentBlock[j] ^ prevBlock[j]
		}

		// encrypt the XORed block
		encryptedBlock := make([]byte, AESBlockSize)
		block.Encrypt(encryptedBlock, xoredBlock)

		// store encrypted block in ciphertext
		copy(ciphertext[i:i+AESBlockSize], encryptedBlock)

		// update previous block for next iteration
		prevBlock = encryptedBlock
	}

	return ciphertext, iv, nil
}

// CBC decryption: given key, iv, and cipher, returns plaintext
func Decrypt(key []byte, iv []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	plaintext := make([]byte, len(ciphertext))
	prevBlock := iv

	// loop of CBC decryption
	for i := 0; i < len(ciphertext); i += AESBlockSize {
		currentCipherBlock := ciphertext[i : i+AESBlockSize]

		decryptedBlock := make([]byte, AESBlockSize)
		block.Decrypt(decryptedBlock, currentCipherBlock)

		// XOR decrypted block with previous ciphertext block (or IV for first block)
		for j := 0; j < AESBlockSize; j++ {
			plaintext[i+j] = decryptedBlock[j] ^ prevBlock[j]
		}

		// update previous block for next iteration
		prevBlock = currentCipherBlock
	}

	// remove padding
	padding := int(plaintext[len(plaintext)-1])
	unpaddedPlaintext := plaintext[:len(plaintext)-padding]

	return unpaddedPlaintext, nil
}
