package algorithms

import (
	"fmt"

	"github.com/romanitalian/gimps-go/internal/math"
)

// ECMResult represents the result of ECM factoring
type ECMResult struct {
	Factor *math.BigInt
	Stage  int // 1 or 2
	Curve  uint64
	Error  error
}

// ECMFactor performs Elliptic Curve Method factoring
// Based on ecm() function from ecm.cpp
// Uses multiple curves to find factors
func ECMFactor(n *math.BigInt, B1, B2 uint64, numCurves uint32) (*ECMResult, error) {
	if n.IsZero() || n.Sign() < 0 {
		return nil, fmt.Errorf("n must be positive")
	}

	if B1 == 0 {
		return nil, fmt.Errorf("B1 must be greater than 0")
	}

	if numCurves == 0 {
		numCurves = 1
	}

	// Try multiple curves
	for curve := uint64(1); curve <= uint64(numCurves); curve++ {
		result, err := ecmSingleCurve(n, B1, B2, curve)
		if err != nil {
			continue
		}

		if result.Factor != nil {
			result.Curve = curve
			return result, nil
		}
	}

	// No factor found after all curves
	return &ECMResult{
		Factor: nil,
		Curve:  uint64(numCurves),
	}, nil
}

// ECMFactorFromParams performs ECM factoring from parameters
func ECMFactorFromParams(k float64, b, n uint64, c int64, B1, B2 uint64, numCurves uint32) (*ECMResult, error) {
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

	if numCurves == 0 {
		numCurves = 1
	}

	return ECMFactor(nBig, B1, B2, numCurves)
}

// ecmSingleCurve performs ECM on a single curve
// This is a simplified implementation - the full version uses elliptic curve arithmetic
func ecmSingleCurve(n *math.BigInt, B1, B2 uint64, sigma uint64) (*ECMResult, error) {
	// Stage 1: Compute point multiplication on elliptic curve
	// Simplified: we'll use a basic approach
	stage1Result, err := ecmStage1(n, B1, sigma)
	if err != nil {
		return nil, err
	}

	// Check for factor after stage 1
	// In ECM, we check GCD of the point coordinates with N
	factor := math.GCD(stage1Result, n)
	
	if !factor.IsOne() && factor.Cmp(n.Int) != 0 {
		return &ECMResult{
			Factor: factor,
			Stage:  1,
			Curve:  sigma,
		}, nil
	}

	// Stage 2: Process additional primes
	if B2 > B1 {
		stage2Result, err := ecmStage2(n, stage1Result, B1, B2, sigma)
		if err != nil {
			return nil, err
		}

		// Check for factor after stage 2
		factor = math.GCD(stage2Result, n)
		
		if !factor.IsOne() && factor.Cmp(n.Int) != 0 {
			return &ECMResult{
				Factor: factor,
				Stage:  2,
				Curve:  sigma,
			}, nil
		}
	}

	// No factor found on this curve
	return &ECMResult{
		Factor: nil,
		Stage:  2,
		Curve:  sigma,
	}, nil
}

// ecmStage1 performs ECM stage 1
// This is a simplified version - full ECM uses elliptic curve point multiplication
func ecmStage1(n *math.BigInt, B1 uint64, sigma uint64) (*math.BigInt, error) {
	// In full ECM, this would compute a point on an elliptic curve
	// For simplicity, we'll use a similar approach to P-1
	
	// Use sigma as a seed for the curve
	start := math.NewBigIntFromUint64(sigma)
	if start.Cmp(n.Int) >= 0 {
		start = math.NewBigInt(3)
	}

	// Compute E = product of all prime powers <= B1
	E := computePrimePowerProduct(B1)

	// Compute start^E mod N (simplified - in real ECM this is point multiplication)
	result := math.ModPow(start, E, n)

	return result, nil
}

// ecmStage2 performs ECM stage 2
func ecmStage2(n *math.BigInt, stage1Result *math.BigInt, B1, B2 uint64, sigma uint64) (*math.BigInt, error) {
	// Simplified version - full ECM uses more complex pairing operations
	result := stage1Result.Clone()

	// Process primes between B1 and B2
	for p := B1 + 1; p <= B2 && p <= B1+1000; p++ {
		if isPrimeSimple(p) {
			pBig := math.NewBigIntFromUint64(p)
			result = math.ModPow(result, pBig, n)
		}
	}

	return result, nil
}

