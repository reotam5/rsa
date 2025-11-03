package bbs

import (
	"crypto/rand"
	"math/big"
)

type BBS struct {
	n *big.Int // n = p * q
	s *big.Int // seed
}

// creates a new BBS struct based on two large primes p and q, both ≡ 3 (mod 4)
func CreateBBSStruct() (*BBS, error) {
	p, err := generateBlumPrime(512)
	if err != nil {
		return nil, err
	}

	q, err := generateBlumPrime(512)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).Mul(p, q)

	// generate seed that is coprime to n (not divisible by p or q)
	var s *big.Int
	for {
		seed, err := rand.Int(rand.Reader, n)
		if err != nil {
			return nil, err
		}
		// ensure seed is not 0, 1, or divisible by p or q
		if seed.Cmp(big.NewInt(0)) != 0 &&
			seed.Cmp(big.NewInt(1)) != 0 &&
			new(big.Int).Mod(seed, p).Cmp(big.NewInt(0)) != 0 &&
			new(big.Int).Mod(seed, q).Cmp(big.NewInt(0)) != 0 {
			s = seed
			break
		}
	}

	return &BBS{
		n: n,
		s: s,
	}, nil
}

// generates a prime number p such that p ≡ 3 (mod 4)
func generateBlumPrime(bits int) (*big.Int, error) {
	for {
		// a random prime number of specified bit length
		prime, err := rand.Prime(rand.Reader, bits)
		if err != nil {
			return nil, err
		}

		// check if prime ≡ 3 (mod 4)
		mod4 := new(big.Int).Mod(prime, big.NewInt(4))
		if mod4.Cmp(big.NewInt(3)) == 0 {
			return prime, nil
		}
	}
}

// generates n random bits using the BBS struct
// each bit is the least significant bit of x_next = (x_prev)^2 mod n
func (b *BBS) GenerateBits(n int) []byte {
	bits := make([]byte, n)
	x := new(big.Int).Set(b.s)

	for i := 0; i < n; i++ {
		x.Mul(x, x)
		x.Mod(x, b.n)

		// get the least significant bit
		if x.Bit(0) == 1 {
			bits[i] = 1
		} else {
			bits[i] = 0
		}
	}

	return bits
}

// generates n random bytes using the BBS generator.
// this will combine 8 bits into a single byte
func (b *BBS) GenerateBytes(n int) []byte {
	bits := b.GenerateBits(n * 8)
	bytes := make([]byte, n)

	for i := 0; i < n; i++ {

		// after this loop, []byte{1, 1, 1, 1, 1, 1, 1, 1} will be []byte{255}
		for j := 0; j < 8; j++ {
			if bits[i*8+j] == 1 {
				bytes[i] |= 1 << (7 - j)
			}
		}
	}

	return bytes
}

// generates a random big integer of specified bit length.
// max is 2^bits (exclusive)
// if bits is 2, 2 ^ 2 = 4, so random number will be in [0, 4)
func (b *BBS) GenerateBigInt(bits int) *big.Int {
	max := new(big.Int).Lsh(big.NewInt(1), uint(bits))
	return b.GenerateBigIntInRange(big.NewInt(0), max)
}

// generates a random big integer in the range [min, max).
func (b *BBS) GenerateBigIntInRange(min, max *big.Int) *big.Int {
	rangeSize := new(big.Int).Sub(max, min)

	// this should never happen, but just in case
	if rangeSize.Cmp(big.NewInt(0)) <= 0 {
		return new(big.Int).Set(min)
	}

	// determine number of bits needed
	bitsNeeded := rangeSize.BitLen()

	// this should never happen, but just in case
	if bitsNeeded == 0 {
		return new(big.Int).Set(min)
	}

	// generate enough bytes
	bytesNeeded := (bitsNeeded + 7) / 8 // +7 to round up to the nearest byte (eg, 1 bit -> (1 + 7) / 8 = 1 byte)
	randomBytes := b.GenerateBytes(bytesNeeded)

	// convert to big.Int and ensure it's in range
	randomInt := new(big.Int).SetBytes(randomBytes)
	randomInt.Mod(randomInt, rangeSize)
	randomInt.Add(randomInt, min)

	return randomInt
}

// getter for n
func (b *BBS) GetN() *big.Int {
	return new(big.Int).Set(b.n)
}

// getter for s
func (b *BBS) GetS() *big.Int {
	return new(big.Int).Set(b.s)
}
