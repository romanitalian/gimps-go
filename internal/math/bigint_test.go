package math

import (
	"math/big"
	"testing"
)

func TestNewBigInt(t *testing.T) {
	bi := NewBigInt(42)
	if bi == nil {
		t.Fatalf("NewBigInt(42) returned nil")
	}
	if bi.Cmp(big.NewInt(42)) != 0 {
		t.Errorf("NewBigInt(42) = %v, expected 42", bi)
	}
}

func TestNewBigIntFromUint64(t *testing.T) {
	bi := NewBigIntFromUint64(12345)
	if bi == nil {
		t.Fatalf("NewBigIntFromUint64(12345) returned nil")
	}
	if bi.Cmp(big.NewInt(12345)) != 0 {
		t.Errorf("NewBigIntFromUint64(12345) = %v, expected 12345", bi)
	}
}

func TestNewBigIntFromString(t *testing.T) {
	bi, err := NewBigIntFromString("12345", 10)
	if err != nil {
		t.Fatalf("NewBigIntFromString(\"12345\", 10) returned error: %v", err)
	}
	if bi.Cmp(big.NewInt(12345)) != 0 {
		t.Errorf("NewBigIntFromString(\"12345\", 10) = %v, expected 12345", bi)
	}
}

func TestNewBigIntFromString_Invalid(t *testing.T) {
	_, err := NewBigIntFromString("invalid", 10)
	if err == nil {
		t.Errorf("NewBigIntFromString(\"invalid\", 10) should return error")
	}
}

func TestClone(t *testing.T) {
	bi1 := NewBigInt(42)
	bi2 := bi1.Clone()
	if bi2 == nil {
		t.Fatalf("Clone() returned nil")
	}
	if bi2.Cmp(bi1.Int) != 0 {
		t.Errorf("Clone() = %v, expected %v", bi2, bi1)
	}
	// Modify original, clone should be unaffected
	bi1.Add(bi1.Int, big.NewInt(1))
	if bi2.Cmp(big.NewInt(42)) != 0 {
		t.Errorf("Clone should be independent, got %v", bi2)
	}
}

func TestModPow(t *testing.T) {
	base := NewBigInt(3)
	exp := NewBigInt(5)
	mod := NewBigInt(7)
	result := ModPow(base, exp, mod)
	// 3^5 mod 7 = 243 mod 7 = 5
	expected := big.NewInt(5)
	if result.Cmp(expected) != 0 {
		t.Errorf("ModPow(3, 5, 7) = %v, expected 5", result)
	}
}

func TestModPowUint64(t *testing.T) {
	base := NewBigInt(3)
	exp := uint64(5)
	mod := NewBigInt(7)
	result := ModPowUint64(base, exp, mod)
	// 3^5 mod 7 = 243 mod 7 = 5
	expected := big.NewInt(5)
	if result.Cmp(expected) != 0 {
		t.Errorf("ModPowUint64(3, 5, 7) = %v, expected 5", result)
	}
}

func TestModMul(t *testing.T) {
	a := NewBigInt(5)
	b := NewBigInt(6)
	m := NewBigInt(7)
	result := ModMul(a, b, m)
	// 5 * 6 mod 7 = 30 mod 7 = 2
	expected := big.NewInt(2)
	if result.Cmp(expected) != 0 {
		t.Errorf("ModMul(5, 6, 7) = %v, expected 2", result)
	}
}

func TestModAdd(t *testing.T) {
	a := NewBigInt(5)
	b := NewBigInt(6)
	m := NewBigInt(7)
	result := ModAdd(a, b, m)
	// 5 + 6 mod 7 = 11 mod 7 = 4
	expected := big.NewInt(4)
	if result.Cmp(expected) != 0 {
		t.Errorf("ModAdd(5, 6, 7) = %v, expected 4", result)
	}
}

