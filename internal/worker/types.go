package worker

// WorkUnitGetter is an interface for getting work units
// This breaks circular dependency between worker and storage
type WorkUnitGetter interface {
	GetUnits() []*WorkUnit
	AddUnit(unit *WorkUnit)
	RemoveUnit(index int) error
	Save() error
}

