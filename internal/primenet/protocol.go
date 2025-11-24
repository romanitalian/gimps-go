package primenet

// Operation types for PrimeNet protocol
const (
	OpUpdateComputerInfo  = 100
	OpProgramOptions      = 101
	OpGetAssignment       = 102
	OpRegisterAssignment  = 103
	OpAssignmentProgress  = 104
	OpAssignmentResult    = 105
	OpAssignmentUnreserve = 106
	OpBenchmarkData       = 107
	OpPingServer          = 108
)

// Work preference values
const (
	WPWhatever        = 0
	WPFactorLMH       = 1
	WPFactor          = 2
	WPPMinus1         = 3
	WPPFactor         = 4
	WPECMSmall        = 5
	WPECMFermat       = 6
	WPECMCunningham   = 7
	WPECMCofactor     = 8
	WPLLFirst         = 100
	WPLLDblChk        = 101
	WPLLWorldRecord   = 102
	WPLL100M          = 104
	WPPRPFirst        = 150
	WPPRPDblChk       = 151
	WPPRPWorldRecord  = 152
	WPPRP100M         = 153
	WPPRPNoPMinus1    = 154
	WPPRPDCProof      = 155
	WPPRPCofactor     = 160
	WPPRPCofactorDblChk = 161
)

// Work type values
const (
	WorkTypeFactor  = 2
	WorkTypePMinus1 = 3
	WorkTypePFactor = 4
	WorkTypeECM     = 5
	WorkTypePPlus1  = 6
	WorkTypeFirstLL = 100
	WorkTypeDblChk  = 101
	WorkTypePRP     = 150
	WorkTypeCert    = 200
)

// UpdateComputerInfo structure for PrimeNet protocol
// Based on struct primenetUpdateComputerInfo from primenet.h:50
type UpdateComputerInfo struct {
	VersionNumber   int32
	ComputerGUID    string // max 33 chars
	HardwareGUID     string // max 33 chars
	WindowsGUID      string // max 33 chars
	Application      string // max 65 chars
	CPUModel         string // max 65 chars
	CPUFeatures      string // max 65 chars
	L1CacheSize      int32  // Cache size in KB
	L2CacheSize      int32  // Cache size in KB
	L3CacheSize      int32  // Cache size in KB
	NumCPUs          uint32
	NumHyperthread   uint32
	MemInstalled     uint32 // Physical memory in MB
	CPUSpeed         uint32 // CPU speed in MHz
	HoursPerDay      uint32
	RollingAverage   uint32
	UserID           string // max 21 chars
	ComputerName     string // max 21 chars

	// Returned by server
	UserName         string // max 33 chars
	OptionsCounter   uint32
}

// ProgramOptions structure for PrimeNet protocol
// Based on struct primenetProgramOptions from primenet.h:113
type ProgramOptions struct {
	VersionNumber   int32
	ComputerGUID    string // max 33 chars
	CPUNum          int32  // CPU number (-1 = all)

	// Read/write parameters. Minus 1 is used to specify parameter not passed in
	NumWorkers      int32
	WorkPreference  int32
	Priority        int32
	DaysOfWork      int32
	DayMemory       int32
	NightMemory     int32
	DayStartTime    int32
	NightStartTime  int32
	RunOnBattery    int32

	// Returned by server
	OptionsCounter  uint32
}

// GetAssignment structure for PrimeNet protocol
// Based on struct primenetGetAssignment from primenet.h:154
type GetAssignment struct {
	VersionNumber   int32
	ComputerGUID    string // max 33 chars
	CPUNum          uint32
	GetCertWork     int32  // If trying to get cert work, set to CertDailyCPULimit
	TempDiskSpace   float32
	MinExp          uint32
	MaxExp          uint32

	// Returned by server
	AssignmentUID   string // max 33 chars
	WorkType        uint32
	K               float64
	B               uint32
	N               uint32
	C               int32
	HasBeenPMinus1ed uint32
	HowFarFactored  float64
	FactorTo        float64
	B1              uint64
	B2              uint64
	TestsSaved      float64
	Curves          uint32
	PRPBase         uint32
	PRPResidueType  uint32
	PRPDblChk       uint32
	NumSquarings    uint32
	NthRun          uint32
	KnownFactors    string // max 2000 chars
}

// RegisterAssignment structure for PrimeNet protocol
// Based on struct primenetRegisterAssignment from primenet.h:193
type RegisterAssignment struct {
	VersionNumber   int32
	ComputerGUID    string // max 33 chars
	CPUNum          uint32
	WorkType        uint32
	K               float64
	B               uint32
	N               uint32
	C               int32
	HowFarFactored  float64
	FactorTo        float64
	HasBeenPMinus1ed uint32
	B1              uint64
	B2              uint64
	TestsSaved      float64
	Curves          uint32
	PRPBase         uint32
	PRPResidueType  uint32
	PRPDblChk       uint32
	NumSquarings    uint32
	NthRun          uint32
}

// AssignmentProgress structure for PrimeNet protocol
type AssignmentProgress struct {
	VersionNumber   int32
	ComputerGUID    string // max 33 chars
	CPUNum          uint32
	AssignmentUID   string // max 33 chars
	Stage           string // max 10 chars
	PctComplete     float64
	CurrentIteration uint64
	Residue         string // 64-bit residue in hex
	FFTLen          uint32
}

// AssignmentResult structure for PrimeNet protocol
type AssignmentResult struct {
	VersionNumber   int32
	ComputerGUID    string // max 33 chars
	CPUNum          uint32
	AssignmentUID   string // max 33 chars
	WorkType        uint32
	K               float64
	B               uint32
	N               uint32
	C               int32
	Residue         string // 64-bit residue in hex
	ResidueType     uint32
	ErrorCode       uint32
	ErrorMessage    string
	ProofFiles      string
}

// AssignmentUnreserve structure for PrimeNet protocol
type AssignmentUnreserve struct {
	VersionNumber   int32
	ComputerGUID    string // max 33 chars
	CPUNum          uint32
	AssignmentUID   string // max 33 chars
}

// BenchmarkData structure for PrimeNet protocol
type BenchmarkData struct {
	VersionNumber   int32
	ComputerGUID    string // max 33 chars
	CPUNum          uint32
	FFTSize         uint32
	IterationTime   float64
}

