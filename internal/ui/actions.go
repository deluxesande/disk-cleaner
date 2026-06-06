package ui

import (
	"log"
	"os"
)

// DeleteItem attempts to remove a file or directory at the given path.
// It uses RemoveAll for directories to recursively clear contents.
func DeleteItem(path string) error {
	info, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if info.IsDir() {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}

	if err != nil {
		log.Printf("Failed to delete %s: %v\n", path, err)
		return err
	}

	return nil
}

// DeleteDuplicates safely removes all redundant instances in a duplicate group.
// It explicitly preserves the first file in the slice and deletes the rest.
func DeleteDuplicates(group []string) []error {
	var errors []error

	if len(group) <= 1 {
		return nil
	}

	// Start iterating at index 1, explicitly skipping the first file
	for i := 1; i < len(group); i++ {
		err := DeleteItem(group[i])
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
