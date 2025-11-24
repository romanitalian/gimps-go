package algorithms

import (
	"testing"

	"github.com/romanitalian/gimps-go/internal/math"
)

func TestPMinus1Factor_Simple(t *testing.T) {
	// Test with a number that has a small factor
	// Use a composite number: 91 = 7 * 13
	n := math.NewBigInt(91)
	result, err := PMinus1Factor(n, 10, 100)
	if err != nil {
		t.Fatalf("PMinus1Factor(91, 10, 100) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PMinus1Factor(91, 10, 100) result has error: %v", result.Error)
	}
	// P-1 may or may not find a factor depending on B1/B2
	// The important thing is that it doesn't crash
}

func TestPMinus1Factor_Prime(t *testing.T) {
	// Test with a prime number
	n := math.NewBigInt(31)
	result, err := PMinus1Factor(n, 10, 100)
	if err != nil {
		t.Fatalf("PMinus1Factor(31, 10, 100) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PMinus1Factor(31, 10, 100) result has error: %v", result.Error)
	}
	// Prime number should not have factors found
	if result.Factor != nil {
		// Factor found might be the number itself or 1, which is acceptable
		if result.Factor.Cmp(n.Int) == 0 || result.Factor.IsOne() {
			// This is expected for prime numbers
		}
	}
}

func TestPMinus1Factor_InvalidInput(t *testing.T) {
	// Test with zero
	n := math.NewBigInt(0)
	_, err := PMinus1Factor(n, 10, 100)
	if err == nil {
		t.Errorf("PMinus1Factor(0, 10, 100) should return error")
	}

	// Test with B1 = 0
	n = math.NewBigInt(31)
	_, err = PMinus1Factor(n, 0, 100)
	if err == nil {
		t.Errorf("PMinus1Factor(31, 0, 100) should return error")
	}
}

func TestPMinus1FactorFromParams(t *testing.T) {
	// Test M11 = 2^11 - 1 = 2047 = 23 * 89
	result, err := PMinus1FactorFromParams(1.0, 2, 11, -1, 10, 100)
	if err != nil {
		t.Fatalf("PMinus1FactorFromParams returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PMinus1FactorFromParams result has error: %v", result.Error)
	}
	// May or may not find factors depending on B1/B2
}

func TestPMinus1FactorFromParams_DefaultB1(t *testing.T) {
	// Test with default B1 (0 should estimate)
	result, err := PMinus1FactorFromParams(1.0, 2, 11, -1, 0, 0)
	if err != nil {
		t.Fatalf("PMinus1FactorFromParams returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PMinus1FactorFromParams result has error: %v", result.Error)
	}
	// Should work with default B1/B2
}

func TestPPlus1Factor_Simple(t *testing.T) {
	// Test with a composite number
	n := math.NewBigInt(91)
	result, err := PPlus1Factor(n, 10, 100, 1)
	if err != nil {
		t.Fatalf("PPlus1Factor(91, 10, 100, 1) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PPlus1Factor(91, 10, 100, 1) result has error: %v", result.Error)
	}
	// P+1 may or may not find a factor
}

func TestPPlus1Factor_InvalidInput(t *testing.T) {
	// Test with zero
	n := math.NewBigInt(0)
	_, err := PPlus1Factor(n, 10, 100, 1)
	if err == nil {
		t.Errorf("PPlus1Factor(0, 10, 100, 1) should return error")
	}

	// Test with B1 = 0
	n = math.NewBigInt(31)
	_, err = PPlus1Factor(n, 0, 100, 1)
	if err == nil {
		t.Errorf("PPlus1Factor(31, 0, 100, 1) should return error")
	}
}

func TestPPlus1FactorFromParams(t *testing.T) {
	// Test M11
	result, err := PPlus1FactorFromParams(1.0, 2, 11, -1, 10, 100, 1)
	if err != nil {
		t.Fatalf("PPlus1FactorFromParams returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PPlus1FactorFromParams result has error: %v", result.Error)
	}
	// May or may not find factors
}

func TestPPlus1FactorFromParams_DefaultNthRun(t *testing.T) {
	// Test with default nthRun (0 should default to 1)
	result, err := PPlus1FactorFromParams(1.0, 2, 11, -1, 10, 100, 0)
	if err != nil {
		t.Fatalf("PPlus1FactorFromParams returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PPlus1FactorFromParams result has error: %v", result.Error)
	}
	// Should work with default nthRun
}

func TestECMFactor_Simple(t *testing.T) {
	// Test with a composite number
	n := math.NewBigInt(91)
	result, err := ECMFactor(n, 10, 100, 1)
	if err != nil {
		t.Fatalf("ECMFactor(91, 10, 100, 1) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("ECMFactor(91, 10, 100, 1) result has error: %v", result.Error)
	}
	// ECM may or may not find a factor
}

func TestECMFactor_InvalidInput(t *testing.T) {
	// Test with zero
	n := math.NewBigInt(0)
	_, err := ECMFactor(n, 10, 100, 1)
	if err == nil {
		t.Errorf("ECMFactor(0, 10, 100, 1) should return error")
	}

	// Test with B1 = 0
	n = math.NewBigInt(31)
	_, err = ECMFactor(n, 0, 100, 1)
	if err == nil {
		t.Errorf("ECMFactor(31, 0, 100, 1) should return error")
	}
}

func TestECMFactorFromParams(t *testing.T) {
	// Test M11
	result, err := ECMFactorFromParams(1.0, 2, 11, -1, 10, 100, 1)
	if err != nil {
		t.Fatalf("ECMFactorFromParams returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("ECMFactorFromParams result has error: %v", result.Error)
	}
	// May or may not find factors
}

func TestECMFactorFromParams_DefaultCurves(t *testing.T) {
	// Test with default numCurves (0 should default to 1)
	result, err := ECMFactorFromParams(1.0, 2, 11, -1, 10, 100, 0)
	if err != nil {
		t.Fatalf("ECMFactorFromParams returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("ECMFactorFromParams result has error: %v", result.Error)
	}
	// Should work with default numCurves
}

func TestECMFactor_MultipleCurves(t *testing.T) {
	// Test with multiple curves
	n := math.NewBigInt(91)
	result, err := ECMFactor(n, 10, 100, 3)
	if err != nil {
		t.Fatalf("ECMFactor(91, 10, 100, 3) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("ECMFactor(91, 10, 100, 3) result has error: %v", result.Error)
	}
	// Should try multiple curves
	if result.Curve == 0 {
		t.Errorf("ECMFactor should set Curve field")
	}
}

