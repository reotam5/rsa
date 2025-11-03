package millerrabin

import (
	"crypto/rand"
	"math/big"
)

// returns true if n is probably prime, false if n is definitely composite
// this only does 1 round. use isPrime function for multiple rounds
func MillerRabinTest(n *big.Int) bool {
	// edge cases: n < 2, n = 2 or 3, n is even
	if n.Cmp(big.NewInt(2)) < 0 {
		return false
	}
	if n.Cmp(big.NewInt(2)) == 0 || n.Cmp(big.NewInt(3)) == 0 {
		return true
	}
	if new(big.Int).Mod(n, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return false
	}

	// find k, q with k > 0, q odd, (n - 1) = (2^k)*q
	nMinus1 := new(big.Int).Sub(n, big.NewInt(1))
	k := 0
	q := new(big.Int).Set(nMinus1)

	// while q is even
	for new(big.Int).Mod(q, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		k++
		q.Rsh(q, 1) // this effectively divides by 2
	}

	// pick a random base a in (1, n-1) or in other words, [2, n-1)
	rangeVal := new(big.Int).Sub(nMinus1, big.NewInt(2))
	a, _ := rand.Int(rand.Reader, rangeVal)
	a.Add(a, big.NewInt(2))

	// x = a^q mod n
	x := new(big.Int).Exp(a, q, n)

	// if x == 1 or x == n-1, n might be prime
	if x.Cmp(big.NewInt(1)) == 0 || x.Cmp(nMinus1) == 0 {
		return true
	}

	for j := 0; j < k; j++ {
		// x = x^2 mod n effectively gives us a^(2^j * q) mod n
		x.Mul(x, x).Mod(x, n)
		if x.Cmp(nMinus1) == 0 {
			return true
		}
	}

	return false
}

// IsPrime performs miller rabin test with 64 rounds
func IsPrime(n *big.Int) bool {
	for i := 0; i < 64; i++ {
		if !MillerRabinTest(n) {
			return false
		}
	}
	return true
}
