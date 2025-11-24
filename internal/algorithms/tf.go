package algorithms

import (
	"fmt"
	"math/big"

	"github.com/romanitalian/gimps-go/internal/math"
)

// TrialFactorResult represents the result of trial factoring
type TrialFactorResult struct {
	Factors []*math.BigInt
	Error   error
}

// TrialFactor performs trial factoring for a Mersenne number
// Based on primeFactor() function from commonb.c
// For Mersenne number M_p = 2^p - 1, factors are of the form 2kp+1
func TrialFactor(exponent uint64, maxBits float64) (*TrialFactorResult, error) {
	if exponent < 2 {
		return nil, fmt.Errorf("exponent must be at least 2")
	}

	// Compute Mersenne number: M_p = 2^p - 1
	mp := math.MersenneNumber(exponent)

	// For Mersenne numbers, factors are of the form 2kp+1 where k is a positive integer
	// We need to test candidates up to sqrt(M_p) or maxBits limit
	
	maxFactor := computeMaxFactor(mp, maxBits)
	
	var factors []*math.BigInt
	
	// Start with k=1, so factor = 2*1*p + 1 = 2p+1
	k := uint64(1)
	
	for {
		// Compute candidate factor: 2kp+1
		factor := computeMersenneFactor(exponent, k)
		
		// Check if factor exceeds maximum
		if factor.Cmp(maxFactor.Int) > 0 {
			break
		}
		
		// Check if factor divides M_p
		// We check: 2^p mod factor == 1
		if isMersenneFactor(mp, factor, exponent) {
			factors = append(factors, factor)
		}
		
		k++
		
		// Safety check to avoid infinite loops
		if k > 1000000 {
			break
		}
	}
	
	return &TrialFactorResult{
		Factors: factors,
	}, nil
}

// TrialFactorFromWorkUnit performs trial factoring from work unit parameters
func TrialFactorFromWorkUnit(k float64, b, n uint64, c int64, factorTo float64) (*TrialFactorResult, error) {
	// Check if it's a Mersenne number
	if k != 1.0 || b != 2 || c != -1 {
		return nil, fmt.Errorf("trial factoring currently only supports Mersenne numbers")
	}
	
	if n == 0 {
		return nil, fmt.Errorf("invalid Mersenne exponent")
	}
	
	maxBits := factorTo
	if maxBits == 0 {
		// Default: factor to reasonable limit
		maxBits = 70.0
	}
	
	return TrialFactor(n, maxBits)
}

// computeMersenneFactor computes a candidate factor 2kp+1 for Mersenne number
func computeMersenneFactor(exponent, k uint64) *math.BigInt {
	// factor = 2*k*exponent + 1
	two := math.NewBigInt(2)
	kBig := math.NewBigIntFromUint64(k)
	pBig := math.NewBigIntFromUint64(exponent)
	
	// 2 * k
	result := &math.BigInt{Int: new(big.Int).Mul(two.Int, kBig.Int)}
	// 2 * k * p
	result = &math.BigInt{Int: new(big.Int).Mul(result.Int, pBig.Int)}
	// 2 * k * p + 1
	result = &math.BigInt{Int: new(big.Int).Add(result.Int, big.NewInt(1))}
	
	return result
}

// isMersenneFactor checks if a candidate factor divides the Mersenne number
// For Mersenne number M_p = 2^p - 1, a factor f divides M_p if 2^p mod f == 1
func isMersenneFactor(mp, factor *math.BigInt, exponent uint64) bool {
	// Check: 2^exponent mod factor == 1
	two := math.NewBigInt(2)
	exp := math.NewBigIntFromUint64(exponent)
	
	result := math.ModPow(two, exp, factor)
	
	return result.IsOne()
}

// computeMaxFactor computes the maximum factor to test based on maxBits
func computeMaxFactor(mp *math.BigInt, maxBits float64) *math.BigInt {
	if maxBits <= 0 {
		// Default: sqrt(M_p)
		sqrt := &math.BigInt{Int: new(big.Int).Sqrt(mp.Int)}
		return sqrt
	}
	
	// Compute 2^maxBits
	maxFactor := &math.BigInt{Int: new(big.Int)}
	maxFactor.Exp(big.NewInt(2), big.NewInt(int64(maxBits)), nil)
	
	// Also check sqrt(M_p) as upper bound
	sqrt := &math.BigInt{Int: new(big.Int).Sqrt(mp.Int)}
	if maxFactor.Cmp(sqrt.Int) > 0 {
		return sqrt
	}
	
	return maxFactor
}

// TrialFactorParallel performs trial factoring in parallel using goroutines
func TrialFactorParallel(exponent uint64, maxBits float64, numWorkers int) (*TrialFactorResult, error) {
	if numWorkers < 1 {
		numWorkers = 1
	}
	
	// Compute Mersenne number
	mp := math.MersenneNumber(exponent)
	maxFactor := computeMaxFactor(mp, maxBits)
	
	// Estimate number of candidates
	// Factors are 2kp+1, so k ranges from 1 to approximately (maxFactor-1)/(2*exponent)
	maxK := uint64(0)
	if maxFactor.Cmp(big.NewInt(2*int64(exponent)+1)) > 0 {
		maxKBig := new(big.Int).Sub(maxFactor.Int, big.NewInt(1))
		maxKBig.Div(maxKBig, big.NewInt(2*int64(exponent)))
		if maxKBig.IsUint64() {
			maxK = maxKBig.Uint64()
		}
	}
	
	if maxK == 0 {
		return &TrialFactorResult{Factors: []*math.BigInt{}}, nil
	}
	
	// Divide work among workers
	candidatesPerWorker := maxK / uint64(numWorkers)
	if candidatesPerWorker == 0 {
		candidatesPerWorker = 1
	}
	
	type result struct {
		factors []*math.BigInt
		err     error
	}
	
	results := make(chan result, numWorkers)
	
	// Launch workers
	for i := 0; i < numWorkers; i++ {
		startK := uint64(i) * candidatesPerWorker + 1
		endK := startK + candidatesPerWorker
		if i == numWorkers-1 {
			endK = maxK + 1
		}
		
		go func(start, end uint64) {
			var factors []*math.BigInt
			for k := start; k < end; k++ {
				factor := computeMersenneFactor(exponent, k)
				if factor.Cmp(maxFactor.Int) > 0 {
					break
				}
				if isMersenneFactor(mp, factor, exponent) {
					factors = append(factors, factor)
				}
			}
			results <- result{factors: factors}
		}(startK, endK)
	}
	
	// Collect results
	var allFactors []*math.BigInt
	for i := 0; i < numWorkers; i++ {
		res := <-results
		if res.err != nil {
			return nil, res.err
		}
		allFactors = append(allFactors, res.factors...)
	}
	
	return &TrialFactorResult{
		Factors: allFactors,
	}, nil
}

