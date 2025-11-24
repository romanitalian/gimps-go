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
	
	if *versionFlag {
		fmt.Printf("gimps version %s\n", version)
		return
	}
	
	if *helpFlag {
		printUsage()
		return
	}
	
	// Load configuration
	cfg := config.NewConfig("prime.ini")
	if err := cfg.Load(); err != nil {
		fmt.Printf("Warning: failed to load config: %v\n", err)
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
		showMenu(cfg)
		return
	}
	
	if *statusFlag {
		showStatus(workToDo)
		return
	}
	
	if *tortureFlag {
		runTortureTest()
		return
	}
	
	if *contactFlag {
		if assignMgr != nil {
			if err := assignMgr.ProcessSpool(); err != nil {
				fmt.Printf("Error processing spool: %v\n", err)
			}
		}
		return
	}
	
	// Default: start workers
	numWorkers := cfg.NumWorkers
	if numWorkers < 1 {
		numWorkers = 1
	}
	
	manager := worker.NewManager(numWorkers, workToDo, assignMgr)
	
	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	// Start workers
	if err := manager.Start(); err != nil {
		fmt.Printf("Failed to start workers: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("GIMPS started with %d worker(s). Press Ctrl+C to stop.\n", numWorkers)
	
	// Wait for signal
	<-sigChan
	fmt.Println("\nShutting down...")
	manager.Stop()
	fmt.Println("Shutdown complete.")
}

func printUsage() {
	fmt.Println("Usage: gimps [options]")
	fmt.Println("Options:")
	fmt.Println("  -m    Menu to configure gimps")
	fmt.Println("  -s    Display status")
	fmt.Println("  -t    Run torture test")
	fmt.Println("  -c    Contact PrimeNet server")
	fmt.Println("  -v    Print version number")
	fmt.Println("  -h    Print this help")
}

func showMenu(cfg *config.Config) {
	fmt.Println("GIMPS Configuration Menu")
	fmt.Println("1. PrimeNet Settings")
	fmt.Println("2. Worker Settings")
	fmt.Println("3. Memory Settings")
	fmt.Println("4. Exit")
	
	var choice int
	fmt.Print("Your choice: ")
	fmt.Scanf("%d", &choice)
	
	switch choice {
	case 1:
		fmt.Print("Use PrimeNet? (0=No, 1=Yes): ")
		var use int
		fmt.Scanf("%d", &use)
		cfg.UsePrimenet = use > 0
		
		if cfg.UsePrimenet {
			fmt.Print("User ID: ")
			fmt.Scanf("%s", &cfg.UserID)
			fmt.Print("Computer Name: ")
			fmt.Scanf("%s", &cfg.ComputerName)
		}
	case 2:
		fmt.Print("Number of workers: ")
		fmt.Scanf("%d", &cfg.NumWorkers)
	case 3:
		fmt.Print("Day Memory (MB): ")
		fmt.Scanf("%d", &cfg.DayMemory)
		fmt.Print("Night Memory (MB): ")
		fmt.Scanf("%d", &cfg.NightMemory)
	}
	
	if err := cfg.Save(); err != nil {
		fmt.Printf("Failed to save config: %v\n", err)
	} else {
		fmt.Println("Configuration saved.")
	}
}

func showStatus(workToDo *storage.WorkToDo) {
	if err := workToDo.Load(); err != nil {
		fmt.Printf("Failed to load worktodo: %v\n", err)
		return
	}
	
	units := workToDo.GetUnits()
	fmt.Printf("Work units in queue: %d\n", len(units))
	
	for i, unit := range units {
		fmt.Printf("%d. %s: n=%d, stage=%s, progress=%.1f%%\n",
			i+1, unit.Stage, unit.N, unit.Stage, unit.PctComplete*100)
	}
}

func runTortureTest() {
	fmt.Println("Torture test not yet implemented")
}

