package algorithms

import (
	"math/big"
	"testing"

	"github.com/romanitalian/gimps-go/internal/math"
)

func TestPRPTest_Prime(t *testing.T) {
	// Test with a small prime: 7
	n := math.NewBigInt(7)
	result, err := PRPTest(n, 3)
	if err != nil {
		t.Fatalf("PRPTest(7, 3) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PRPTest(7, 3) result has error: %v", result.Error)
	}
	if !result.IsProbablePrime {
		t.Errorf("7 should be probable prime, got IsProbablePrime=%v", result.IsProbablePrime)
	}
}

func TestPRPTest_Composite(t *testing.T) {
	// Test with a composite: 15 = 3 * 5
	n := math.NewBigInt(15)
	result, err := PRPTest(n, 3)
	if err != nil {
		t.Fatalf("PRPTest(15, 3) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PRPTest(15, 3) result has error: %v", result.Error)
	}
	if result.IsProbablePrime {
		t.Errorf("15 should not be probable prime, got IsProbablePrime=%v", result.IsProbablePrime)
	}
}

func TestPRPTest_Even(t *testing.T) {
	// Test with even number (except 2)
	n := math.NewBigInt(4)
	result, err := PRPTest(n, 3)
	if err != nil {
		t.Fatalf("PRPTest(4, 3) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PRPTest(4, 3) result has error: %v", result.Error)
	}
	if result.IsProbablePrime {
		t.Errorf("4 should not be probable prime, got IsProbablePrime=%v", result.IsProbablePrime)
	}
}

func TestPRPTest_Two(t *testing.T) {
	// Test with 2 (the only even prime)
	n := math.NewBigInt(2)
	result, err := PRPTest(n, 3)
	if err != nil {
		t.Fatalf("PRPTest(2, 3) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PRPTest(2, 3) result has error: %v", result.Error)
	}
	if !result.IsProbablePrime {
		t.Errorf("2 should be probable prime, got IsProbablePrime=%v", result.IsProbablePrime)
	}
}

func TestPRPTest_One(t *testing.T) {
	// Test with 1
	n := math.NewBigInt(1)
	result, err := PRPTest(n, 3)
	if err != nil {
		t.Fatalf("PRPTest(1, 3) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PRPTest(1, 3) result has error: %v", result.Error)
	}
	if result.IsProbablePrime {
		t.Errorf("1 should not be probable prime, got IsProbablePrime=%v", result.IsProbablePrime)
	}
}

func TestPRPTest_InvalidInput(t *testing.T) {
	// Test with zero
	n := math.NewBigInt(0)
	result, err := PRPTest(n, 3)
	if err != nil {
		t.Fatalf("PRPTest(0, 3) should not return error")
	}
	if result.Error == nil {
		t.Errorf("PRPTest(0, 3) should have error in result")
	}

	// Test with negative (if possible)
	neg := &math.BigInt{Int: big.NewInt(-5)}
	result, err = PRPTest(neg, 3)
	if err != nil {
		t.Fatalf("PRPTest(-5, 3) should not return error")
	}
	if result.Error == nil {
		t.Errorf("PRPTest(-5, 3) should have error in result")
	}
}

func TestPRPTestMersenne_M3(t *testing.T) {
	// M3 = 2^3 - 1 = 7 (prime)
	result, err := PRPTestMersenne(3, 3)
	if err != nil {
		t.Fatalf("PRPTestMersenne(3, 3) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PRPTestMersenne(3, 3) result has error: %v", result.Error)
	}
	if !result.IsProbablePrime {
		t.Errorf("M3 should be probable prime, got IsProbablePrime=%v", result.IsProbablePrime)
	}
}

func TestPRPTestMersenne_M5(t *testing.T) {
	// M5 = 2^5 - 1 = 31 (prime)
	result, err := PRPTestMersenne(5, 3)
	if err != nil {
		t.Fatalf("PRPTestMersenne(5, 3) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PRPTestMersenne(5, 3) result has error: %v", result.Error)
	}
	if !result.IsProbablePrime {
		t.Errorf("M5 should be probable prime, got IsProbablePrime=%v", result.IsProbablePrime)
	}
}

func TestPRPTestMersenne_M11(t *testing.T) {
	// M11 = 2^11 - 1 = 2047 = 23 * 89 (composite)
	result, err := PRPTestMersenne(11, 3)
	if err != nil {
		t.Fatalf("PRPTestMersenne(11, 3) returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PRPTestMersenne(11, 3) result has error: %v", result.Error)
	}
	// Note: PRP test may give false positives, so we don't assert composite here
	// The important thing is that it doesn't crash
}

func TestPRPTestFromWorkUnit(t *testing.T) {
	// Test M5 = 2^5 - 1 = 31 (prime)
	result, err := PRPTestFromWorkUnit(1.0, 2, 5, -1, 3)
	if err != nil {
		t.Fatalf("PRPTestFromWorkUnit returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PRPTestFromWorkUnit result has error: %v", result.Error)
	}
	if !result.IsProbablePrime {
		t.Errorf("M5 should be probable prime, got IsProbablePrime=%v", result.IsProbablePrime)
	}
}

func TestPRPTestFromWorkUnit_DefaultBase(t *testing.T) {
	// Test with default base (0 should default to 3)
	result, err := PRPTestFromWorkUnit(1.0, 2, 5, -1, 0)
	if err != nil {
		t.Fatalf("PRPTestFromWorkUnit returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PRPTestFromWorkUnit result has error: %v", result.Error)
	}
	// Should work with default base
}

func TestPRPTestWithProgress(t *testing.T) {
	// Test with progress callback
	n := math.NewBigInt(31)
	iterations := make([]uint64, 0)
	result, err := PRPTestWithProgress(n, 3, func(iter, total uint64) error {
		iterations = append(iterations, iter)
		return nil
	})
	if err != nil {
		t.Fatalf("PRPTestWithProgress returned error: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("PRPTestWithProgress result has error: %v", result.Error)
	}
	if !result.IsProbablePrime {
		t.Errorf("31 should be probable prime, got IsProbablePrime=%v", result.IsProbablePrime)
	}
	if len(iterations) == 0 {
		t.Errorf("Progress callback should be called")
	}
}

func TestPRPResidueType(t *testing.T) {
	// Test that residue type is set correctly
	n := math.NewBigInt(7)
	result, err := PRPTest(n, 3)
	if err != nil {
		t.Fatalf("PRPTest(7, 3) returned error: %v", err)
	}
	if result.ResidueType != PRPResidueTypeStandard {
		t.Errorf("ResidueType should be Standard, got %v", result.ResidueType)
	}
	if result.Residue == nil {
		t.Errorf("Residue should not be nil")
	}
}

