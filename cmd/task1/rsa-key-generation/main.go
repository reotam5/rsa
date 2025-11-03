package main

import (
	"fmt"
	"math/big"

	"reotamai/assign4/internal/bbs"
	"reotamai/assign4/internal/millerrabin"
)

func main() {
	p := generatePrime(256)
	q := generatePrime(256)

	fmt.Println("Generated 256-bit primes:")
	fmt.Printf("  p: %s\n", p.String())
	fmt.Printf("  q: %s\n", q.String())

	// commonly used public exponent
	e := big.NewInt(65537)

	// n = p * q
	n := new(big.Int).Mul(p, q)

	// d = e^(-1) mod (p-1)(q-1)
	d := new(big.Int).ModInverse(e, new(big.Int).Mul(new(big.Int).Sub(p, big.NewInt(1)), new(big.Int).Sub(q, big.NewInt(1))))

	fmt.Println("\nRSA Public Key Components:")
	fmt.Printf("  n (modulus): %s\n", n.String())
	fmt.Printf("  e (exponent): %s\n", e.String())

	fmt.Println("\nRSA Private Key Components:")
	fmt.Printf("  n (modulus): %s\n", n.String())
	fmt.Printf("  d (private exponent): %s\n", d.String())
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
