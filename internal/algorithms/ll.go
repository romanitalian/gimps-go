package algorithms

import (
	"fmt"
	"math/big"

	"github.com/romanitalian/gimps-go/internal/math"
)

// LLResult represents the result of a Lucas-Lehmer test
type LLResult struct {
	IsPrime bool
	Residue *math.BigInt // Final residue (64-bit)
	Error   error
}

// LucasLehmerTest performs the Lucas-Lehmer primality test for Mersenne number 2^p - 1
// Based on prime() function from commonb.c:6451
// Algorithm: s_0 = 4, s_i = (s_{i-1}^2 - 2) mod (2^p - 1) for i = 1 to p-2
// If s_{p-2} == 0, then 2^p - 1 is prime
func LucasLehmerTest(exponent uint64) (*LLResult, error) {
	if exponent < 2 {
		return &LLResult{IsPrime: false, Error: fmt.Errorf("exponent must be at least 2")}, nil
	}

	// Compute Mersenne number: M_p = 2^p - 1
	mp := math.MersenneNumber(exponent)

	// Initialize: s_0 = 4
	s := math.NewBigInt(4)

	// Perform p-2 iterations
	// s_i = (s_{i-1}^2 - 2) mod M_p
	for i := uint64(1); i < exponent-1; i++ {
		// s = s^2 mod M_p
		s = math.ModSqr(s, mp)
		// s = s - 2 mod M_p
		two := math.NewBigInt(2)
		s = math.ModSub(s, two, mp)
	}

	// Check if final residue is zero
	isPrime := s.IsZero()

	// Extract 64-bit residue
	residue := extractResidue64(s, exponent)

	return &LLResult{
		IsPrime: isPrime,
		Residue: residue,
	}, nil
}

// LucasLehmerTestWithProgress performs Lucas-Lehmer test with progress callback
// callback is called after each iteration with (iteration, total, current_residue)
func LucasLehmerTestWithProgress(exponent uint64, callback func(uint64, uint64, *math.BigInt) error) (*LLResult, error) {
	if exponent < 2 {
		return &LLResult{IsPrime: false, Error: fmt.Errorf("exponent must be at least 2")}, nil
	}

	// Compute Mersenne number: M_p = 2^p - 1
	mp := math.MersenneNumber(exponent)

	// Initialize: s_0 = 4
	s := math.NewBigInt(4)

	total := exponent - 2

	// Perform p-2 iterations
	for i := uint64(1); i < exponent-1; i++ {
		// s = s^2 mod M_p
		s = math.ModSqr(s, mp)
		// s = s - 2 mod M_p
		two := math.NewBigInt(2)
		s = math.ModSub(s, two, mp)

		// Call progress callback
		if callback != nil {
			if err := callback(i, total, s); err != nil {
				return nil, err
			}
		}
	}

	// Check if final residue is zero
	isPrime := s.IsZero()

	// Extract 64-bit residue
	residue := extractResidue64(s, exponent)

	return &LLResult{
		IsPrime: isPrime,
		Residue: residue,
	}, nil
}

// LucasLehmerTestFromWorkUnit performs Lucas-Lehmer test from work unit parameters
func LucasLehmerTestFromWorkUnit(k float64, b, n uint64, c int64) (*LLResult, error) {
	// Check if it's a Mersenne number (k=1, b=2, c=-1)
	if k != 1.0 || b != 2 || c != -1 {
		return nil, fmt.Errorf("work unit is not a Mersenne number")
	}

	if n == 0 {
		return nil, fmt.Errorf("invalid Mersenne exponent")
	}

	return LucasLehmerTest(n)
}

// extractResidue64 extracts a 64-bit residue from a BigInt
// This simulates the behavior of generateResidue64 from commonb.c:5929
func extractResidue64(value *math.BigInt, exponent uint64) *math.BigInt {
	// For simplicity, we'll extract the low 64 bits
	// In the original code, there's more complex logic involving shift counts
	mask := new(big.Int)
	mask.SetUint64(0xFFFFFFFFFFFFFFFF)
	result := &math.BigInt{Int: new(big.Int)}
	result.And(value.Int, mask)
	return result
}

// KnownMersennePrimes returns a list of known Mersenne prime exponents
// These are the currently known Mersenne primes (as of 2024)
func KnownMersennePrimes() []uint64 {
	return []uint64{
		2, 3, 5, 7, 13, 17, 19, 31, 61, 89, 107, 127,
		521, 607, 1279, 2203, 2281, 3217, 4253, 4423, 9689, 9941, 11213, 19937,
		21701, 23209, 44497, 86243, 110503, 132049, 216091, 756839, 859433,
		1257787, 1398269, 2976221, 3021377, 6972593, 13466917, 20996011,
		24036583, 25964951, 30402457, 32582657, 37156667, 42643801, 43112609,
		57885161, 74207281, 77232917, 82589933,
	}
}

// IsKnownMersennePrime checks if an exponent is a known Mersenne prime
func IsKnownMersennePrime(exponent uint64) bool {
	known := KnownMersennePrimes()
	for _, p := range known {
		if p == exponent {
			return true
		}
	}
	return false
}

