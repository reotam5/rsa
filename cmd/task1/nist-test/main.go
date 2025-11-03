package main

import (
	"fmt"
	"log"
	"reotamai/assign4/internal/bbs"
	"reotamai/assign4/internal/nist"
)

func main() {
	fmt.Println("\nCreating BBS generator for randomness testing...")
	bbsGen, err := bbs.CreateBBSStruct()
	if err != nil {
		log.Fatal("Failed to create BBS generator:", err)
	}

	// frequency test
	fmt.Println("\nFrequency Test")
	fmt.Println("-----------------------------------------------------------------------")
	testBits := bbsGen.GenerateBits(100000)
	result := nist.FrequencyTest(testBits)
	fmt.Printf("P-value: %.6f\n", result.PValue)
	if result.Passed {
		fmt.Println("PASSED (p-value >= 0.01)")
	} else {
		fmt.Println("FAILED (p-value < 0.01)")
	}

	// runs test
	fmt.Println("\nRuns Test")
	fmt.Println("-----------------------------------------------------------------------")
	runsResult := nist.RunsTest(testBits)
	fmt.Printf("P-value: %.6f\n", runsResult.PValue)
	if runsResult.Passed {
		fmt.Println("PASSED (p-value >= 0.01)")
	} else {
		fmt.Println("FAILED (p-value < 0.01)")
	}

	// Maurer's Universal Statistical Test
	fmt.Println("\nMaurer's Universal Statistical Test")
	fmt.Println("-----------------------------------------------------------------------")
	maurerResult := nist.MaurerUniversalTest(testBits)
	fmt.Printf("P-value: %.6f\n", maurerResult.PValue)
	if maurerResult.Passed {
		fmt.Println("PASSED (p-value >= 0.01)")
	} else {
		fmt.Println("FAILED (p-value < 0.01)")
	}

	// summary
	fmt.Println("\n---------------------------------------------------------------------")
	fmt.Println("Summary: ")
	fmt.Printf("Frequency Test:  %s\n", passFail(result.Passed))
	fmt.Printf("Runs Test:       %s\n", passFail(runsResult.Passed))
	fmt.Printf("Maurer's Test:   %s\n", passFail(maurerResult.Passed))
}

func passFail(passed bool) string {
	if passed {
		return "PASSED"
	}
	return "FAILED"
}
