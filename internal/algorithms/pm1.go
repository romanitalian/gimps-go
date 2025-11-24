package algorithms

import (
	"fmt"

	"github.com/romanitalian/gimps-go/internal/math"
)

// PMinus1Result represents the result of P-1 factoring
type PMinus1Result struct {
	Factor *math.BigInt
	Stage  int // 1 or 2
	Error  error
}

// PMinus1Factor performs P-1 factoring
// Based on pminus1() function from ecm.cpp
// Stage 1: Compute a^E mod N where E is product of prime powers <= B1
// Stage 2: Use additional primes between B1 and B2
func PMinus1Factor(n *math.BigInt, B1, B2 uint64) (*PMinus1Result, error) {
	if n.IsZero() || n.Sign() < 0 {
		return nil, fmt.Errorf("n must be positive")
	}

	if B1 == 0 {
		return nil, fmt.Errorf("B1 must be greater than 0")
	}

	// Stage 1: Compute a^E mod N
	// E = product of all prime powers <= B1
	stage1Result, err := pm1Stage1(n, B1)
	if err != nil {
		return nil, err
	}

	// Check for factor after stage 1
	// GCD(x-1, N)
	xMinus1 := math.ModSub(stage1Result, math.NewBigInt(1), n)
	factor := math.GCD(xMinus1, n)
	
	if !factor.IsOne() && factor.Cmp(n.Int) != 0 {
		return &PMinus1Result{
			Factor: factor,
			Stage:  1,
		}, nil
	}

	// Stage 2: Use primes between B1 and B2
	if B2 > B1 {
		stage2Result, err := pm1Stage2(n, stage1Result, B1, B2)
		if err != nil {
			return nil, err
		}

		// Check for factor after stage 2
		// GCD(stage2Result-1, N)
		stage2Minus1 := math.ModSub(stage2Result, math.NewBigInt(1), n)
		factor = math.GCD(stage2Minus1, n)
		
		if !factor.IsOne() && factor.Cmp(n.Int) != 0 {
			return &PMinus1Result{
				Factor: factor,
				Stage:  2,
			}, nil
		}
	}

	// No factor found
	return &PMinus1Result{
		Factor: nil,
		Stage:  2,
	}, nil
}

// PMinus1FactorFromParams performs P-1 factoring from parameters
func PMinus1FactorFromParams(k float64, b, n uint64, c int64, B1, B2 uint64) (*PMinus1Result, error) {
	// Construct the number k*b^n+c
	nBig := constructNumberFromParams(k, b, n, c)
	if nBig == nil {
		return nil, fmt.Errorf("failed to construct number")
	}

	if B1 == 0 {
		B1 = estimateB1(nBig)
	}

	if B2 == 0 {
		B2 = B1 * 20
	}

	return PMinus1Factor(nBig, B1, B2)
}

// PMinus1FactorFromWorkUnit performs P-1 factoring from work unit parameters
func PMinus1FactorFromWorkUnit(k float64, b, n uint64, c int64, B1, B2 uint64) (*PMinus1Result, error) {
	return PMinus1FactorFromParams(k, b, n, c, B1, B2)
}

// pm1Stage1 performs P-1 stage 1
// Computes a^E mod N where E is the product of all prime powers <= B1
func pm1Stage1(n *math.BigInt, B1 uint64) (*math.BigInt, error) {
	// Start with a = 3 (or any small prime)
	a := math.NewBigInt(3)

	// Compute E = product of all prime powers <= B1
	E := computePrimePowerProduct(B1)

	// Compute a^E mod N
	result := math.ModPow(a, E, n)

	return result, nil
}

// pm1Stage2 performs P-1 stage 2
// Uses additional primes between B1 and B2
func pm1Stage2(n *math.BigInt, stage1Result *math.BigInt, B1, B2 uint64) (*math.BigInt, error) {
	// For simplicity, we'll use a basic approach
	// In the full implementation, this uses pairing and Lucas chains
	
	result := stage1Result.Clone()

	// Process primes between B1 and B2
	// This is a simplified version - the full implementation is more complex
	for p := B1 + 1; p <= B2 && p <= B1+1000; p++ {
		if isPrimeSimple(p) {
			// result = result^p mod N
			pBig := math.NewBigIntFromUint64(p)
			result = math.ModPow(result, pBig, n)
		}
	}

	return result, nil
}


