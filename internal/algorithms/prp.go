package algorithms

import (
	"fmt"
	"math/big"

	"github.com/romanitalian/gimps-go/internal/math"
)

// PRPResidueType represents the type of PRP residue
type PRPResidueType int32

const (
	PRPResidueTypeStandard PRPResidueType = 1
	PRPResidueTypeFermat   PRPResidueType = 2
	PRPResidueTypeFermatVar PRPResidueType = 3
	PRPResidueTypeCofactor PRPResidueType = 5
)

// PRPResult represents the result of a PRP test
type PRPResult struct {
	IsProbablePrime bool
	Residue         *math.BigInt
	ResidueType     PRPResidueType
	Error           error
}

// PRPTest performs a Probable Prime test
// Based on prp() function from commonb.c
// Computes base^(N-1) mod N. If result == 1, then N is probably prime
func PRPTest(n *math.BigInt, base int) (*PRPResult, error) {
	if n.IsZero() || n.Sign() < 0 {
		return &PRPResult{Error: fmt.Errorf("n must be positive")}, nil
	}

	// Check if n is even (except 2)
	if n.IsEven() && n.Cmp(big.NewInt(2)) != 0 {
		return &PRPResult{IsProbablePrime: false}, nil
	}

	// Check if n is 1
	if n.IsOne() {
		return &PRPResult{IsProbablePrime: false}, nil
	}

	// Compute base^(N-1) mod N
	baseBig := math.NewBigInt(int64(base))
	nMinus1 := &math.BigInt{Int: new(big.Int).Sub(n.Int, big.NewInt(1))}
	
	result := math.ModPow(baseBig, nMinus1, n)

	// If result == 1, then N is probably prime
	isProbablePrime := result.IsOne()

	return &PRPResult{
		IsProbablePrime: isProbablePrime,
		Residue:         result,
		ResidueType:     PRPResidueTypeStandard,
	}, nil
}

// PRPTestMersenne performs PRP test optimized for Mersenne numbers
// For Mersenne number M_p = 2^p - 1, we compute 3^(M_p-1) mod M_p
// Based on PRP test optimizations for Mersenne numbers
func PRPTestMersenne(exponent uint64, base int) (*PRPResult, error) {
	// Compute Mersenne number: M_p = 2^p - 1
	mp := math.MersenneNumber(exponent)

	// For Mersenne numbers, we can optimize the exponentiation
	// M_p - 1 = 2^p - 2 = 2 * (2^(p-1) - 1)
	// But for simplicity, we'll use the standard method

	return PRPTest(mp, base)
}

// PRPTestFromWorkUnit performs PRP test from work unit parameters
func PRPTestFromWorkUnit(k float64, b, n uint64, c int64, prpBase uint32) (*PRPResult, error) {
	// Construct the number k*b^n+c
	nBig := constructNumberFromParams(k, b, n, c)
	if nBig == nil {
		return nil, fmt.Errorf("failed to construct number from work unit")
	}

	base := int(prpBase)
	if base == 0 {
		base = 3 // Default base for PRP tests
	}

	return PRPTest(nBig, base)
}

// PRPTestWithProgress performs PRP test with progress callback
func PRPTestWithProgress(n *math.BigInt, base int, callback func(uint64, uint64) error) (*PRPResult, error) {
	if n.IsZero() || n.Sign() < 0 {
		return &PRPResult{Error: fmt.Errorf("n must be positive")}, nil
	}

	// Check if n is even (except 2)
	if n.IsEven() && n.Cmp(big.NewInt(2)) != 0 {
		return &PRPResult{IsProbablePrime: false}, nil
	}

	// Check if n is 1
	if n.IsOne() {
		return &PRPResult{IsProbablePrime: false}, nil
	}

	// Compute base^(N-1) mod N using binary exponentiation with progress
	baseBig := math.NewBigInt(int64(base))
	nMinus1 := &math.BigInt{Int: new(big.Int).Sub(n.Int, big.NewInt(1))}
	
	// Get bit length for progress tracking
	bitLen := nMinus1.BitLen()
	
	result := math.NewBigInt(1)
	power := baseBig.Clone()
	
	for i := 0; i < bitLen; i++ {
		if nMinus1.Bit(i) == 1 {
			result = math.ModMul(result, power, n)
		}
		power = math.ModSqr(power, n)
		
		if callback != nil {
			if err := callback(uint64(i+1), uint64(bitLen)); err != nil {
				return nil, err
			}
		}
	}

	isProbablePrime := result.IsOne()

	return &PRPResult{
		IsProbablePrime: isProbablePrime,
		Residue:         result,
		ResidueType:     PRPResidueTypeStandard,
	}, nil
}


