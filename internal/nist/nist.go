package nist

import (
	"math"
)

// holds test results
type NISTTestResults struct {
	FrequencyTest       FrequencyTestResult
	RunsTest            RunsTestResult
	MaurerUniversalTest MaurerUniversalTestResult
}

// checks if ones and zeros occur with the same proportion
type FrequencyTestResult struct {
	Passed bool
	PValue float64 // PValue >= 0.01 is accepted as random
}

// checks for number and average length of runs
type RunsTestResult struct {
	Passed bool
	PValue float64 // PValue >= 0.01 is accepted as random
}

// checks whether the sequence can be significantly compressed
type MaurerUniversalTestResult struct {
	Passed bool
	PValue float64 // PValue >= 0.01 is accepted as random
}

// runs all tests
func RunAllTests(bits []byte) NISTTestResults {
	return NISTTestResults{
		FrequencyTest:       FrequencyTest(bits),
		RunsTest:            RunsTest(bits),
		MaurerUniversalTest: MaurerUniversalTest(bits),
	}
}

// performs frequency test
func FrequencyTest(bits []byte) FrequencyTestResult {
	n := len(bits)
	if n == 0 {
		return FrequencyTestResult{
			Passed: false,
		}
	}

	//sum of normalized bits (+1 for 1, -1 for 0)
	sn := float64(0)
	for _, bit := range bits {
		if bit == 1 {
			sn++
		} else {
			sn--
		}
	}

	// from spec doc
	sObs := math.Abs(sn) / math.Sqrt(float64(n))
	pValue := math.Erfc(sObs / math.Sqrt(2))

	// test passes if PValue >= 0.01
	passed := pValue >= 0.01

	return FrequencyTestResult{
		Passed: passed,
		PValue: pValue,
	}
}

// performs the Runs Test
func RunsTest(bits []byte) RunsTestResult {
	n := len(bits)
	if n < 2 {
		return RunsTestResult{
			Passed: false,
		}
	}

	// count 1s
	ones := 0
	for _, bit := range bits {
		if bit == 1 {
			ones++
		}
	}

	// pi is the proportion of 1s
	pi := float64(ones) / float64(n)

	// count the number of runs
	vn := 1
	for i := 1; i < n; i++ {
		if bits[i] != bits[i-1] {
			vn++
		}
	}

	// from spec doc
	num := math.Abs(float64(vn) - 2.0*float64(n)*pi*(1.0-pi))
	den := 2.0 * math.Sqrt(2.0*float64(n)) * pi * (1.0 - pi)
	pValue := math.Erfc(num / den)

	// test passes if PValue >= 0.01
	passed := pValue >= 0.01

	return RunsTestResult{
		Passed: passed,
		PValue: pValue,
	}
}

// performs the Maurer's Universal Statistical Test
func MaurerUniversalTest(bits []byte) MaurerUniversalTestResult {
	n := len(bits)

	// following recommendation for L size from spec doc
	var L int
	if n >= 387840 {
		L = 6
	} else if n >= 904960 {
		L = 7
	} else if n >= 2068480 {
		L = 8
	} else if n >= 4654080 {
		L = 9
	} else if n >= 10342400 {
		L = 10
	} else if n >= 22753280 {
		L = 11
	} else if n >= 49643520 {
		L = 12
	} else if n >= 107560960 {
		L = 13
	} else if n >= 231669760 {
		L = 14
	} else if n >= 496435200 {
		L = 15
	} else if n >= 1059061760 {
		L = 16
	} else {
		L = 6
	}

	// Q = 10 * 2^L (number of initialization blocks)
	Q := 10 * (1 << L)

	// K = floor((n / L) - Q) (number of test blocks)
	K := (n / L) - Q

	if K <= 0 || Q <= 0 {
		return MaurerUniversalTestResult{
			Passed: false,
		}
	}

	// n bits to L-bits blocks (discards remainder)
	blocks := make([]int, n/L)
	for i := 0; i < n/L; i++ {
		block := 0
		for j := 0; j < L && i*L+j < n; j++ {
			if bits[i*L+j] == 1 {
				block |= 1 << (L - j - 1)
			}
		}
		blocks[i] = block
	}

	// table T stores last occurence index (table size is 2^L)
	// eg) index 0(1) stores last occurrence of 00 if L is 2
	T := make([]int, 1<<L)

	// first Q blocks prepares table T
	for i := 0; i < Q; i++ {
		T[blocks[i]] = i + 1 // 1 indexed array
	}

	// remaining K blocks are similar but with running sum of log2 distance
	sum := 0.0
	for i := Q; i < Q+K; i++ {
		pattern := blocks[i]
		lastOccurrence := T[pattern]

		distance := i + 1 - lastOccurrence
		sum += math.Log2(float64(distance))

		// update last occurrence
		T[pattern] = i + 1
	}

	// calculate fn = (1/K) * sum of log2 distances
	fn := sum / float64(K)

	// lookup expected value and variance based on L
	expectedValue, variance := getExpectedValueAndVariance(L)
	if expectedValue == 0 {
		// L value not in lookup table
		return MaurerUniversalTestResult{
			Passed: false,
		}
	}

	// calculate c = 0.7 - (0.8 / L) + (4 + 32/L) * (K^(-3/L)) / 15)
	c := 0.7 - (0.8 / float64(L)) + (4.0+32.0/float64(L))*(math.Pow(float64(K), -3.0/float64(L))/15.0)

	// calculate sigma = c * sqrt(variance / K)
	sigma := c * math.Sqrt(variance/float64(K))

	// calculate P-value = erfc(|fn - expectedValue / sqrt(2 * sigma)|)
	pValue := math.Erfc(math.Abs((fn - expectedValue) / math.Sqrt(2.0*sigma)))

	// test passes if pValue >= 0.01
	passed := pValue >= 0.01

	return MaurerUniversalTestResult{
		Passed: passed,
		PValue: pValue,
	}
}

// returns the precomputed expected value and variance (taken from spec doc)
func getExpectedValueAndVariance(L int) (expectedValue float64, variance float64) {
	values := map[int]struct {
		expected float64
		variance float64
	}{
		6:  {5.2177052, 2.954},
		7:  {6.1962507, 3.125},
		8:  {7.1836656, 3.238},
		9:  {8.1764248, 3.311},
		10: {9.1723243, 3.356},
		11: {10.170032, 3.384},
		12: {11.168765, 3.401},
		13: {12.168070, 3.410},
		14: {13.167693, 3.416},
		15: {14.167488, 3.419},
		16: {15.167379, 3.421},
	}

	if val, ok := values[L]; ok {
		return val.expected, val.variance
	}
	return 0, 0
}
