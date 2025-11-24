package algorithms

import (
	"math/big"
	"testing"

	"github.com/romanitalian/gimps-go/internal/math"
)

func TestLucasLehmerTest_M3(t *testing.T) {
	// M3 = 2^3 - 1 = 7 (prime)
	result, err := LucasLehmerTest(3)
	if err != nil {
		t.Fatalf("LucasLehmerTest(3) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("LucasLehmerTest(3) result has error: %v", result.Error)
	}
	if !result.IsPrime {
		t.Errorf("M3 should be prime, got IsPrime=%v", result.IsPrime)
	}
}

func TestLucasLehmerTest_M5(t *testing.T) {
	// M5 = 2^5 - 1 = 31 (prime)
	result, err := LucasLehmerTest(5)
	if err != nil {
		t.Fatalf("LucasLehmerTest(5) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("LucasLehmerTest(5) result has error: %v", result.Error)
	}
	if !result.IsPrime {
		t.Errorf("M5 should be prime, got IsPrime=%v", result.IsPrime)
	}
}

func TestLucasLehmerTest_M11(t *testing.T) {
	// M11 = 2^11 - 1 = 2047 = 23 * 89 (composite)
	result, err := LucasLehmerTest(11)
	if err != nil {
		t.Fatalf("LucasLehmerTest(11) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("LucasLehmerTest(11) result has error: %v", result.Error)
	}
	if result.IsPrime {
		t.Errorf("M11 should be composite, got IsPrime=%v", result.IsPrime)
	}
}

func TestLucasLehmerTest_M13(t *testing.T) {
	// M13 = 2^13 - 1 = 8191 (prime)
	result, err := LucasLehmerTest(13)
	if err != nil {
		t.Fatalf("LucasLehmerTest(13) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("LucasLehmerTest(13) result has error: %v", result.Error)
	}
	if !result.IsPrime {
		t.Errorf("M13 should be prime, got IsPrime=%v", result.IsPrime)
	}
}

func TestLucasLehmerTest_M17(t *testing.T) {
	// M17 = 2^17 - 1 = 131071 (prime)
	result, err := LucasLehmerTest(17)
	if err != nil {
		t.Fatalf("LucasLehmerTest(17) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("LucasLehmerTest(17) result has error: %v", result.Error)
	}
	if !result.IsPrime {
		t.Errorf("M17 should be prime, got IsPrime=%v", result.IsPrime)
	}
}

func TestLucasLehmerTest_InvalidExponent(t *testing.T) {
	// Exponent < 2 should return error
	result, err := LucasLehmerTest(1)
	if err != nil {
		t.Fatalf("LucasLehmerTest(1) should not return error, got: %v", err)
	}
	if result.Error == nil {
		t.Errorf("LucasLehmerTest(1) should have error in result")
	}
	if result.IsPrime {
		t.Errorf("LucasLehmerTest(1) should return IsPrime=false")
	}
}

func TestLucasLehmerTestFromWorkUnit(t *testing.T) {
	// Test M5 = 2^5 - 1 = 31 (prime)
	result, err := LucasLehmerTestFromWorkUnit(1.0, 2, 5, -1)
	if err != nil {
		t.Fatalf("LucasLehmerTestFromWorkUnit returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("LucasLehmerTestFromWorkUnit result has error: %v", result.Error)
	}
	if !result.IsPrime {
		t.Errorf("M5 should be prime, got IsPrime=%v", result.IsPrime)
	}
}

func TestLucasLehmerTestFromWorkUnit_Invalid(t *testing.T) {
	// Not a Mersenne number
	_, err := LucasLehmerTestFromWorkUnit(2.0, 2, 5, -1)
	if err == nil {
		t.Errorf("LucasLehmerTestFromWorkUnit should return error for non-Mersenne number")
	}
}

func TestLucasLehmerTestWithProgress(t *testing.T) {
	// Test M5 with progress callback
	iterations := make([]uint64, 0)
	result, err := LucasLehmerTestWithProgress(5, func(iter, total uint64, residue *math.BigInt) error {
		iterations = append(iterations, iter)
		return nil
	})
	if err != nil {
		t.Fatalf("LucasLehmerTestWithProgress(5) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("LucasLehmerTestWithProgress(5) result has error: %v", result.Error)
	}
	if !result.IsPrime {
		t.Errorf("M5 should be prime, got IsPrime=%v", result.IsPrime)
	}
	if len(iterations) == 0 {
		t.Errorf("Progress callback should be called")
	}
}

func TestKnownMersennePrimes(t *testing.T) {
	known := KnownMersennePrimes()
	if len(known) == 0 {
		t.Errorf("KnownMersennePrimes should return non-empty list")
	}

	// Test first few known primes
	expected := []uint64{2, 3, 5, 7, 13, 17, 19, 31}
	for _, exp := range expected {
		found := false
		for _, p := range known {
			if p == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("KnownMersennePrimes should include %d", exp)
		}
	}
}

func TestIsKnownMersennePrime(t *testing.T) {
	if !IsKnownMersennePrime(3) {
		t.Errorf("IsKnownMersennePrime(3) should return true")
	}
	if !IsKnownMersennePrime(5) {
		t.Errorf("IsKnownMersennePrime(5) should return true")
	}
	if IsKnownMersennePrime(4) {
		t.Errorf("IsKnownMersennePrime(4) should return false")
	}
	if IsKnownMersennePrime(11) {
		t.Errorf("IsKnownMersennePrime(11) should return false (composite)")
	}
}

func TestExtractResidue64(t *testing.T) {
	// Test residue extraction
	value := &struct {
		Int *big.Int
	}{
		Int: big.NewInt(0x1234567890ABCDEF),
	}
	
	// Create a mock BigInt-like structure for testing
	// Since we can't directly test extractResidue64, we test through LucasLehmerTest
	result, err := LucasLehmerTest(5)
	if err != nil {
		t.Fatalf("LucasLehmerTest(5) returned error: %v", err)
	}
	if result.Residue == nil {
		t.Errorf("Residue should not be nil")
	}
	
	// Residue should be a 64-bit value
	if result.Residue.Int.BitLen() > 64 {
		t.Errorf("Residue should be 64-bit or less, got %d bits", result.Residue.Int.BitLen())
	}
	
	_ = value // Suppress unused variable warning
}

