package primenet

import (
	"fmt"

	"github.com/romanitalian/gimps-go/internal/worker"
)

// AssignmentManager manages PrimeNet assignments
type AssignmentManager struct {
	client    *Client
	spool     *Spool
	workToDo  worker.WorkUnitGetter
	computerGUID string
}

// NewAssignmentManager creates a new assignment manager
func NewAssignmentManager(client *Client, spool *Spool, workToDo worker.WorkUnitGetter, computerGUID string) *AssignmentManager {
	return &AssignmentManager{
		client:      client,
		spool:       spool,
		workToDo:    workToDo,
		computerGUID: computerGUID,
	}
}

// GetAssignment requests a new assignment from PrimeNet server
func (am *AssignmentManager) GetAssignment(cpuNum uint32) (*worker.WorkUnit, error) {
	req := &GetAssignment{
		VersionNumber: 5,
		ComputerGUID:  am.computerGUID,
		CPUNum:        cpuNum,
	}
	
	resp, err := am.client.GetAssignment(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignment: %w", err)
	}
	
	// Convert response to WorkUnit
	unit := am.convertAssignmentToWorkUnit(resp)
	
	// Add to worktodo
	am.workToDo.AddUnit(unit)
	if err := am.workToDo.Save(); err != nil {
		return nil, fmt.Errorf("failed to save worktodo: %w", err)
	}
	
	return unit, nil
}

// RegisterAssignment registers an assignment with the server
func (am *AssignmentManager) RegisterAssignment(unit *worker.WorkUnit, cpuNum uint32) error {
	req := &RegisterAssignment{
		VersionNumber: 5,
		ComputerGUID:  am.computerGUID,
		CPUNum:        cpuNum,
		WorkType:      uint32(unit.WorkType),
		K:             unit.K,
		B:             uint32(unit.B),
		N:             uint32(unit.N),
		C:             int32(unit.C),
		HowFarFactored: unit.SieveDepth,
		FactorTo:      unit.FactorTo,
		HasBeenPMinus1ed: boolToUint32(unit.PMinus1ed),
		B1:            unit.B1,
		B2:            unit.B2,
		TestsSaved:    unit.TestsSaved,
		Curves:        unit.CurvesToDo,
		PRPBase:       unit.PRPBase,
		PRPResidueType: uint32(unit.PRPResidueType),
		PRPDblChk:     boolToUint32(unit.PRPDblChk),
	}
	
	if err := am.client.RegisterAssignment(req); err != nil {
		// Add to spool for retry
		data := serializeRegisterAssignment(req)
		am.spool.AddMessage(OpRegisterAssignment, data)
		return err
	}
	
	// AssignmentUID is set when we get assignment from server
	return nil
}

// SendProgress sends progress update for an assignment
func (am *AssignmentManager) SendProgress(unit *worker.WorkUnit, cpuNum uint32, iteration uint64) error {
	req := &AssignmentProgress{
		VersionNumber:   5,
		ComputerGUID:    am.computerGUID,
		CPUNum:          cpuNum,
		AssignmentUID:   unit.AssignmentUID,
		Stage:           unit.Stage,
		PctComplete:     unit.PctComplete,
		CurrentIteration: iteration,
		FFTLen:          uint32(unit.FFTLen),
	}
	
	if err := am.client.AssignmentProgress(req); err != nil {
		// Add to spool for retry
		data := serializeAssignmentProgress(req)
		am.spool.AddMessage(OpAssignmentProgress, data)
		return err
	}
	
	return nil
}

// SendResult sends result of an assignment
func (am *AssignmentManager) SendResult(unit *worker.WorkUnit, cpuNum uint32, residue string, isPrime bool, errorCode uint32, errorMsg string) error {
	req := &AssignmentResult{
		VersionNumber: 5,
		ComputerGUID:  am.computerGUID,
		CPUNum:        cpuNum,
		AssignmentUID: unit.AssignmentUID,
		WorkType:      uint32(unit.WorkType),
		K:             unit.K,
		B:             uint32(unit.B),
		N:             uint32(unit.N),
		C:             int32(unit.C),
		Residue:       residue,
		ResidueType:   uint32(unit.PRPResidueType),
		ErrorCode:     errorCode,
		ErrorMessage:  errorMsg,
	}
	
	if err := am.client.AssignmentResult(req); err != nil {
		// Add to spool for retry
		data := serializeAssignmentResult(req)
		am.spool.AddMessage(OpAssignmentResult, data)
		return err
	}
	
	return nil
}

// ProcessSpool processes messages in the spool file
func (am *AssignmentManager) ProcessSpool() error {
	messages, err := am.spool.ReadMessages()
	if err != nil {
		return err
	}
	
	for i, msg := range messages {
		var err error
		switch msg.Type {
		case OpRegisterAssignment:
			req := deserializeRegisterAssignment(msg.Data)
			err = am.client.RegisterAssignment(req)
		case OpAssignmentProgress:
			req := deserializeAssignmentProgress(msg.Data)
			err = am.client.AssignmentProgress(req)
		case OpAssignmentResult:
			req := deserializeAssignmentResult(msg.Data)
			err = am.client.AssignmentResult(req)
		}
		
		if err == nil {
			// Successfully sent, remove from spool
			am.spool.RemoveMessage(i)
		}
	}
	
	return nil
}

// convertAssignmentToWorkUnit converts PrimeNet assignment to WorkUnit
func (am *AssignmentManager) convertAssignmentToWorkUnit(assign *GetAssignment) *worker.WorkUnit {
	unit := &worker.WorkUnit{
		WorkType:        worker.WorkType(assign.WorkType),
		AssignmentUID:   assign.AssignmentUID,
		K:               assign.K,
		B:               uint64(assign.B),
		N:               uint64(assign.N),
		C:               int64(assign.C),
		SieveDepth:      assign.HowFarFactored,
		FactorTo:        assign.FactorTo,
		PMinus1ed:       assign.HasBeenPMinus1ed > 0,
		B1:              assign.B1,
		B2:              assign.B2,
		TestsSaved:      assign.TestsSaved,
		CurvesToDo:      assign.Curves,
		PRPBase:         assign.PRPBase,
		PRPResidueType:  int32(assign.PRPResidueType),
		PRPDblChk:       assign.PRPDblChk > 0,
		KnownFactors:    assign.KnownFactors,
	}
	
	return unit
}

// Helper functions for serialization (simplified - would need proper encoding)
func serializeRegisterAssignment(req *RegisterAssignment) []byte {
	// Simplified - would use proper encoding like gob or json
	return []byte{}
}

func serializeAssignmentProgress(req *AssignmentProgress) []byte {
	return []byte{}
}

func serializeAssignmentResult(req *AssignmentResult) []byte {
	return []byte{}
}

func deserializeRegisterAssignment(data []byte) *RegisterAssignment {
	return &RegisterAssignment{}
}

func deserializeAssignmentProgress(data []byte) *AssignmentProgress {
	return &AssignmentProgress{}
}

func deserializeAssignmentResult(data []byte) *AssignmentResult {
	return &AssignmentResult{}
}

func boolToUint32(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

