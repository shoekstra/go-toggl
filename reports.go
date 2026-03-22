package toggl

import (
	"context"
	"fmt"
)

// ReportsService handles operations related to reports.
// It uses the Reports API v3 base path (/reports/api/v3/) rather than the
// standard API v9 path used by other services.
type ReportsService struct {
	client *Client
}

const reportsBasePath = "/reports/api/v3"

// ReportFilters holds the common filtering fields shared across all report
// types. Embed this struct in a report-specific options type.
type ReportFilters struct {
	// StartDate is the report start date (YYYY-MM-DD).
	StartDate *string `json:"start_date,omitempty"`
	// EndDate is the report end date (YYYY-MM-DD). Must be after StartDate.
	EndDate *string `json:"end_date,omitempty"`
	// ClientIDs filters by client. Pass a slice containing 0 to include
	// entries with no client.
	ClientIDs []int `json:"client_ids,omitempty"`
	// ProjectIDs filters by project. Pass a slice containing 0 to include
	// entries with no project.
	ProjectIDs []int `json:"project_ids,omitempty"`
	// UserIDs filters by user.
	UserIDs []int `json:"user_ids,omitempty"`
	// TagIDs filters by tag. Pass a slice containing 0 to include entries
	// with no tags.
	TagIDs []int `json:"tag_ids,omitempty"`
	// GroupIDs filters by user group.
	GroupIDs []int `json:"group_ids,omitempty"`
	// TaskIDs filters by task.
	TaskIDs []int `json:"task_ids,omitempty"`
	// TimeEntryIDs filters to specific time entry IDs.
	TimeEntryIDs []int `json:"time_entry_ids,omitempty"`
	// Description filters by description text.
	Description *string `json:"description,omitempty"`
	// Billable filters by billable status (premium feature).
	Billable *bool `json:"billable,omitempty"`
	// Rounding controls time rounding. Defaults to user preferences.
	Rounding *int `json:"rounding,omitempty"`
	// RoundingMinutes sets the rounding interval in minutes.
	// Valid values: 0, 1, 5, 6, 10, 12, 15, 30, 60, 240.
	RoundingMinutes *int `json:"rounding_minutes,omitempty"`
	// MinDurationSeconds filters out entries shorter than this value (Time Audit).
	MinDurationSeconds *int `json:"min_duration_seconds,omitempty"`
	// MaxDurationSeconds filters out entries longer than this value (Time Audit).
	MaxDurationSeconds *int `json:"max_duration_seconds,omitempty"`
}

// SummaryAudit configures the Time Audit view for a summary report.
type SummaryAudit struct {
	// GroupFilter restricts which groups appear in the audit view.
	GroupFilter *SummaryAuditFilter `json:"group_filter,omitempty"`
	// ShowEmptyGroups includes groups with no tracked time.
	ShowEmptyGroups *bool `json:"show_empty_groups,omitempty"`
	// ShowTrackedGroups includes groups with tracked time.
	ShowTrackedGroups *bool `json:"show_tracked_groups,omitempty"`
}

// SummaryAuditFilter constrains entries shown in the Time Audit group.
type SummaryAuditFilter struct {
	// Currency filters by currency code.
	Currency *string `json:"currency,omitempty"`
	// MaxAmountCents is the upper billable amount bound in cents.
	MaxAmountCents *int `json:"max_amount_cents,omitempty"`
	// MaxDurationSeconds is the upper duration bound in seconds.
	MaxDurationSeconds *int `json:"max_duration_seconds,omitempty"`
	// MinAmountCents is the lower billable amount bound in cents.
	MinAmountCents *int `json:"min_amount_cents,omitempty"`
	// MinDurationSeconds is the lower duration bound in seconds.
	MinDurationSeconds *int `json:"min_duration_seconds,omitempty"`
}

// SummaryReportOptions specifies the parameters to ReportsService.SummaryReport.
type SummaryReportOptions struct {
	ReportFilters
	// Grouping sets the top-level grouping (e.g. "projects", "clients", "users").
	Grouping *string `json:"grouping,omitempty"`
	// SubGrouping sets the secondary grouping.
	SubGrouping *string `json:"sub_grouping,omitempty"`
	// DistinguishRates creates sub-groups for each billing rate.
	DistinguishRates *bool `json:"distinguish_rates,omitempty"`
	// IncludeTimeEntryIDs includes time entry IDs in the response.
	IncludeTimeEntryIDs *bool `json:"include_time_entry_ids,omitempty"`
	// Audit configures the Time Audit view.
	Audit *SummaryAudit `json:"audit,omitempty"`
}

