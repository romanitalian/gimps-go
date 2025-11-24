package math

import (
	"math/big"
)

// BigInt wraps big.Int with additional utility methods
type BigInt struct {
	*big.Int
}

// NewBigInt creates a new BigInt from int64
func NewBigInt(x int64) *BigInt {
	return &BigInt{Int: big.NewInt(x)}
}

// NewBigIntFromUint64 creates a new BigInt from uint64
func NewBigIntFromUint64(x uint64) *BigInt {
	bi := &BigInt{Int: new(big.Int)}
	bi.SetUint64(x)
	return bi
}

// NewBigIntFromString creates a new BigInt from string
func NewBigIntFromString(s string, base int) (*BigInt, error) {
	bi := &BigInt{Int: new(big.Int)}
	_, ok := bi.SetString(s, base)
	if !ok {
		return nil, &ParseError{Value: s}
	}
	return bi, nil
}

// Clone creates a copy of the BigInt
func (bi *BigInt) Clone() *BigInt {
	return &BigInt{Int: new(big.Int).Set(bi.Int)}
}

// ModPow computes (base^exp) mod m
// This is a wrapper around big.Int.Exp for modular exponentiation
func ModPow(base, exp, m *BigInt) *BigInt {
	result := &BigInt{Int: new(big.Int)}
	result.Exp(base.Int, exp.Int, m.Int)
	return result
}

// ModPowUint64 computes (base^exp) mod m where exp is uint64
func ModPowUint64(base *BigInt, exp uint64, m *BigInt) *BigInt {
	expBig := NewBigIntFromUint64(exp)
	return ModPow(base, expBig, m)
}

// ModMul computes (a * b) mod m
func ModMul(a, b, m *BigInt) *BigInt {
	result := &BigInt{Int: new(big.Int)}
	result.Mul(a.Int, b.Int)
	result.Mod(result.Int, m.Int)
	return result
}

// ModAdd computes (a + b) mod m
func ModAdd(a, b, m *BigInt) *BigInt {
	result := &BigInt{Int: new(big.Int)}
	result.Add(a.Int, b.Int)
	result.Mod(result.Int, m.Int)
	return result
}

// ModSub computes (a - b) mod m
func ModSub(a, b, m *BigInt) *BigInt {
	result := &BigInt{Int: new(big.Int)}
	result.Sub(a.Int, b.Int)
	result.Mod(result.Int, m.Int)
	return result
}

// ModSqr computes (a^2) mod m
func ModSqr(a, m *BigInt) *BigInt {
	result := &BigInt{Int: new(big.Int)}
	result.Mul(a.Int, a.Int)
	result.Mod(result.Int, m.Int)
	return result
}

// IsZero returns true if the BigInt is zero
func (bi *BigInt) IsZero() bool {
	return bi.Sign() == 0
}

// IsOne returns true if the BigInt is one
func (bi *BigInt) IsOne() bool {
	return bi.Cmp(big.NewInt(1)) == 0
}

// IsEven returns true if the BigInt is even
func (bi *BigInt) IsEven() bool {
	return new(big.Int).And(bi.Int, big.NewInt(1)).Sign() == 0
}

// BitLen returns the length of the absolute value of x in bits
func (bi *BigInt) BitLen() int {
	return bi.Int.BitLen()
}

// GetUint64 returns the uint64 representation of the BigInt
// Returns 0 and false if the value cannot be represented as uint64
func (bi *BigInt) GetUint64() (uint64, bool) {
	if !bi.IsUint64() {
		return 0, false
	}
	return bi.Uint64(), true
}

// MersenneNumber computes 2^p - 1 for given exponent p
func MersenneNumber(p uint64) *BigInt {
	// 2^p
	two := big.NewInt(2)
	exp := big.NewInt(int64(p))
	result := &BigInt{Int: new(big.Int)}
	result.Exp(two, exp, nil)
	// 2^p - 1
	result.Sub(result.Int, big.NewInt(1))
	return result
}

// GCD computes the greatest common divisor of a and b
func GCD(a, b *BigInt) *BigInt {
	result := &BigInt{Int: new(big.Int)}
	result.GCD(nil, nil, a.Int, b.Int)
	return result
}

// ParseError represents an error parsing a BigInt
type ParseError struct {
	Value string
}

func (e *ParseError) Error() string {
	return "failed to parse BigInt: " + e.Value
}
