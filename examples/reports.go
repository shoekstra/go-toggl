//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	toggl "github.com/shoekstra/go-toggl"
)

func main() {
	client, err := toggl.NewClient(os.Getenv("TOGGL_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	wsID, err := strconv.Atoi(os.Getenv("TOGGL_WORKSPACE_ID"))
	if err != nil {
		log.Fatal("TOGGL_WORKSPACE_ID must be set to a valid workspace ID (run examples/workspaces.go to find yours)")
	}

	ctx := context.Background()

	filters := toggl.ReportFilters{
		StartDate: toggl.String("2026-01-01"),
		EndDate:   toggl.String("2026-03-01"),
	}

	// Summary report — grouped by project.
	summary, _, err := client.Reports.SummaryReport(ctx, wsID, &toggl.SummaryReportOptions{
		ReportFilters: filters,
		Grouping:      toggl.String("projects"),
		SubGrouping:   toggl.String("users"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Summary groups: %d\n", len(summary.Groups))
	for _, g := range summary.Groups {
		fmt.Printf("  group %v: %ds\n", g.ID, g.Seconds)
	}

	// Detailed report — individual time entries with cursor-based pagination.
	fmt.Println("\nDetailed entries:")
	var firstID, firstRowNumber *int
	for {
		entries, resp, err := client.Reports.DetailedReport(ctx, wsID, &toggl.DetailedReportOptions{
			ReportFilters:  filters,
			FirstID:        firstID,
			FirstRowNumber: firstRowNumber,
		})
		if err != nil {
			log.Fatal(err)
		}

		for _, e := range entries {
			fmt.Printf("  row %d: project %v\n", e.RowNumber, e.ProjectID)
		}

		if resp.Pagination.NextID == 0 {
			break
		}
		firstID = toggl.Int(resp.Pagination.NextID)
		firstRowNumber = toggl.Int(resp.Pagination.NextRowNumber)
	}

	// Totals report.
	totals, _, err := client.Reports.DetailedReportTotals(ctx, wsID, &toggl.DetailedReportOptions{
		ReportFilters: filters,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nTotal: %ds across %d days\n", totals.Seconds, totals.TrackedDays)

	// Weekly report.
	weekly, _, err := client.Reports.WeeklyReport(ctx, wsID, &toggl.WeeklyReportOptions{
		ReportFilters: filters,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Weekly entries: %d\n", len(weekly))

	// Export summary as PDF (returns an error if the workspace has no tracked time).
	pdf, _, err := client.Reports.ExportSummaryPDF(ctx, wsID, &toggl.SummaryExportOptions{
		SummaryReportOptions: toggl.SummaryReportOptions{ReportFilters: filters},
		DateFormat:           toggl.String("YYYY-MM-DD"),
	})
	if err != nil {
		fmt.Printf("PDF export skipped: %v\n", err)
	} else {
		fmt.Printf("PDF size: %d bytes\n", len(pdf))
	}

	// Export detailed report as CSV (returns an error if the workspace has no tracked time).
	csv, _, err := client.Reports.ExportDetailedCSV(ctx, wsID, &toggl.DetailedExportOptions{
		DetailedReportOptions: toggl.DetailedReportOptions{ReportFilters: filters},
	})
	if err != nil {
		fmt.Printf("CSV export skipped: %v\n", err)
	} else {
		fmt.Printf("CSV size: %d bytes\n", len(csv))
	}
}