// SummaryExportOptions specifies the parameters to ReportsService.ExportSummaryPDF
// and ReportsService.ExportSummaryCSV.
type SummaryExportOptions struct {
	SummaryReportOptions
	// DateFormat controls the date format in the export.
	// Valid values: "MM/DD/YYYY", "DD-MM-YYYY", "MM-DD-YYYY", "YYYY-MM-DD",
	// "DD/MM/YYYY", "DD.MM.YYYY". Defaults to "MM/DD/YYYY".
	DateFormat *string `json:"date_format,omitempty"`
	// DurationFormat controls the duration format.
	// Valid values: "classic", "decimal", "improved". Defaults to "classic".
	DurationFormat *string `json:"duration_format,omitempty"`
	// OrderBy sets the sort column. Valid values: "title", "duration".
	OrderBy *string `json:"order_by,omitempty"`
	// OrderDir sets the sort direction. Valid values: "ASC", "DESC".
	OrderDir *string `json:"order_dir,omitempty"`
	// Collapse collapses other entries. Defaults to false.
	Collapse *bool `json:"collapse,omitempty"`
	// HideRates omits rate information from the export.
	HideRates *bool `json:"hide_rates,omitempty"`
	// HideAmounts omits monetary amounts from the export.
	HideAmounts *bool `json:"hide_amounts,omitempty"`
	// CentsSeparator sets the separator character for cent values.
	CentsSeparator *string `json:"cents_separator,omitempty"`
	// Resolution sets the graph resolution (PDF only).
	Resolution *string `json:"resolution,omitempty"`
}

// DetailedReportOptions specifies the parameters to ReportsService.DetailedReport
// and ReportsService.DetailedReportTotals.
type DetailedReportOptions struct {
	ReportFilters
	// OrderBy sets the sort column.
	// Valid values: "date", "user", "duration", "description", "last_update".
	// Defaults to "date".
	OrderBy *string `json:"order_by,omitempty"`
	// OrderDir sets the sort direction. Valid values: "ASC", "DESC".
	OrderDir *string `json:"order_dir,omitempty"`
	// PageSize sets the number of entries per page. Defaults to 50.
	PageSize *int `json:"page_size,omitempty"`
	// Grouped groups time entries together. Defaults to false.
	Grouped *bool `json:"grouped,omitempty"`
	// EnrichResponse returns the maximum amount of information. Defaults to false.
	EnrichResponse *bool `json:"enrich_response,omitempty"`
	// HideAmounts omits monetary amounts from the response.
	HideAmounts *bool `json:"hide_amounts,omitempty"`
	// FirstID is the cursor ID for pagination (from X-Next-ID response header).
	FirstID *int `json:"first_id,omitempty"`
	// FirstRowNumber is the cursor row number for pagination
	// (from X-Next-Row-Number response header).
	FirstRowNumber *int `json:"first_row_number,omitempty"`
	// FirstTimestamp is the cursor timestamp for pagination.
	FirstTimestamp *int `json:"first_timestamp,omitempty"`
}

// DetailedExportOptions specifies the parameters to ReportsService.ExportDetailedPDF
// and ReportsService.ExportDetailedCSV.
type DetailedExportOptions struct {
	DetailedReportOptions
	// DateFormat controls the date format in the export.
	// Valid values: "MM/DD/YYYY", "DD-MM-YYYY", "MM-DD-YYYY", "YYYY-MM-DD",
	// "DD/MM/YYYY", "DD.MM.YYYY".
	DateFormat *string `json:"date_format,omitempty"`
	// DurationFormat controls the duration format.
	// Valid values: "classic", "decimal", "improved". Defaults to "classic".
	DurationFormat *string `json:"duration_format,omitempty"`
	// DisplayMode controls how date/time is displayed in the PDF.
	// Valid values: "date_only", "time_only", "date_time", "date_and_time".
	DisplayMode *string `json:"display_mode,omitempty"`
	// HourFormat controls the hour format in the PDF.
	HourFormat *string `json:"hour_format,omitempty"`
	// CentsSeparator sets the separator character for cent values.
	CentsSeparator *string `json:"cents_separator,omitempty"`
}

// WeeklyReportOptions specifies the parameters to ReportsService.WeeklyReport.
type WeeklyReportOptions struct {
	ReportFilters
}

