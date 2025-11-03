package rsa

import (
	"math/big"

	"reotamai/assign4/internal/bbs"
	"reotamai/assign4/internal/millerrabin"
)

type KeyPair struct {
	PublicKey  PublicKey
	PrivateKey PrivateKey
}

type PublicKey struct {
	N *big.Int
	E *big.Int
}

type PrivateKey struct {
	N *big.Int
	D *big.Int
	P *big.Int
	Q *big.Int
}

// randomly generates public and private key pair
func GenerateKeyPair(bits int) (*KeyPair, error) {
	// for example, 2048 key requires p and q to be 1024 each
	primeBits := bits / 2

	p := generatePrime(primeBits)
	q := generatePrime(primeBits)

	// commonly used public exponent
	e := big.NewInt(65537)

	// n = p * q
	n := new(big.Int).Mul(p, q)

	// d = e^(-1) mod (p-1)(q-1)
	d := new(big.Int).ModInverse(e, new(big.Int).Mul(new(big.Int).Sub(p, big.NewInt(1)), new(big.Int).Sub(q, big.NewInt(1))))

	return &KeyPair{
		PublicKey: PublicKey{
			N: n,
			E: e,
		},
		PrivateKey: PrivateKey{
			N: n,
			D: d,
			P: p,
			Q: q,
		},
	}, nil
}

func generatePrime(bits int) *big.Int {
	for {
		bbsGen, _ := bbs.CreateBBSStruct()
		candidate := bbsGen.GenerateBigInt(bits)
		if millerrabin.IsPrime(candidate) {
			return candidate
		}
	}
}

// rsa encryption
func Encrypt(pubKey *PublicKey, message []byte) ([]byte, error) {
	// convert message to big integer
	m := new(big.Int).SetBytes(message)

	// cipher = message ^ e (mod n)
	ciphertext := new(big.Int).Exp(m, pubKey.E, pubKey.N)

	// convert back to bytes
	return ciphertext.Bytes(), nil
}

// rsa decryption
func Decrypt(privKey *PrivateKey, ciphertext []byte) ([]byte, error) {
	// convert ciphertext to big integer
	c := new(big.Int).SetBytes(ciphertext)

	// plaintext = cipher^d (mod n)
	plaintext := new(big.Int).Exp(c, privKey.D, privKey.N)

	// convert back to bytes and return
	return plaintext.Bytes(), nil
}

// rsa decryption using Chinese Remainder Theorem optimization
// step 1: compute dp = d mod (p - 1) and dq = d mod (q - 1)
// step 2: decrypt with mp = c^dp mod p
// step 3: decrypt with mq = c^dq mod q
// step 4: combine mp and mq with h = (mp - mq) * q inverse (mod p)
// step 5: final result is m = mq + q * h
func DecryptCRT(privKey *PrivateKey, ciphertext []byte) []byte {
	c := new(big.Int).SetBytes(ciphertext)

	p := privKey.P
	q := privKey.Q
	d := privKey.D

	pMinus1 := new(big.Int).Sub(p, big.NewInt(1))
	qMinus1 := new(big.Int).Sub(q, big.NewInt(1))
	dp := new(big.Int).Mod(d, pMinus1)
	dq := new(big.Int).Mod(d, qMinus1)

	// Compute mp and mq
	mp := new(big.Int).Exp(c, dp, p)
	mq := new(big.Int).Exp(c, dq, q)

	// combine with CRT
	qInv := new(big.Int).ModInverse(q, p)
	h := new(big.Int).Sub(mp, mq)
	h.Mul(h, qInv)
	h.Mod(h, p)

	plaintext := new(big.Int).Set(mq)
	plaintext.Add(plaintext, new(big.Int).Mul(h, q))

	return plaintext.Bytes()
}
