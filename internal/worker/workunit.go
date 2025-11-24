package worker

// WorkType represents the type of work to be performed
type WorkType int32

const (
	// WorkTypeFactor - Trial factoring
	WorkTypeFactor WorkType = 2
	// WorkTypePMinus1 - P-1 factoring
	WorkTypePMinus1 WorkType = 3
	// WorkTypePFactor - P-1 factoring before LL test
	WorkTypePFactor WorkType = 4
	// WorkTypeECM - ECM factoring
	WorkTypeECM WorkType = 5
	// WorkTypePPlus1 - P+1 factoring
	WorkTypePPlus1 WorkType = 6
	// WorkTypeFirstLL - First time Lucas-Lehmer test
	WorkTypeFirstLL WorkType = 100
	// WorkTypeDblChk - Double check Lucas-Lehmer test
	WorkTypeDblChk WorkType = 101
	// WorkTypePRP - PRP test
	WorkTypePRP WorkType = 150
	// WorkTypeCert - Certification
	WorkTypeCert WorkType = 200
)

// WorkUnit represents a single work assignment
// Based on struct work_unit from commonc.h:415
type WorkUnit struct {
	// Work assignment fields
	WorkType        WorkType
	AssignmentUID   string // Primenet assignment ID (max 33 chars)
	Extension       string // Optional save file extension (max 9 chars)
	K               float64
	B               uint64
	N               uint64
	C               int64
	MinimumFFTLen   uint64
	SieveDepth      float64 // How far it has been trial factored
	FactorTo        float64 // How far we should trial factor to
	PMinus1ed       bool    // TRUE if has been P-1 factored
	B1              uint64  // ECM, P-1, P+1 - Stage 1 bound
	B2              uint64  // ECM, P-1, P+1 - Stage 2 bound
	B2Start         uint64  // P-1 - Stage 2 start
	NthRun          int32   // P+1 - 1 for start 2/7, 2 for start 6/5, 3+ for random start
	SkipCurves      uint32  // ECM - number of curves from gmp_ecm_file to skip over
	CurvesToDo      uint32  // ECM - curves to try
	Curve           uint64  // ECM - Specific curve to test (debug tool)
	TestsSaved      float64 // Pfactor - primality tests saved if a factor is found
	PRPBase         uint32  // PRP base to use
	PRPResidueType  int32   // PRP residue to output
	PRPDblChk       bool    // True if this is a doublecheck of a previous PRP
	CertSquarings   int32   // Number of squarings required for PRP proof certification
	GMPECMFile      string  // Save file from GMP-ECM to run stage 2 on
	KnownFactors    string  // ECM, P-1, P+1, PRP - list of known factors
	Comment         string  // Comment line in worktodo.txt

	// Runtime variables
	InUseCount      int32   // Count of threads accessing this work unit
	HighMemoryUsage bool    // Set if we are using a lot of memory
	Stage           string  // Test stage (e.g. TF,P-1,LL) (max 10 chars)
	PctComplete     float64 // Percent complete (0.0 to 1.0)
	FFTLen          uint64  // FFT length in use
	RAFailed        bool    // Set when register assignment fails
}

// IsMersenne returns true if this work unit represents a Mersenne number (2^n - 1)
func (w *WorkUnit) IsMersenne() bool {
	return w.K == 1.0 && w.B == 2 && w.C == -1
}

// GetMersenneExponent returns the exponent for Mersenne number, or 0 if not a Mersenne
func (w *WorkUnit) GetMersenneExponent() uint64 {
	if w.IsMersenne() {
		return w.N
	}
	return 0
}