// WeeklyExportOptions specifies the parameters to ReportsService.ExportWeeklyPDF
// and ReportsService.ExportWeeklyCSV.
type WeeklyExportOptions struct {
	WeeklyReportOptions
	// Calculate sets the calculation mode. Valid values: "by time", "amounts".
	Calculate *string `json:"calculate,omitempty"`
	// GroupByTask groups entries by planned task. Defaults to false.
	GroupByTask *bool `json:"group_by_task,omitempty"`
	// Grouping sets the grouping option.
	Grouping *string `json:"grouping,omitempty"`
	// DateFormat controls the date format in the PDF export.
	// Valid values: "MM/DD/YYYY", "DD-MM-YYYY", "MM-DD-YYYY", "YYYY-MM-DD",
	// "DD/MM/YYYY", "DD.MM.YYYY". Defaults to "MM/DD/YYYY".
	DateFormat *string `json:"date_format,omitempty"`
	// DurationFormat controls the duration format.
	// Valid values: "classic", "decimal", "improved". Defaults to "classic".
	DurationFormat *string `json:"duration_format,omitempty"`
	// LogoURL sets the URL of the logo to include in the PDF.
	LogoURL *string `json:"logo_url,omitempty"`
	// CentsSeparator sets the separator character for cent values.
	CentsSeparator *string `json:"cents_separator,omitempty"`
}

// SummaryReport returns aggregated time data for the given workspace grouped
// according to the provided options.
//
// API: POST /reports/api/v3/workspace/{workspace_id}/summary/time_entries
//
// See: https://engineering.toggl.com/docs/reports/summary
func (s *ReportsService) SummaryReport(ctx context.Context, workspaceID int, opts *SummaryReportOptions) (*SummaryReportData, *Response, error) {
	if opts == nil {
		opts = &SummaryReportOptions{}
	}

	path := fmt.Sprintf("%s/workspace/%d/summary/time_entries", reportsBasePath, workspaceID)

	data := new(SummaryReportData)
	resp, err := s.client.post(ctx, path, opts, data)
	if err != nil {
		return nil, resp, err
	}

	return data, resp, nil
}

// ExportSummaryPDF downloads a summary report as a PDF file.
//
// API: POST /reports/api/v3/workspace/{workspace_id}/summary/time_entries.pdf
//
// See: https://engineering.toggl.com/docs/reports/summary
func (s *ReportsService) ExportSummaryPDF(ctx context.Context, workspaceID int, opts *SummaryExportOptions) ([]byte, *Response, error) {
	if opts == nil {
		opts = &SummaryExportOptions{}
	}

	path := fmt.Sprintf("%s/workspace/%d/summary/time_entries.pdf", reportsBasePath, workspaceID)

	resp, err := s.client.post(ctx, path, opts, nil)
	if err != nil {
		return nil, resp, err
	}

	return resp.Body, resp, nil
}

// ExportSummaryCSV downloads a summary report as a CSV file.
// To download as XLSX, use the same options but the underlying path ends in
// .xlsx — see the Toggl Reports API documentation.
//
// API: POST /reports/api/v3/workspace/{workspace_id}/summary/time_entries.csv
//
// See: https://engineering.toggl.com/docs/reports/summary
func (s *ReportsService) ExportSummaryCSV(ctx context.Context, workspaceID int, opts *SummaryExportOptions) ([]byte, *Response, error) {
	if opts == nil {
		opts = &SummaryExportOptions{}
	}

	path := fmt.Sprintf("%s/workspace/%d/summary/time_entries.csv", reportsBasePath, workspaceID)

	resp, err := s.client.post(ctx, path, opts, nil)
	if err != nil {
		return nil, resp, err
	}

	return resp.Body, resp, nil
}

// DetailedReport returns individual time entries for the given workspace
// matching the provided filters. Pagination is cursor-based: read the
// X-Next-ID and X-Next-Row-Number response headers and pass them back as
// FirstID and FirstRowNumber in subsequent calls.
//
// API: POST /reports/api/v3/workspace/{workspace_id}/search/time_entries
//
// See: https://engineering.toggl.com/docs/reports/detailed
func (s *ReportsService) DetailedReport(ctx context.Context, workspaceID int, opts *DetailedReportOptions) ([]*DetailedTimeEntry, *Response, error) {
	if opts == nil {
		opts = &DetailedReportOptions{}
	}

	path := fmt.Sprintf("%s/workspace/%d/search/time_entries", reportsBasePath, workspaceID)

	var entries []*DetailedTimeEntry
	resp, err := s.client.post(ctx, path, opts, &entries)
	if err != nil {
		return nil, resp, err
	}

	return entries, resp, nil
}

