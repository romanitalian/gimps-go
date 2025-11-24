package storage

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/romanitalian/gimps-go/internal/worker"
)

// WorkToDo manages the worktodo.txt file
type WorkToDo struct {
	filePath string
	units   []*worker.WorkUnit
}

// NewWorkToDo creates a new WorkToDo manager
func NewWorkToDo(filePath string) *WorkToDo {
	return &WorkToDo{
		filePath: filePath,
		units:   []*worker.WorkUnit{},
	}
}

// Load loads work units from worktodo.txt file
func (wtd *WorkToDo) Load() error {
	file, err := os.Open(wtd.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist, that's OK
		}
		return fmt.Errorf("failed to open worktodo file: %w", err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	wtd.units = []*worker.WorkUnit{}
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		unit, err := parseWorkToDoLine(line)
		if err != nil {
			continue // Skip invalid lines
		}
		
		wtd.units = append(wtd.units, unit)
	}
	
	return scanner.Err()
}

// Save saves work units to worktodo.txt file
func (wtd *WorkToDo) Save() error {
	file, err := os.Create(wtd.filePath)
	if err != nil {
		return fmt.Errorf("failed to create worktodo file: %w", err)
	}
	defer file.Close()
	
	for _, unit := range wtd.units {
		line := formatWorkToDoLine(unit)
		if _, err := file.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	
	return nil
}

// AddUnit adds a work unit
func (wtd *WorkToDo) AddUnit(unit *worker.WorkUnit) {
	wtd.units = append(wtd.units, unit)
}

// GetUnits returns all work units
func (wtd *WorkToDo) GetUnits() []*worker.WorkUnit {
	return wtd.units
}

// RemoveUnit removes a work unit by index
func (wtd *WorkToDo) RemoveUnit(index int) error {
	if index < 0 || index >= len(wtd.units) {
		return fmt.Errorf("invalid index")
	}
	wtd.units = append(wtd.units[:index], wtd.units[index+1:]...)
	return nil
}

// parseWorkToDoLine parses a line from worktodo.txt
// Format examples:
// Test=12345678
// DoubleCheck=12345678
// Factor=12345678,70
// Pminus1=1,2,12345678,-1,50000,5000000
func parseWorkToDoLine(line string) (*worker.WorkUnit, error) {
	unit := &worker.WorkUnit{
		K: 1.0,
		B: 2,
		C: -1,
	}
	
	parts := strings.Split(line, "=")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid line format")
	}
	
	workTypeStr := strings.TrimSpace(parts[0])
	paramsStr := strings.TrimSpace(parts[1])
	
	// Parse work type
	switch workTypeStr {
	case "Test":
		unit.WorkType = worker.WorkTypeFirstLL
	case "DoubleCheck":
		unit.WorkType = worker.WorkTypeDblChk
	case "Factor":
		unit.WorkType = worker.WorkTypeFactor
	case "Pminus1":
		unit.WorkType = worker.WorkTypePMinus1
	case "Pplus1":
		unit.WorkType = worker.WorkTypePPlus1
	case "ECM":
		unit.WorkType = worker.WorkTypeECM
	case "PRP":
		unit.WorkType = worker.WorkTypePRP
	default:
		return nil, fmt.Errorf("unknown work type: %s", workTypeStr)
	}
	
	// Parse parameters
	params := strings.Split(paramsStr, ",")
	
	if unit.WorkType == worker.WorkTypeFirstLL || unit.WorkType == worker.WorkTypeDblChk {
		// Test=12345678 or DoubleCheck=12345678
		if len(params) >= 1 {
			if n, err := strconv.ParseUint(params[0], 10, 64); err == nil {
				unit.N = n
			}
		}
	} else if unit.WorkType == worker.WorkTypeFactor {
		// Factor=12345678,70
		if len(params) >= 1 {
			if n, err := strconv.ParseUint(params[0], 10, 64); err == nil {
				unit.N = n
			}
		}
		if len(params) >= 2 {
			if factorTo, err := strconv.ParseFloat(params[1], 64); err == nil {
				unit.FactorTo = factorTo
			}
		}
	} else if unit.WorkType == worker.WorkTypePMinus1 || unit.WorkType == worker.WorkTypePPlus1 {
		// Pminus1=1,2,12345678,-1,50000,5000000
		if len(params) >= 1 {
			if k, err := strconv.ParseFloat(params[0], 64); err == nil {
				unit.K = k
			}
		}
		if len(params) >= 2 {
			if b, err := strconv.ParseUint(params[1], 10, 64); err == nil {
				unit.B = b
			}
		}
		if len(params) >= 3 {
			if n, err := strconv.ParseUint(params[2], 10, 64); err == nil {
				unit.N = n
			}
		}
		if len(params) >= 4 {
			if c, err := strconv.ParseInt(params[3], 10, 64); err == nil {
				unit.C = c
			}
		}
		if len(params) >= 5 {
			if b1, err := strconv.ParseUint(params[4], 10, 64); err == nil {
				unit.B1 = b1
			}
		}
		if len(params) >= 6 {
			if b2, err := strconv.ParseUint(params[5], 10, 64); err == nil {
				unit.B2 = b2
			}
		}
	}
	
	return unit, nil
}

// formatWorkToDoLine formats a work unit as a worktodo.txt line
func formatWorkToDoLine(unit *worker.WorkUnit) string {
	var sb strings.Builder
	
	switch unit.WorkType {
	case worker.WorkTypeFirstLL:
		sb.WriteString(fmt.Sprintf("Test=%d", unit.N))
	case worker.WorkTypeDblChk:
		sb.WriteString(fmt.Sprintf("DoubleCheck=%d", unit.N))
	case worker.WorkTypeFactor:
		sb.WriteString(fmt.Sprintf("Factor=%d,%.0f", unit.N, unit.FactorTo))
	case worker.WorkTypePMinus1:
		sb.WriteString(fmt.Sprintf("Pminus1=%.0f,%d,%d,%d,%d,%d",
			unit.K, unit.B, unit.N, unit.C, unit.B1, unit.B2))
	case worker.WorkTypePPlus1:
		sb.WriteString(fmt.Sprintf("Pplus1=%.0f,%d,%d,%d,%d,%d",
			unit.K, unit.B, unit.N, unit.C, unit.B1, unit.B2))
	case worker.WorkTypeECM:
		sb.WriteString(fmt.Sprintf("ECM=%.0f,%d,%d,%d,%d,%d,%d",
			unit.K, unit.B, unit.N, unit.C, unit.B1, unit.B2, unit.CurvesToDo))
	case worker.WorkTypePRP:
		sb.WriteString(fmt.Sprintf("PRP=%.0f,%d,%d,%d",
			unit.K, unit.B, unit.N, unit.C))
	}
	
	if unit.Comment != "" {
		sb.WriteString(" ; " + unit.Comment)
	}
	
	return sb.String()
}

