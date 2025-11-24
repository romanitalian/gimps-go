package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/romanitalian/gimps-go/internal/algorithms"
	"github.com/romanitalian/gimps-go/internal/math"
	"github.com/romanitalian/gimps-go/pkg/logger"
)

// AssignmentManager is an interface for assignment management
// This breaks circular dependency
type AssignmentManager interface {
	GetAssignment(cpuNum uint32) (*WorkUnit, error)
	RegisterAssignment(unit *WorkUnit, cpuNum uint32) error
	SendProgress(unit *WorkUnit, cpuNum uint32, iteration uint64) error
	SendResult(unit *WorkUnit, cpuNum uint32, residue string, isPrime bool, errorCode uint32, errorMsg string) error
}

// Manager manages worker threads
type Manager struct {
	workers      []*Worker
	workToDo     WorkUnitGetter
	assignMgr    AssignmentManager
	numWorkers   int
	logger       *logger.Logger
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	mu           sync.Mutex
}

// NewManager creates a new worker manager
func NewManager(numWorkers int, workToDo WorkUnitGetter, assignMgr AssignmentManager, l *logger.Logger) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		workers:    make([]*Worker, 0, numWorkers),
		workToDo:   workToDo,
		assignMgr:  assignMgr,
		numWorkers: numWorkers,
		logger:     l,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start starts all workers
func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Start workers
	for i := 0; i < m.numWorkers; i++ {
		worker := NewWorker(i, m.workToDo, m.assignMgr, m.ctx, m.logger)
		m.workers = append(m.workers, worker)
		
		m.wg.Add(1)
		go func(w *Worker) {
			defer m.wg.Done()
			w.Run()
		}(worker)
	}
	
	return nil
}

// Stop stops all workers
func (m *Manager) Stop() {
	m.cancel()
	m.wg.Wait()
}

// GetStatus returns status of all workers
func (m *Manager) GetStatus() []WorkerStatus {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	status := make([]WorkerStatus, len(m.workers))
	for i, worker := range m.workers {
		status[i] = worker.GetStatus()
	}
	return status
}

// WorkerStatus represents the status of a worker
type WorkerStatus struct {
	WorkerNum    int
	IsRunning    bool
	CurrentWork  *WorkUnit
	Progress     float64
	Stage        string
}

// Worker represents a single worker thread
type Worker struct {
	num          int
	workToDo     WorkUnitGetter
	assignMgr    AssignmentManager
	ctx          context.Context
	logger       *logger.Logger
	currentWork  *WorkUnit
	mu           sync.Mutex
}

// NewWorker creates a new worker
func NewWorker(num int, workToDo WorkUnitGetter, assignMgr AssignmentManager, ctx context.Context, l *logger.Logger) *Worker {
	return &Worker{
		num:       num,
		workToDo:  workToDo,
		assignMgr: assignMgr,
		ctx:       ctx,
		logger:    l,
	}
}

// Run runs the worker main loop
func (w *Worker) Run() {
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
			// Get next work unit
			unit := w.getNextWork()
			if unit == nil {
				// Try to get assignment from server
				if w.assignMgr != nil {
					newUnit, err := w.assignMgr.GetAssignment(uint32(w.num))
					if err != nil {
						time.Sleep(30 * time.Second)
						continue
					}
					unit = newUnit
				} else {
					time.Sleep(5 * time.Second)
					continue
				}
			}
			
			// Process work unit
			w.processWork(unit)
		}
	}
}

// getNextWork gets the next work unit from worktodo
func (w *Worker) getNextWork() *WorkUnit {
	units := w.workToDo.GetUnits()
	for _, unit := range units {
		if unit.InUseCount == 0 {
			unit.InUseCount = 1
			return unit
		}
	}
	return nil
}