// DetailedReportTotals returns the aggregate totals for a detailed report,
// including total seconds, billable amounts, and a day-by-day graph.
//
// API: POST /reports/api/v3/workspace/{workspace_id}/search/time_entries/totals
//
// See: https://engineering.toggl.com/docs/reports/detailed
func (s *ReportsService) DetailedReportTotals(ctx context.Context, workspaceID int, opts *DetailedReportOptions) (*TotalsReport, *Response, error) {
	if opts == nil {
		opts = &DetailedReportOptions{}
	}

	path := fmt.Sprintf("%s/workspace/%d/search/time_entries/totals", reportsBasePath, workspaceID)

	totals := new(TotalsReport)
	resp, err := s.client.post(ctx, path, opts, totals)
	if err != nil {
		return nil, resp, err
	}

	return totals, resp, nil
}

// ExportDetailedPDF downloads a detailed report as a PDF file.
//
// API: POST /reports/api/v3/workspace/{workspace_id}/search/time_entries.pdf
//
// See: https://engineering.toggl.com/docs/reports/detailed
func (s *ReportsService) ExportDetailedPDF(ctx context.Context, workspaceID int, opts *DetailedExportOptions) ([]byte, *Response, error) {
	if opts == nil {
		opts = &DetailedExportOptions{}
	}

	path := fmt.Sprintf("%s/workspace/%d/search/time_entries.pdf", reportsBasePath, workspaceID)

	resp, err := s.client.post(ctx, path, opts, nil)
	if err != nil {
		return nil, resp, err
	}

	return resp.Body, resp, nil
}

// ExportDetailedCSV downloads a detailed report as a CSV file.
//
// API: POST /reports/api/v3/workspace/{workspace_id}/search/time_entries.csv
//
// See: https://engineering.toggl.com/docs/reports/detailed
func (s *ReportsService) ExportDetailedCSV(ctx context.Context, workspaceID int, opts *DetailedExportOptions) ([]byte, *Response, error) {
	if opts == nil {
		opts = &DetailedExportOptions{}
	}

	path := fmt.Sprintf("%s/workspace/%d/search/time_entries.csv", reportsBasePath, workspaceID)

	resp, err := s.client.post(ctx, path, opts, nil)
	if err != nil {
		return nil, resp, err
	}

	return resp.Body, resp, nil
}

// WeeklyReport returns a week-by-week breakdown of tracked time for the given
// workspace.
//
// API: POST /reports/api/v3/workspace/{workspace_id}/weekly/time_entries
//
// See: https://engineering.toggl.com/docs/reports/weekly
func (s *ReportsService) WeeklyReport(ctx context.Context, workspaceID int, opts *WeeklyReportOptions) ([]*WeeklyReportEntry, *Response, error) {
	if opts == nil {
		opts = &WeeklyReportOptions{}
	}

	path := fmt.Sprintf("%s/workspace/%d/weekly/time_entries", reportsBasePath, workspaceID)

	var entries []*WeeklyReportEntry
	resp, err := s.client.post(ctx, path, opts, &entries)
	if err != nil {
		return nil, resp, err
	}

	return entries, resp, nil
}

// ExportWeeklyPDF downloads a weekly report as a PDF file.
//
// API: POST /reports/api/v3/workspace/{workspace_id}/weekly/time_entries.pdf
//
// See: https://engineering.toggl.com/docs/reports/weekly
func (s *ReportsService) ExportWeeklyPDF(ctx context.Context, workspaceID int, opts *WeeklyExportOptions) ([]byte, *Response, error) {
	if opts == nil {
		opts = &WeeklyExportOptions{}
	}

	path := fmt.Sprintf("%s/workspace/%d/weekly/time_entries.pdf", reportsBasePath, workspaceID)

	resp, err := s.client.post(ctx, path, opts, nil)
	if err != nil {
		return nil, resp, err
	}

	return resp.Body, resp, nil
}

// ExportWeeklyCSV downloads a weekly report as a CSV file.
//
// API: POST /reports/api/v3/workspace/{workspace_id}/weekly/time_entries.csv
//
// See: https://engineering.toggl.com/docs/reports/weekly
func (s *ReportsService) ExportWeeklyCSV(ctx context.Context, workspaceID int, opts *WeeklyExportOptions) ([]byte, *Response, error) {
	if opts == nil {
		opts = &WeeklyExportOptions{}
	}

	path := fmt.Sprintf("%s/workspace/%d/weekly/time_entries.csv", reportsBasePath, workspaceID)

	resp, err := s.client.post(ctx, path, opts, nil)
	if err != nil {
		return nil, resp, err
	}

	return resp.Body, resp, nil
}
