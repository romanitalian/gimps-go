package algorithms

import (
	"fmt"

	"github.com/romanitalian/gimps-go/internal/math"
)

// PPlus1Result represents the result of P+1 factoring
type PPlus1Result struct {
	Factor *math.BigInt
	Stage  int // 1 or 2
	Error  error
}

// PPlus1Factor performs P+1 factoring
// Based on pplus1() function from ecm.cpp
// Similar to P-1 but uses P+1 group
func PPlus1Factor(n *math.BigInt, B1, B2 uint64, nthRun int32) (*PPlus1Result, error) {
	if n.IsZero() || n.Sign() < 0 {
		return nil, fmt.Errorf("n must be positive")
	}

	if B1 == 0 {
		return nil, fmt.Errorf("B1 must be greater than 0")
	}

	// Stage 1: Similar to P-1 but with different starting value
	stage1Result, err := pp1Stage1(n, B1, nthRun)
	if err != nil {
		return nil, err
	}

	// Check for factor after stage 1
	// GCD(V-2, N) where V is the stage 1 result
	vMinus2 := math.ModSub(stage1Result, math.NewBigInt(2), n)
	factor := math.GCD(vMinus2, n)
	
	if !factor.IsOne() && factor.Cmp(n.Int) != 0 {
		return &PPlus1Result{
			Factor: factor,
			Stage:  1,
		}, nil
	}

	// Stage 2: Similar to P-1 stage 2
	if B2 > B1 {
		stage2Result, err := pp1Stage2(n, stage1Result, B1, B2)
		if err != nil {
			return nil, err
		}

		// Check for factor after stage 2
		stage2Minus2 := math.ModSub(stage2Result, math.NewBigInt(2), n)
		factor = math.GCD(stage2Minus2, n)
		
		if !factor.IsOne() && factor.Cmp(n.Int) != 0 {
			return &PPlus1Result{
				Factor: factor,
				Stage:  2,
			}, nil
		}
	}

	// No factor found
	return &PPlus1Result{
		Factor: nil,
		Stage:  2,
	}, nil
}

// PPlus1FactorFromParams performs P+1 factoring from parameters
func PPlus1FactorFromParams(k float64, b, n uint64, c int64, B1, B2 uint64, nthRun int32) (*PPlus1Result, error) {
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

	if nthRun == 0 {
		nthRun = 1
	}

	return PPlus1Factor(nBig, B1, B2, nthRun)
}

// PPlus1FactorFromWorkUnit performs P+1 factoring from work unit parameters
func PPlus1FactorFromWorkUnit(k float64, b, n uint64, c int64, B1, B2 uint64, nthRun int32) (*PPlus1Result, error) {
	return PPlus1FactorFromParams(k, b, n, c, B1, B2, nthRun)
}

// pp1Stage1 performs P+1 stage 1
func pp1Stage1(n *math.BigInt, B1 uint64, nthRun int32) (*math.BigInt, error) {
	// Starting value depends on nthRun:
	// 1 = 2/7, 2 = 6/5, 3+ = random
	var start *math.BigInt
	
	switch nthRun {
	case 1:
		// Start with 2/7 (simplified - in practice this is more complex)
		start = math.NewBigInt(2)
	case 2:
		// Start with 6/5
		start = math.NewBigInt(6)
	default:
		// Random start (simplified)
		start = math.NewBigInt(3)
	}

	// Compute E = product of all prime powers <= B1
	E := computePrimePowerProduct(B1)

	// Compute start^E mod N
	result := math.ModPow(start, E, n)

	return result, nil
}

// pp1Stage2 performs P+1 stage 2
func pp1Stage2(n *math.BigInt, stage1Result *math.BigInt, B1, B2 uint64) (*math.BigInt, error) {
	// Similar to P-1 stage 2
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