// processWork processes a work unit
func (w *Worker) processWork(unit *WorkUnit) {
	w.mu.Lock()
	w.currentWork = unit
	w.mu.Unlock()
	
	defer func() {
		w.mu.Lock()
		w.currentWork = nil
		w.mu.Unlock()
		unit.InUseCount = 0
	}()
	
	// Register assignment if needed
	if unit.AssignmentUID == "" && w.assignMgr != nil {
		w.assignMgr.RegisterAssignment(unit, uint32(w.num))
	}
	
	// Process based on work type
	var err error
	switch unit.WorkType {
	case WorkTypeFirstLL, WorkTypeDblChk:
		err = w.processLL(unit.K, unit.B, unit.N, unit.C)
	case WorkTypePRP:
		err = w.processPRP(unit.K, unit.B, unit.N, unit.C, unit.PRPBase)
	case WorkTypeFactor:
		err = w.processTF(unit.K, unit.B, unit.N, unit.C, unit.FactorTo)
	case WorkTypePMinus1:
		err = w.processPM1(unit.K, unit.B, unit.N, unit.C, unit.B1, unit.B2)
	case WorkTypePPlus1:
		err = w.processPP1(unit.K, unit.B, unit.N, unit.C, unit.B1, unit.B2, unit.NthRun)
	case WorkTypeECM:
		err = w.processECM(unit.K, unit.B, unit.N, unit.C, unit.B1, unit.B2, unit.CurvesToDo)
	}
	
	if err != nil {
		var stage logger.Stage = logger.StageIdle
		var exponent uint64
		if w.currentWork != nil {
			stage = logger.Stage(w.currentWork.Stage)
			exponent = w.currentWork.N
		}
		w.logger.WorkerError(w.num, stage, exponent, "Error processing work: "+err.Error())
	}
}

// processLL processes Lucas-Lehmer test
func (w *Worker) processLL(k float64, b, n uint64, c int64) error {
	w.currentWork.Stage = "LL"
	stage := logger.StageLL
	
	// Use progress callback for logging
	result, err := algorithms.LucasLehmerTestWithProgress(n, func(iteration, total uint64, residue *math.BigInt) error {
		// Log progress periodically (every 1% or every 1000 iterations, whichever is more frequent)
		if total > 0 && (iteration%1000 == 0 || iteration*100/total != (iteration-1)*100/total) {
			w.logger.WorkerProgress(w.num, stage, n, iteration, total, 
				fmt.Sprintf("Lucas-Lehmer iteration %d/%d", iteration, total))
		}
		return nil
	})
	
	if err != nil {
		return err
	}
	
	// Log completion
	w.logger.WorkerProgress(w.num, stage, n, n-2, n-2, 
		fmt.Sprintf("Lucas-Lehmer test completed: isPrime=%v", result.IsPrime))
	
	// Send result
	if w.assignMgr != nil {
		residue := "0000000000000000"
		if result.Residue != nil {
			residue = result.Residue.Text(16)
		}
		w.assignMgr.SendResult(w.currentWork, uint32(w.num), residue, result.IsPrime, 0, "")
	}
	
	return nil
}

// processPRP processes PRP test
func (w *Worker) processPRP(k float64, b, n uint64, c int64, prpBase uint32) error {
	w.currentWork.Stage = "PRP"
	stage := logger.StagePRP
	
	result, err := algorithms.PRPTestFromWorkUnit(k, b, n, c, prpBase)
	if err != nil {
		return err
	}
	
	// Log completion
	w.logger.WorkerProgress(w.num, stage, n, 100, 100, fmt.Sprintf("PRP test completed for exponent %d: isProbablePrime=%v", 
		n, result.IsProbablePrime))
	
	// Send result
	if w.assignMgr != nil {
		residue := "0000000000000000"
		if result.Residue != nil {
			residue = result.Residue.Text(16)
		}
		w.assignMgr.SendResult(w.currentWork, uint32(w.num), residue, result.IsProbablePrime, 0, "")
	}
	
	return nil
}

// processTF processes Trial Factoring
func (w *Worker) processTF(k float64, b, n uint64, c int64, factorTo float64) error {
	w.currentWork.Stage = "TF"
	stage := logger.StageTF
	
	w.logger.WorkerProgress(w.num, stage, n, 0, 100, fmt.Sprintf("Starting trial factoring for exponent %d", n))
	
	result, err := algorithms.TrialFactorFromWorkUnit(k, b, n, c, factorTo)
	if err != nil {
		return err
	}
	
	if len(result.Factors) > 0 {
		w.logger.WorkerProgress(w.num, stage, n, 100, 100, fmt.Sprintf("Factor found for exponent %d: %s", 
			n, result.Factors[0].Text(10)))
		// Factor found
		if w.assignMgr != nil {
			factorStr := result.Factors[0].Text(10)
			w.assignMgr.SendResult(w.currentWork, uint32(w.num), factorStr, false, 0, "")
		}
	} else {
		w.logger.WorkerProgress(w.num, stage, n, 100, 100, fmt.Sprintf("No factors found for exponent %d", n))
	}
	
	return nil
}

