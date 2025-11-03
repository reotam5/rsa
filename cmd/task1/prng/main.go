package main

import (
	"fmt"
	"log"
	"reotamai/assign4/internal/bbs"
)

func main() {
	fmt.Println("\nCreating BBS generator...")
	bbsGen, err := bbs.CreateBBSStruct()
	if err != nil {
		log.Fatal("Failed to create BBS generator:", err)
	}

	fmt.Println("BBS Generator created successfully!")
	fmt.Printf("Modulus n: %s\n", bbsGen.GetN().String())
	fmt.Printf("Initial seed s: %s\n", bbsGen.GetS().String())
	fmt.Println()

	// generate some random bits
	fmt.Println("Generating 1000 random bits...")
	bits := bbsGen.GenerateBits(1000)
	fmt.Printf("First 100 bits: ")
	for i := 0; i < 100 && i < len(bits); i++ {
		fmt.Printf("%d", bits[i])
		if (i+1)%10 == 0 {
			fmt.Printf(" ")
		}
	}
	fmt.Println()

	// generate random integer with 256 bit
	fmt.Println("\nGenerating a random 256-bit number...")
	randomNum := bbsGen.GenerateBigInt(256)
	fmt.Printf("Random 256-bit number: %s\n", randomNum.String())
}