func TestModSub(t *testing.T) {
	a := NewBigInt(5)
	b := NewBigInt(6)
	m := NewBigInt(7)
	result := ModSub(a, b, m)
	// 5 - 6 mod 7 = -1 mod 7 = 6
	expected := big.NewInt(6)
	if result.Cmp(expected) != 0 {
		t.Errorf("ModSub(5, 6, 7) = %v, expected 6", result)
	}
}

func TestModSqr(t *testing.T) {
	a := NewBigInt(5)
	m := NewBigInt(7)
	result := ModSqr(a, m)
	// 5^2 mod 7 = 25 mod 7 = 4
	expected := big.NewInt(4)
	if result.Cmp(expected) != 0 {
		t.Errorf("ModSqr(5, 7) = %v, expected 4", result)
	}
}

func TestIsZero(t *testing.T) {
	bi := NewBigInt(0)
	if !bi.IsZero() {
		t.Errorf("IsZero() should return true for 0")
	}
	bi = NewBigInt(1)
	if bi.IsZero() {
		t.Errorf("IsZero() should return false for 1")
	}
}

func TestIsOne(t *testing.T) {
	bi := NewBigInt(1)
	if !bi.IsOne() {
		t.Errorf("IsOne() should return true for 1")
	}
	bi = NewBigInt(0)
	if bi.IsOne() {
		t.Errorf("IsOne() should return false for 0")
	}
}

func TestIsEven(t *testing.T) {
	bi := NewBigInt(2)
	if !bi.IsEven() {
		t.Errorf("IsEven() should return true for 2")
	}
	bi = NewBigInt(3)
	if bi.IsEven() {
		t.Errorf("IsEven() should return false for 3")
	}
}

func TestBitLen(t *testing.T) {
	bi := NewBigInt(8)
	// 8 = 1000 in binary, so 4 bits
	if bi.BitLen() != 4 {
		t.Errorf("BitLen(8) = %d, expected 4", bi.BitLen())
	}
}

func TestGetUint64(t *testing.T) {
	bi := NewBigIntFromUint64(12345)
	val, ok := bi.GetUint64()
	if !ok {
		t.Errorf("GetUint64() should return ok=true for 12345")
	}
	if val != 12345 {
		t.Errorf("GetUint64() = %d, expected 12345", val)
	}
}

func TestGetUint64_TooLarge(t *testing.T) {
	// Create a number larger than uint64 max
	bi, err := NewBigIntFromString("18446744073709551616", 10) // 2^64
	if err != nil {
		t.Fatalf("Failed to create large number: %v", err)
	}
	if bi == nil {
		t.Fatalf("NewBigIntFromString returned nil")
	}
	_, ok := bi.GetUint64()
	if ok {
		t.Errorf("GetUint64() should return ok=false for number > uint64 max")
	}
}

func TestMersenneNumber(t *testing.T) {
	// M3 = 2^3 - 1 = 7
	mp := MersenneNumber(3)
	if mp.Cmp(big.NewInt(7)) != 0 {
		t.Errorf("MersenneNumber(3) = %v, expected 7", mp)
	}
	
	// M5 = 2^5 - 1 = 31
	mp = MersenneNumber(5)
	if mp.Cmp(big.NewInt(31)) != 0 {
		t.Errorf("MersenneNumber(5) = %v, expected 31", mp)
	}
}

func TestGCD(t *testing.T) {
	a := NewBigInt(48)
	b := NewBigInt(18)
	result := GCD(a, b)
	// GCD(48, 18) = 6
	expected := big.NewInt(6)
	if result.Cmp(expected) != 0 {
		t.Errorf("GCD(48, 18) = %v, expected 6", result)
	}
}

func TestGCD_Coprime(t *testing.T) {
	a := NewBigInt(7)
	b := NewBigInt(5)
	result := GCD(a, b)
	// GCD(7, 5) = 1
	expected := big.NewInt(1)
	if result.Cmp(expected) != 0 {
		t.Errorf("GCD(7, 5) = %v, expected 1", result)
	}
}

func TestParseError(t *testing.T) {
	err := &ParseError{Value: "invalid"}
	if err.Error() == "" {
		t.Errorf("ParseError.Error() should return non-empty string")
	}
}

