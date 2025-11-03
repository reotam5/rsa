package main

import (
	"fmt"

	"reotamai/assign4/internal/rsa"
)

func main() {
	keyPair, _ := rsa.GenerateKeyPair(512)
	public := keyPair.PublicKey
	private := keyPair.PrivateKey

	fmt.Println("Generated 256-bit primes:")
	fmt.Printf("  p: %s\n", private.P.String())
	fmt.Printf("  q: %s\n", private.Q.String())

	fmt.Println("\nRSA Public Key Components:")
	fmt.Printf("  n: %s\n", public.N.String())
	fmt.Printf("  e: %s\n", public.E.String())

	fmt.Println("\nRSA Private Key Components:")
	fmt.Printf("  n: %s\n", private.N.String())
	fmt.Printf("  d: %s\n", private.D.String())
}
