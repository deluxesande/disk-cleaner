package models

import "time"

type SpaceWaster struct {
	Path       string
	Size       int64
	LastAccess time.Time
	Category   string
}

type DuplicateGroup struct {
	Checksum  string
	FileSize  int64
	Instances []string
}

type DiskReport struct {
	DevArtifacts []SpaceWaster
	AppCaches    []SpaceWaster
	Duplicates   []DuplicateGroup
	TempFiles    []SpaceWaster
	TotalSavings int64
}
