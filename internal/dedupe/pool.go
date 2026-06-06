package dedupe

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sync"

	"github.com/deluxesande/disk-cleaner/internal/config"
	"github.com/deluxesande/disk-cleaner/internal/models"
)

type fileJob struct {
	path string
	size int64
}

type hashResult struct {
	job  fileJob
	hash string
}

// Run executes the complete multi-pass deduplication pipeline.
func Run(cfg *config.AppConfig) []models.DuplicateGroup {
	// PASS 1: The Size Sieve (I/O Bound)
	sizeGroups := make(map[int64][]string)

	err := filepath.WalkDir(cfg.TargetDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || config.IsExcluded(path) {
			if d != nil && d.IsDir() && config.IsExcluded(path) {
				return filepath.SkipDir
			}
			return nil
		}

		info, err := d.Info()
		if err == nil && info.Mode().IsRegular() {
			sizeGroups[info.Size()] = append(sizeGroups[info.Size()], path)
		}
		return nil
	})

	if err != nil {
		log.Printf("Error during dedupe size traversal: %v\n", err)
	}

	var candidateJobs []fileJob
	for size, paths := range sizeGroups {
		if len(paths) > 1 {
			for _, p := range paths {
				candidateJobs = append(candidateJobs, fileJob{path: p, size: size})
			}
		}
	}

	if len(candidateJobs) == 0 {
		return nil
	}

	// PASS 2: Fast-Hash Filter (First 4KB)
	fastHashGroups := processHashPool(candidateJobs, cfg.Workers, computeFastHash)

	var fullHashJobs []fileJob
	for _, jobs := range fastHashGroups {
		if len(jobs) > 1 {
			fullHashJobs = append(fullHashJobs, jobs...)
		}
	}

	if len(fullHashJobs) == 0 {
		return nil
	}

	// PASS 3: Full Cryptographic Hash (Entire File)
	fullHashGroups := processHashPool(fullHashJobs, cfg.Workers, computeFullHash)

	// PASS 4: Final Purge and Formatting
	var finalDuplicates []models.DuplicateGroup
	for hashKey, jobs := range fullHashGroups {
		if len(jobs) > 1 {
			group := models.DuplicateGroup{
				Checksum: hashKey,
				FileSize: jobs[0].size,
			}
			for _, j := range jobs {
				group.Instances = append(group.Instances, j.path)
			}
			finalDuplicates = append(finalDuplicates, group)
		}
	}

	return finalDuplicates
}

// processHashPool is a generic worker pool that accepts either computeFastHash or computeFullHash.
func processHashPool(jobs []fileJob, workerCount int, hashFunc func(string) (string, error)) map[string][]fileJob {
	jobsChan := make(chan fileJob, len(jobs))
	resultsChan := make(chan hashResult, len(jobs))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobsChan {
				hash, err := hashFunc(job.path)
				if err == nil {
					// Combine size and hash to guarantee uniqueness across different sizes
					// that might coincidentally share a 4KB prefix.
					combinedKey := fmt.Sprintf("%d-%s", job.size, hash)
					resultsChan <- hashResult{job: job, hash: combinedKey}
				}
			}
		}()
	}

	// Feed jobs into the channel
	for _, job := range jobs {
		jobsChan <- job
	}
	close(jobsChan)

	// Wait for workers in the background and close results
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Aggregate final results
	grouped := make(map[string][]fileJob)
	for res := range resultsChan {
		grouped[res.hash] = append(grouped[res.hash], res.job)
	}

	return grouped
}
