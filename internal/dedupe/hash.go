package dedupe

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

const fastHashChunkSize = 4096 // 4KB

// computeFastHash reads only the first 4KB of a file. This is extremely fast
// and allows us to filter out files that share a byte size but have completely
// different headers/content, saving us from reading gigabytes of data.
func computeFastHash(filePath string) (string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return "", err
	}

	defer file.Close()

	hash := sha256.New()

	// CopyN ensures we only read the first chunk
	_, err = io.CopyN(hash, file, fastHashChunkSize)

	if err != nil && err != io.EOF {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// computeFullHash streams the entire file from disk through the SHA-256 algorithm.
// It uses io.Copy, which prevents loading the whole file into RAM, ensuring
// the memory footprint stays tiny even when hashing 50GB ISO files.
func computeFullHash(filePath string) (string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return "", err
	}

	defer file.Close()

	hash := sha256.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
