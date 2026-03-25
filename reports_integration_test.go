//go:build integration

package toggl_test

import (
	"testing"
	"time"

	toggl "github.com/shoekstra/go-toggl"
)

func TestIntegration_Reports_SummaryReport(t *testing.T) {
	client := integrationClient(t)
	wsID := integrationWorkspaceID(t)
	ctx := integrationCtx(t)

	filters := recentFilters()

	data, _, err := client.Reports.SummaryReport(ctx, wsID, &toggl.SummaryReportOptions{
		ReportFilters: filters,
		Grouping:      toggl.String("projects"),
	})
	if err != nil {
		t.Fatalf("SummaryReport: %v", err)
	}
	if data == nil {
		t.Fatal("SummaryReport returned nil data")
	}
	// An empty workspace has no groups; that's valid.
	t.Logf("summary groups: %d", len(data.Groups))
}

func TestIntegration_Reports_DetailedReport(t *testing.T) {
	client := integrationClient(t)
	wsID := integrationWorkspaceID(t)
	ctx := integrationCtx(t)

	filters := recentFilters()

	entries, resp, err := client.Reports.DetailedReport(ctx, wsID, &toggl.DetailedReportOptions{
		ReportFilters: filters,
	})
	if err != nil {
		t.Fatalf("DetailedReport: %v", err)
	}
	t.Logf("detailed entries: %d, next_id: %d", len(entries), resp.Pagination.NextID)
}

func TestIntegration_Reports_DetailedReportTotals(t *testing.T) {
	client := integrationClient(t)
	wsID := integrationWorkspaceID(t)
	ctx := integrationCtx(t)

	totals, _, err := client.Reports.DetailedReportTotals(ctx, wsID, &toggl.DetailedReportOptions{
		ReportFilters: recentFilters(),
	})
	if err != nil {
		t.Fatalf("DetailedReportTotals: %v", err)
	}
	if totals == nil {
		t.Fatal("DetailedReportTotals returned nil data")
	}
	t.Logf("total seconds: %d, tracked days: %d", totals.Seconds, totals.TrackedDays)
}

func TestIntegration_Reports_WeeklyReport(t *testing.T) {
	client := integrationClient(t)
	wsID := integrationWorkspaceID(t)
	ctx := integrationCtx(t)

	entries, _, err := client.Reports.WeeklyReport(ctx, wsID, &toggl.WeeklyReportOptions{
		ReportFilters: recentFilters(),
	})
	if err != nil {
		t.Fatalf("WeeklyReport: %v", err)
	}
	t.Logf("weekly entries: %d", len(entries))
}

func TestIntegration_Reports_ExportDetailedCSV(t *testing.T) {
	client := integrationClient(t)
	wsID := integrationWorkspaceID(t)
	ctx := integrationCtx(t)

	data, _, err := client.Reports.ExportDetailedCSV(ctx, wsID, &toggl.DetailedExportOptions{
		DetailedReportOptions: toggl.DetailedReportOptions{ReportFilters: recentFilters()},
	})
	if err != nil {
		t.Fatalf("ExportDetailedCSV: %v", err)
	}
	if len(data) == 0 {
		t.Error("ExportDetailedCSV returned empty bytes")
	}
	t.Logf("CSV size: %d bytes", len(data))
}

// recentFilters returns a ReportFilters covering the last 30 days.
// Toggl rejects queries older than ~3 months.
func recentFilters() toggl.ReportFilters {
	now := time.Now().UTC()
	return toggl.ReportFilters{
		StartDate: toggl.String(now.AddDate(0, -1, 0).Format("2006-01-02")),
		EndDate:   toggl.String(now.Format("2006-01-02")),
	}
}