// processPM1 processes P-1 factoring
func (w *Worker) processPM1(k float64, b, n uint64, c int64, B1, B2 uint64) error {
	w.currentWork.Stage = "P-1"
	stage := logger.StagePM1
	
	w.logger.WorkerProgress(w.num, stage, n, 0, 100, fmt.Sprintf("Starting P-1 factoring for exponent %d (B1=%d, B2=%d)", 
		n, B1, B2))
	
	result, err := algorithms.PMinus1FactorFromWorkUnit(k, b, n, c, B1, B2)
	if err != nil {
		return err
	}
	
	if result.Factor != nil {
		w.logger.WorkerProgress(w.num, stage, n, 100, 100, fmt.Sprintf("P-1 factor found for exponent %d (stage %d): %s", 
			n, result.Stage, result.Factor.Text(10)))
		// Factor found
		if w.assignMgr != nil {
			factorStr := result.Factor.Text(10)
			w.assignMgr.SendResult(w.currentWork, uint32(w.num), factorStr, false, 0, "")
		}
	} else {
		w.logger.WorkerProgress(w.num, stage, n, 100, 100, fmt.Sprintf("No P-1 factor found for exponent %d", n))
	}
	
	return nil
}

// processPP1 processes P+1 factoring
func (w *Worker) processPP1(k float64, b, n uint64, c int64, B1, B2 uint64, nthRun int32) error {
	w.currentWork.Stage = "P+1"
	stage := logger.StagePP1
	
	w.logger.WorkerProgress(w.num, stage, n, 0, 100, fmt.Sprintf("Starting P+1 factoring for exponent %d (B1=%d, B2=%d, run=%d)", 
		n, B1, B2, nthRun))
	
	result, err := algorithms.PPlus1FactorFromWorkUnit(k, b, n, c, B1, B2, nthRun)
	if err != nil {
		return err
	}
	
	if result.Factor != nil {
		w.logger.WorkerProgress(w.num, stage, n, 100, 100, fmt.Sprintf("P+1 factor found for exponent %d (stage %d): %s", 
			n, result.Stage, result.Factor.Text(10)))
		// Factor found
		if w.assignMgr != nil {
			factorStr := result.Factor.Text(10)
			w.assignMgr.SendResult(w.currentWork, uint32(w.num), factorStr, false, 0, "")
		}
	} else {
		w.logger.WorkerProgress(w.num, stage, n, 100, 100, fmt.Sprintf("No P+1 factor found for exponent %d", n))
	}
	
	return nil
}

// processECM processes ECM factoring
func (w *Worker) processECM(k float64, b, n uint64, c int64, B1, B2 uint64, numCurves uint32) error {
	w.currentWork.Stage = "ECM"
	stage := logger.StageECM
	
	w.logger.WorkerProgress(w.num, stage, n, 0, 100, fmt.Sprintf("Starting ECM factoring for exponent %d (B1=%d, B2=%d, curves=%d)", 
		n, B1, B2, numCurves))
	
	result, err := algorithms.ECMFactorFromParams(k, b, n, c, B1, B2, numCurves)
	if err != nil {
		return err
	}
	
	if result.Factor != nil {
		w.logger.WorkerProgress(w.num, stage, n, 100, 100, fmt.Sprintf("ECM factor found for exponent %d (curve %d, stage %d): %s", 
			n, result.Curve, result.Stage, result.Factor.Text(10)))
		// Factor found
		if w.assignMgr != nil {
			factorStr := result.Factor.Text(10)
			w.assignMgr.SendResult(w.currentWork, uint32(w.num), factorStr, false, 0, "")
		}
	} else {
		w.logger.WorkerProgress(w.num, stage, n, 100, 100, fmt.Sprintf("No ECM factor found for exponent %d after %d curves", 
			n, numCurves))
	}
	
	return nil
}

// GetStatus returns the current status of the worker
func (w *Worker) GetStatus() WorkerStatus {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	status := WorkerStatus{
		WorkerNum: w.num,
		IsRunning: w.currentWork != nil,
		CurrentWork: w.currentWork,
	}
	
	if w.currentWork != nil {
		status.Progress = w.currentWork.PctComplete
		status.Stage = w.currentWork.Stage
	}
	
	return status
}

