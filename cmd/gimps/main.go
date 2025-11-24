package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/romanitalian/gimps-go/internal/primenet"
	"github.com/romanitalian/gimps-go/internal/storage"
	"github.com/romanitalian/gimps-go/internal/worker"
	"github.com/romanitalian/gimps-go/pkg/config"
	"github.com/romanitalian/gimps-go/pkg/logger"
)

const (
	version = "1.0.0"
)

func main() {
	var (
		menuFlag      = flag.Bool("m", false, "Menu to configure gimps")
		statusFlag    = flag.Bool("s", false, "Display status")
		tortureFlag   = flag.Bool("t", false, "Run torture test")
		contactFlag   = flag.Bool("c", false, "Contact PrimeNet server")
		versionFlag   = flag.Bool("v", false, "Print version number")
		helpFlag      = flag.Bool("h", false, "Print help")
	)
	flag.Parse()
	
	// Load configuration first (needed for logger)
	cfg := config.NewConfig("prime.ini")
	if err := cfg.Load(); err != nil {
		// Use fmt for error before logger is initialized
		fmt.Printf("Warning: failed to load config: %v\n", err)
	}
	
	// Initialize logger
	appLogger, err := initLogger(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer appLogger.Close()
	
	if *versionFlag {
		appLogger.Info("gimps version " + version)
		return
	}
	
	if *helpFlag {
		printUsage(appLogger)
		return
	}
	
	// Log config warning if needed
	if err := cfg.Load(); err != nil {
		appLogger.Warning("Failed to load config: " + err.Error())
	}
	
	// Initialize components
	workToDo := storage.NewWorkToDo("worktodo.txt")
	
	var assignMgr *primenet.AssignmentManager
	if cfg.UsePrimenet {
		client := primenet.NewClient(cfg.ComputerGUID)
		spool := primenet.NewSpool("prime.spl")
		assignMgr = primenet.NewAssignmentManager(client, spool, workToDo, cfg.ComputerGUID)
	}
	
	// Handle different modes
	if *menuFlag {
		showMenu(cfg, appLogger)
		return
	}
	
	if *statusFlag {
		showStatus(workToDo, appLogger)
		return
	}
	
	if *tortureFlag {
		runTortureTest(appLogger)
		return
	}
	
	if *contactFlag {
		if assignMgr != nil {
			if err := assignMgr.ProcessSpool(); err != nil {
				appLogger.Error("Error processing spool: " + err.Error())
			}
		}
		return
	}
	
	// Default: start workers
	numWorkers := cfg.NumWorkers
	if numWorkers < 1 {
		numWorkers = 1
	}
	
	manager := worker.NewManager(numWorkers, workToDo, assignMgr, appLogger)
	
	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	// Start workers
	if err := manager.Start(); err != nil {
		appLogger.Error("Failed to start workers: " + err.Error())
		os.Exit(1)
	}
	
	appLogger.Info("GIMPS started with " + fmt.Sprintf("%d", numWorkers) + " worker(s). Press Ctrl+C to stop.")
	
	// Wait for signal
	<-sigChan
	appLogger.Info("Shutting down...")
	manager.Stop()
	appLogger.Info("Shutdown complete.")
}

// initLogger initializes logger based on config
func initLogger(cfg *config.Config) (*logger.Logger, error) {
	if cfg.LogFile != "" {
		return logger.NewFileLogger(cfg.LogFile)
	}
	return logger.NewLogger(), nil
}

func printUsage(l *logger.Logger) {
	l.Info("Usage: gimps [options]")
	l.Info("Options:")
	l.Info("  -m    Menu to configure gimps")
	l.Info("  -s    Display status")
	l.Info("  -t    Run torture test")
	l.Info("  -c    Contact PrimeNet server")
	l.Info("  -v    Print version number")
	l.Info("  -h    Print this help")
}

func showMenu(cfg *config.Config, l *logger.Logger) {
	l.Info("GIMPS Configuration Menu")
	l.Info("1. PrimeNet Settings")
	l.Info("2. Worker Settings")
	l.Info("3. Memory Settings")
	l.Info("4. Exit")
	
	var choice int
	fmt.Print("Your choice: ")
	fmt.Scanf("%d", &choice)
	l.Info("Menu choice selected: " + fmt.Sprintf("%d", choice))
	
	switch choice {
	case 1:
		fmt.Print("Use PrimeNet? (0=No, 1=Yes): ")
		var use int
		fmt.Scanf("%d", &use)
		cfg.UsePrimenet = use > 0
		l.Info("PrimeNet setting changed: " + fmt.Sprintf("%v", cfg.UsePrimenet))
		
		if cfg.UsePrimenet {
			fmt.Print("User ID: ")
			fmt.Scanf("%s", &cfg.UserID)
			fmt.Print("Computer Name: ")
			fmt.Scanf("%s", &cfg.ComputerName)
			l.Info("PrimeNet configured: UserID=" + cfg.UserID + ", ComputerName=" + cfg.ComputerName)
		}
	case 2:
		fmt.Print("Number of workers: ")
		fmt.Scanf("%d", &cfg.NumWorkers)
		l.Info("Workers setting changed: " + fmt.Sprintf("%d", cfg.NumWorkers))
	case 3:
		fmt.Print("Day Memory (MB): ")
		fmt.Scanf("%d", &cfg.DayMemory)
		fmt.Print("Night Memory (MB): ")
		fmt.Scanf("%d", &cfg.NightMemory)
		l.Info("Memory settings changed: DayMemory=" + fmt.Sprintf("%d", cfg.DayMemory) + ", NightMemory=" + fmt.Sprintf("%d", cfg.NightMemory))
	}
	
	if err := cfg.Save(); err != nil {
		l.Error("Failed to save config: " + err.Error())
	} else {
		l.Info("Configuration saved.")
	}
}

func showStatus(workToDo *storage.WorkToDo, l *logger.Logger) {
	if err := workToDo.Load(); err != nil {
		l.Error("Failed to load worktodo: " + err.Error())
		return
	}
	
	units := workToDo.GetUnits()
	l.Info("Work units in queue: " + fmt.Sprintf("%d", len(units)))
	
	for i, unit := range units {
		stage := logger.Stage(unit.Stage)
		exponent := unit.N
		progress := unit.PctComplete
		message := fmt.Sprintf("Work unit %d: stage=%s, n=%d, progress=%.1f%%", 
			i+1, unit.Stage, unit.N, progress*100)
		l.WorkerProgress(i+1, stage, exponent, uint64(progress*100), 100, message)
	}
}

func runTortureTest(l *logger.Logger) {
	l.Info("Torture test not yet implemented")
}

