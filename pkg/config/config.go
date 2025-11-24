package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config represents application configuration
type Config struct {
	// PrimeNet settings
	UsePrimenet      bool
	UserID           string
	ComputerName     string
	ComputerGUID     string
	HardwareGUID     string

	// Worker settings
	NumWorkers       int
	WorkPreference   int
	Priority         int

	// Memory settings
	DayMemory        int // MB
	NightMemory      int // MB
	DayStartTime     int // Hour (0-23)
	NightStartTime   int // Hour (0-23)

	// Power settings
	RunOnBattery     bool

	// Algorithm settings
	SkipTrialFactoring bool
	JacobiErrorCheck   bool

	// Other settings
	StressTester    bool
	Verbose         bool
	NoGUI           bool

	// Logging settings
	LogFile         string // Path to log file (empty = stdout)
	LogLevel        string // Log level: debug, info, warning, error

	// Internal
	filePath        string
	values          map[string]string
}

// NewConfig creates a new config instance
func NewConfig(filePath string) *Config {
	return &Config{
		filePath: filePath,
		values:   make(map[string]string),
		// Defaults
		NumWorkers:     1,
		WorkPreference: 0, // Whatever makes most sense
		Priority:       1,
		DayMemory:      8,
		NightMemory:     8,
		DayStartTime:    0,
		NightStartTime:  0,
		RunOnBattery:    false,
		SkipTrialFactoring: false,
		JacobiErrorCheck:  true,
		StressTester:    false,
		Verbose:         false,
		NoGUI:           false,
		LogFile:         "", // Default: stdout
		LogLevel:        "info", // Default: info
	}
}

// Load loads configuration from INI file
func (c *Config) Load() error {
	data, err := os.ReadFile(c.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, use defaults
			return nil
		}
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		c.values[key] = value

		// Parse known keys
		c.parseKey(key, value)
	}

	return nil
}

// parseKey parses a specific configuration key
func (c *Config) parseKey(key, value string) {
	switch strings.ToLower(key) {
	case "useprimenet":
		c.UsePrimenet = c.getBool(value, false)
	case "userid":
		c.UserID = value
	case "computername":
		c.ComputerName = value
	case "computerguid":
		c.ComputerGUID = value
	case "hardwareguid":
		c.HardwareGUID = value
	case "workers":
		c.NumWorkers = c.getInt(value, 1)
	case "workpreference":
		c.WorkPreference = c.getInt(value, 0)
	case "priority":
		c.Priority = c.getInt(value, 1)
	case "daymemory":
		c.DayMemory = c.getInt(value, 8)
	case "nightmemory":
		c.NightMemory = c.getInt(value, 8)
	case "daystarttime":
		c.DayStartTime = c.getInt(value, 0)
	case "nightstarttime":
		c.NightStartTime = c.getInt(value, 0)
	case "runonbattery":
		c.RunOnBattery = c.getBool(value, false)
	case "skiptrialfactoring":
		c.SkipTrialFactoring = c.getBool(value, false)
	case "jacobierrorcheck":
		c.JacobiErrorCheck = c.getBool(value, true)
	case "stresstester":
		c.StressTester = c.getBool(value, false)
	case "verbose":
		c.Verbose = c.getBool(value, false)
	case "nogui":
		c.NoGUI = c.getBool(value, false)
	case "logfile":
		c.LogFile = value
	case "loglevel":
		c.LogLevel = value
	}
}

// Save saves configuration to INI file
func (c *Config) Save() error {
	var sb strings.Builder

	// Write known settings
	c.writeKey(&sb, "UsePrimenet", c.getBoolStr(c.UsePrimenet))
	c.writeKey(&sb, "UserID", c.UserID)
	c.writeKey(&sb, "ComputerName", c.ComputerName)
	c.writeKey(&sb, "ComputerGUID", c.ComputerGUID)
	c.writeKey(&sb, "HardwareGUID", c.HardwareGUID)
	c.writeKey(&sb, "Workers", fmt.Sprintf("%d", c.NumWorkers))
	c.writeKey(&sb, "WorkPreference", fmt.Sprintf("%d", c.WorkPreference))
	c.writeKey(&sb, "Priority", fmt.Sprintf("%d", c.Priority))
	c.writeKey(&sb, "DayMemory", fmt.Sprintf("%d", c.DayMemory))
	c.writeKey(&sb, "NightMemory", fmt.Sprintf("%d", c.NightMemory))
	c.writeKey(&sb, "DayStartTime", fmt.Sprintf("%d", c.DayStartTime))
	c.writeKey(&sb, "NightStartTime", fmt.Sprintf("%d", c.NightStartTime))
	c.writeKey(&sb, "RunOnBattery", c.getBoolStr(c.RunOnBattery))
	c.writeKey(&sb, "SkipTrialFactoring", c.getBoolStr(c.SkipTrialFactoring))
	c.writeKey(&sb, "JacobiErrorCheck", c.getBoolStr(c.JacobiErrorCheck))
	c.writeKey(&sb, "StressTester", c.getBoolStr(c.StressTester))
	c.writeKey(&sb, "Verbose", c.getBoolStr(c.Verbose))
	c.writeKey(&sb, "NoGUI", c.getBoolStr(c.NoGUI))
	c.writeKey(&sb, "LogFile", c.LogFile)
	c.writeKey(&sb, "LogLevel", c.LogLevel)

	// Write other values
	for k, v := range c.values {
		if !c.isKnownKey(k) {
			c.writeKey(&sb, k, v)
		}
	}

	return os.WriteFile(c.filePath, []byte(sb.String()), 0644)
}

// GetInt returns integer value for key, or default if not found
func (c *Config) GetInt(key string, defaultValue int) int {
	if val, ok := c.values[key]; ok {
		return c.getInt(val, defaultValue)
	}
	return defaultValue
}

// GetBool returns boolean value for key, or default if not found
func (c *Config) GetBool(key string, defaultValue bool) bool {
	if val, ok := c.values[key]; ok {
		return c.getBool(val, defaultValue)
	}
	return defaultValue
}

// GetString returns string value for key, or default if not found
func (c *Config) GetString(key string, defaultValue string) string {
	if val, ok := c.values[key]; ok {
		return val
	}
	return defaultValue
}

// Helper methods
func (c *Config) getInt(s string, defaultValue int) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return val
}

func (c *Config) getBool(s string, defaultValue bool) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "1" || s == "true" || s == "yes"
}

func (c *Config) getBoolStr(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func (c *Config) writeKey(sb *strings.Builder, key, value string) {
	sb.WriteString(fmt.Sprintf("%s=%s\n", key, value))
}

func (c *Config) isKnownKey(key string) bool {
	knownKeys := []string{
		"useprimenet", "userid", "computername", "computerguid", "hardwareguid",
		"workers", "workpreference", "priority",
		"daymemory", "nightmemory", "daystarttime", "nightstarttime",
		"runonbattery", "skiptrialfactoring", "jacobierrorcheck",
		"stresstester", "verbose", "nogui",
		"logfile", "loglevel",
	}
	keyLower := strings.ToLower(key)
	for _, k := range knownKeys {
		if keyLower == k {
			return true
		}
	}
	return false
}

