package config

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
)

// AppConfig holds the runtime configuration for the scanning engine.
type AppConfig struct {
	TargetDir string
	Workers   int
	Mode      string
}

// Load parses command line flags and returns a populated AppConfig.
func Load() *AppConfig {
	var (
		dirFlag     string
		workersFlag int
		modeFlag    string
	)

	// Default to scanning the current working directory to prevent
	// accidental full-drive scans without explicit user intent.
	defaultDir, err := os.Getwd()

	if err != nil {
		defaultDir = "."
	}

	flag.StringVar(&dirFlag, "dir", defaultDir, "Target directory or drive to scan (e.g., C:\\ or ./)")
	flag.IntVar(&workersFlag, "workers", runtime.NumCPU(), "Number of concurrent hashing workers")
	flag.StringVar(&modeFlag, "mode", "all", "Execution mode: 'all', 'sweep' (caches only), or 'dedupe' (duplicates only)")

	// Parse the flags if they haven't been parsed yet
	if !flag.Parsed() {
		flag.Parse()
	}

	// Clean the input directory path to ensure consistent separator usage
	cleanDir := filepath.Clean(dirFlag)

	return &AppConfig{
		TargetDir: cleanDir,
		Workers:   workersFlag,
		Mode:      modeFlag,
	}
}
