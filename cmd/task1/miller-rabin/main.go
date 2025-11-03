package main

import (
	"fmt"
	"math/big"
	"reotamai/assign4/internal/bbs"
	"reotamai/assign4/internal/millerrabin"
)

func main() {
	// generate and test a large prime candidate
	fmt.Println("\nGenerating a 256-bit prime candidate using BBS + Miller-Rabin...")
	var prime *big.Int
	attempts := 0

	for {
		attempts++
		bbsGen, _ := bbs.CreateBBSStruct()
		candidate := bbsGen.GenerateBigInt(256)
		if millerrabin.IsPrime(candidate) {
			prime = candidate
			break
		}
		if attempts%10 == 0 {
			fmt.Printf("  Attempt %d...\n", attempts)
		}
	}
	fmt.Printf("Found 256-bit prime after %d attempts: %s\n", attempts, prime.String())
}
