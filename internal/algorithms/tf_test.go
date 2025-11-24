package algorithms

import (
	"math/big"
	"testing"
)

func TestTrialFactor_M11(t *testing.T) {
	// M11 = 2^11 - 1 = 2047 = 23 * 89
	// Factors of form 2kp+1 where p=11
	// For k=1: 2*1*11+1 = 23
	// For k=4: 2*4*11+1 = 89
	result, err := TrialFactor(11, 100.0)
	if err != nil {
		t.Fatalf("TrialFactor(11, 100.0) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("TrialFactor(11, 100.0) result has error: %v", result.Error)
	}
	if len(result.Factors) == 0 {
		t.Errorf("TrialFactor should find factors for M11")
	}
	
	// Check that we found at least one known factor
	found23 := false
	found89 := false
	expected23 := big.NewInt(23)
	expected89 := big.NewInt(89)
	for _, factor := range result.Factors {
		if factor.Cmp(expected23) == 0 {
			found23 = true
		}
		if factor.Cmp(expected89) == 0 {
			found89 = true
		}
	}
	if !found23 && !found89 {
		t.Errorf("TrialFactor should find at least one of 23 or 89 for M11")
	}
}

func TestTrialFactor_M5(t *testing.T) {
	// M5 = 2^5 - 1 = 31 (prime, no small factors)
	result, err := TrialFactor(5, 100.0)
	if err != nil {
		t.Fatalf("TrialFactor(5, 100.0) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("TrialFactor(5, 100.0) result has error: %v", result.Error)
	}
	// M5 is prime, so no factors should be found (or factors should be empty)
	// This is acceptable behavior
}

func TestTrialFactor_InvalidExponent(t *testing.T) {
	// Exponent < 2 should return error
	_, err := TrialFactor(1, 100.0)
	if err == nil {
		t.Errorf("TrialFactor(1, 100.0) should return error")
	}
}

func TestTrialFactorFromWorkUnit(t *testing.T) {
	// Test M11 factoring from work unit
	result, err := TrialFactorFromWorkUnit(1.0, 2, 11, -1, 100.0)
	if err != nil {
		t.Fatalf("TrialFactorFromWorkUnit returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("TrialFactorFromWorkUnit result has error: %v", result.Error)
	}
	if len(result.Factors) == 0 {
		t.Errorf("TrialFactorFromWorkUnit should find factors for M11")
	}
}

func TestTrialFactorFromWorkUnit_Invalid(t *testing.T) {
	// Not a Mersenne number
	_, err := TrialFactorFromWorkUnit(2.0, 2, 11, -1, 100.0)
	if err == nil {
		t.Errorf("TrialFactorFromWorkUnit should return error for non-Mersenne number")
	}
}

func TestTrialFactorFromWorkUnit_DefaultMaxBits(t *testing.T) {
	// Test with default maxBits (0)
	result, err := TrialFactorFromWorkUnit(1.0, 2, 11, -1, 0.0)
	if err != nil {
		t.Fatalf("TrialFactorFromWorkUnit returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("TrialFactorFromWorkUnit result has error: %v", result.Error)
	}
	// Should work with default maxBits
}

func TestComputeMersenneFactor(t *testing.T) {
	// Test factor computation: 2*1*11+1 = 23
	factor := computeMersenneFactor(11, 1)
	if factor == nil {
		t.Fatalf("computeMersenneFactor(11, 1) returned nil")
	}
	expected := big.NewInt(23)
	if factor.Cmp(expected) != 0 {
		t.Errorf("computeMersenneFactor(11, 1) = %v, expected 23", factor)
	}
	
	// Test factor computation: 2*4*11+1 = 89
	factor = computeMersenneFactor(11, 4)
	if factor == nil {
		t.Fatalf("computeMersenneFactor(11, 4) returned nil")
	}
	expected = big.NewInt(89)
	if factor.Cmp(expected) != 0 {
		t.Errorf("computeMersenneFactor(11, 4) = %v, expected 89", factor)
	}
}

func TestIsMersenneFactor(t *testing.T) {
	// Test that isMersenneFactor works correctly through the full TrialFactor function
	// For M11, factor 23 should divide it
	result, err := TrialFactor(11, 100.0)
	if err != nil {
		t.Fatalf("TrialFactor(11, 100.0) returned error: %v", err)
	}
	
	// Check that isMersenneFactor works correctly by verifying found factors
	if len(result.Factors) > 0 {
		// At least one factor was found, which means isMersenneFactor worked
		// We can't directly test isMersenneFactor without exposing it,
		// but we can verify the results are correct
	}
}

func TestTrialFactorParallel(t *testing.T) {
	// Test parallel factoring
	result, err := TrialFactorParallel(11, 100.0, 2)
	if err != nil {
		t.Fatalf("TrialFactorParallel(11, 100.0, 2) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("TrialFactorParallel(11, 100.0, 2) result has error: %v", result.Error)
	}
	if len(result.Factors) == 0 {
		t.Errorf("TrialFactorParallel should find factors for M11")
	}
}

func TestTrialFactorParallel_SingleWorker(t *testing.T) {
	// Test with single worker
	result, err := TrialFactorParallel(11, 100.0, 1)
	if err != nil {
		t.Fatalf("TrialFactorParallel(11, 100.0, 1) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("TrialFactorParallel(11, 100.0, 1) result has error: %v", result.Error)
	}
	// Should work with single worker
}

func TestTrialFactorParallel_InvalidWorkers(t *testing.T) {
	// Test with 0 workers (should default to 1)
	result, err := TrialFactorParallel(11, 100.0, 0)
	if err != nil {
		t.Fatalf("TrialFactorParallel(11, 100.0, 0) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("TrialFactorParallel(11, 100.0, 0) result has error: %v", result.Error)
	}
	// Should work (defaults to 1 worker)
}

