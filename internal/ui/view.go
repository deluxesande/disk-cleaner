package ui

import (
	"fmt"
	"strings"

	"github.com/deluxesande/disk-cleaner/internal/models"
)

// RenderSweepReport formats the results of a quick sweep using Lip Gloss colors.
func RenderSweepReport(report models.DiskReport) string {
	var b strings.Builder

	// Main Header
	b.WriteString("\n")
	b.WriteString(HeaderStyle.Render(" SWEEP RESULTS "))
	b.WriteString("\n\n")

	// Helper function to render a category
	renderCategory := func(title string, items []models.SpaceWaster) {
		if len(items) == 0 {
			return
		}

		totalMB := float64(calculateTotalSavings(items)) / (1024 * 1024)
		titleText := fmt.Sprintf("%s (Total: %.2f MB)", title, totalMB)

		b.WriteString(CategoryTitleStyle.Render(titleText))
		b.WriteString("\n")

		for _, item := range items {
			mb := float64(item.Size) / (1024 * 1024)

			sizeStr := FormatSize(mb)
			pathStr := PathStyle.Render(item.Path)

			b.WriteString(fmt.Sprintf("  %s  %s\n", sizeStr, pathStr))
		}
		b.WriteString("\n")
	}

	renderCategory("Development Artifacts", report.DevArtifacts)
	renderCategory("Application Caches", report.AppCaches)
	renderCategory("Temporary Files", report.TempFiles)

	b.WriteString(SubtleStyle.Render(strings.Repeat("-", 60)))
	b.WriteString("\n")

	totalSavingsMB := float64(report.TotalSavings) / (1024 * 1024)
	b.WriteString(fmt.Sprintf("SPACE RECOVERED: %s\n", FormatSize(totalSavingsMB)))

	return b.String()
}

// RenderDedupeReport formats the verified duplicate groups using Lip Gloss colors.
func RenderDedupeReport(duplicates []models.DuplicateGroup) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(HeaderStyle.Render(" DEDUPE RESULTS "))
	b.WriteString("\n\n")

	if len(duplicates) == 0 {
		b.WriteString(SubtleStyle.Render("No duplicates found."))
		b.WriteString("\n")
		return b.String()
	}

	var totalWasted int64
	for i, group := range duplicates {
		wasted := int64(len(group.Instances)-1) * group.FileSize
		totalWasted += wasted

		titleText := fmt.Sprintf("Group %d: %d identical files (Wasted: %.2f MB)", i+1, len(group.Instances), float64(wasted)/(1024*1024))

		b.WriteString(CategoryTitleStyle.Render(titleText))
		b.WriteString("\n")

		for _, path := range group.Instances {
			b.WriteString(fmt.Sprintf("  - %s\n", PathStyle.Render(path)))
		}
		b.WriteString("\n")
	}

	b.WriteString(SubtleStyle.Render(strings.Repeat("-", 60)))
	b.WriteString("\n")

	totalWastedMB := float64(totalWasted) / (1024 * 1024)
	b.WriteString(fmt.Sprintf("TOTAL REDUNDANT SPACE: %s\n", FormatSize(totalWastedMB)))

	return b.String()
}

func calculateTotalSavings(items []models.SpaceWaster) int64 {
	var total int64
	for _, item := range items {
		total += item.Size
	}
	return total
}
