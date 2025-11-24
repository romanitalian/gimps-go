package algorithms

import (
	"math/big"

	"github.com/romanitalian/gimps-go/internal/math"
)

// computePrimePowerProduct computes the product of all prime powers <= bound
func computePrimePowerProduct(bound uint64) *math.BigInt {
	result := math.NewBigInt(1)

	for p := uint64(2); p <= bound; p++ {
		if isPrimeSimple(p) {
			// Include the highest power of p <= bound
			power := p
			for power*p <= bound {
				power *= p
			}
			powerBig := math.NewBigIntFromUint64(power)
			result = &math.BigInt{Int: new(big.Int).Mul(result.Int, powerBig.Int)}
		}
	}

	return result
}

// isPrimeSimple checks if a number is prime (simple trial division)
func isPrimeSimple(n uint64) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	for i := uint64(3); i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// estimateB1 estimates a reasonable B1 bound based on the number size
func estimateB1(n *math.BigInt) uint64 {
	bits := n.BitLen()
	
	// Rough estimation based on number of bits
	if bits < 100 {
		return 1000
	} else if bits < 200 {
		return 5000
	} else if bits < 500 {
		return 50000
	} else if bits < 1000 {
		return 500000
	} else {
		return 1000000
	}
}

// constructNumberFromParams constructs the number k*b^n+c from parameters
func constructNumberFromParams(k float64, b, n uint64, c int64) *math.BigInt {
	// b^n
	bBig := math.NewBigIntFromUint64(b)
	nBig := math.NewBigIntFromUint64(n)
	bToN := math.ModPow(bBig, nBig, nil) // No modulus for intermediate calculation
	
	// k * b^n
	if k != 1.0 {
		kBig := math.NewBigInt(int64(k))
		bToN = &math.BigInt{Int: new(big.Int).Mul(bToN.Int, kBig.Int)}
	}
	
	// k * b^n + c
	if c != 0 {
		cBig := math.NewBigInt(c)
		bToN = &math.BigInt{Int: new(big.Int).Add(bToN.Int, cBig.Int)}
	}
	
	return bToN
}

