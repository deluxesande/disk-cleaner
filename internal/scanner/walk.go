package scanner

import (
	"io/fs"
	"log"
	"path/filepath"
	"time"

	"github.com/deluxesande/disk-cleaner/internal/config"
	"github.com/deluxesande/disk-cleaner/internal/models"
)

// RunSweep executes a fast traversal of the target directory to locate,
// measure, and group known space-wasting directories.
func RunSweep(cfg *config.AppConfig) models.DiskReport {
	report := models.DiskReport{}

	err := filepath.WalkDir(cfg.TargetDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() && config.IsExcluded(path) {
			return filepath.SkipDir
		}

		category := IdentifyCategory(path, d.IsDir())

		if category != "" {
			size, lastAccess := calculateDirMetrics(path)

			waster := models.SpaceWaster{
				Path:       path,
				Size:       size,
				LastAccess: lastAccess,
				Category:   category,
			}

			switch category {
			case "DevArtifacts":
				report.DevArtifacts = append(report.DevArtifacts, waster)
			case "AppCaches":
				report.AppCaches = append(report.AppCaches, waster)
			case "TempFiles":
				report.TempFiles = append(report.TempFiles, waster)
			}

			report.TotalSavings += size

			// Skip descending into the junk folder to save I/O cycles
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		log.Printf("Scanner encountered a fatal error: %v\n", err)
	}

	return report
}

// calculateDirMetrics recursively calculates the total byte size of a directory
// and finds the most recent modification time among its contents.
func calculateDirMetrics(dirPath string) (int64, time.Time) {
	var totalSize int64
	var lastAccess time.Time

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		info, err := d.Info()
		if err == nil && !d.IsDir() {
			totalSize += info.Size()
			if info.ModTime().After(lastAccess) {
				lastAccess = info.ModTime()
			}
		}
		return nil
	})

	if err != nil {
		// Log silently or ignore, returning whatever size was calculated before the permission error
		_ = err
	}

	return totalSize, lastAccess
}
