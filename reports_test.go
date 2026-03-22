package toggl

import (
	"context"
	"net/http"
	"testing"
)

// ── Fixtures ──────────────────────────────────────────────────────────────────

const summaryReportJSON = `{
	"groups": [
		{
			"id": 1,
			"seconds": 3600,
			"rates": [{"billable_seconds": 1800, "currency": "USD", "hourly_rate_in_cents": 10000}],
			"sub_groups": [
				{
					"id": 2,
					"seconds": 1800,
					"rates": [{"billable_seconds": 900, "currency": "USD", "hourly_rate_in_cents": 10000}],
					"title": "My Project"
				}
			]
		}
	]
}`

const summaryReportEmptyJSON = `{"groups": []}`

const detailedTimeEntryJSON = `{
	"billable": true,
	"billable_amount_in_cents": 5000,
	"client_name": "Acme Corp",
	"currency": "USD",
	"description": "Development work",
	"hourly_rate_in_cents": 10000,
	"project_color": "#06aaf5",
	"project_hex": "#06aaf5",
	"project_id": 42,
	"project_name": "Website Redesign",
	"row_number": 1,
	"tag_ids": [1, 2],
	"tag_names": ["backend", "api"],
	"user_id": 7,
	"username": "alice@example.com"
}`

const detailedTimeEntryNullsJSON = `{
	"billable": false,
	"billable_amount_in_cents": null,
	"client_name": null,
	"currency": null,
	"description": null,
	"hourly_rate_in_cents": null,
	"project_color": null,
	"project_hex": null,
	"project_id": null,
	"project_name": null,
	"row_number": 2,
	"tag_ids": [],
	"tag_names": [],
	"user_id": null,
	"username": null
}`

const totalsReportJSON = `{
	"billable_amount_in_cents": 10000,
	"labour_cost_in_cents": 5000,
	"seconds": 7200,
	"tracked_days": 2,
	"resolution": "day",
	"graph": [
		{"billable_amount_in_cents": 5000, "labour_cost_in_cents": 2500, "seconds": 3600}
	],
	"rates": [
		{"billable_seconds": 3600, "currency": "USD", "hourly_rate_in_cents": 10000}
	]
}`

const weeklyReportEntryJSON = `{
	"project_id": 42,
	"client_id": 10,
	"user_id": 7,
	"title": "Website Redesign",
	"seconds": [3600, 0, 1800, 0, 7200, 0, 0]
}`

const weeklyReportEntryNullsJSON = `{
	"project_id": null,
	"client_id": null,
	"user_id": null,
	"title": null,
	"seconds": []
}`

// ── Helpers ───────────────────────────────────────────────────────────────────

// reportHandler returns a handler that asserts POST method, exact path, and
// responds with the given status code and body.
func reportHandler(t *testing.T, wantPath string, statusCode int, body string) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.URL.Path != wantPath {
			t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if body != "" {
			w.Write([]byte(body))
		}
	}
}

// ── Summary Report ────────────────────────────────────────────────────────────

func TestReportsService_SummaryReport(t *testing.T) {
	wantPath := "/reports/api/v3/workspace/100/summary/time_entries"

	tests := []struct {
		name       string
		opts       *SummaryReportOptions
		statusCode int
		response   string
		wantGroups int
		wantErr    bool
	}{
		{
			name:       "success — nil opts",
			opts:       nil,
			statusCode: http.StatusOK,
			response:   summaryReportJSON,
			wantGroups: 1,
			wantErr:    false,
		},
		{
			name:       "success — with opts",
			opts:       &SummaryReportOptions{Grouping: String("projects"), SubGrouping: String("users")},
			statusCode: http.StatusOK,
			response:   summaryReportJSON,
			wantGroups: 1,
			wantErr:    false,
		},
		{
			name:       "empty groups",
			opts:       nil,
			statusCode: http.StatusOK,
			response:   summaryReportEmptyJSON,
			wantGroups: 0,
			wantErr:    false,
		},
		{
			name:       "unauthorized",
			opts:       nil,
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			opts:       nil,
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.URL.Path != wantPath {
					t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if tt.opts != nil {
					body := assertBody(t, r)
					if tt.opts.Grouping != nil {
						if got, ok := body["grouping"].(string); !ok || got != *tt.opts.Grouping {
							t.Errorf("grouping = %v, want %q", body["grouping"], *tt.opts.Grouping)
						}
					}
					if tt.opts.SubGrouping != nil {
						if got, ok := body["sub_grouping"].(string); !ok || got != *tt.opts.SubGrouping {
							t.Errorf("sub_grouping = %v, want %q", body["sub_grouping"], *tt.opts.SubGrouping)
						}
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			client := testClient(t, handler)
			data, _, err := client.Reports.SummaryReport(context.Background(), 100, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("SummaryReport() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(data.Groups) != tt.wantGroups {
				t.Errorf("SummaryReport() groups = %d, want %d", len(data.Groups), tt.wantGroups)
			}
		})
	}
}

func TestReportsService_SummaryReport_BodyFilters(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := assertBody(t, r)
		want := map[string]interface{}{
			"start_date":   "2024-01-01",
			"end_date":     "2024-01-31",
			"grouping":     "projects",
			"sub_grouping": "users",
			"billable":     true,
		}
		for k, v := range want {
			if got := body[k]; got != v {
				t.Errorf("body[%q] = %v, want %v", k, got, v)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(summaryReportEmptyJSON))
	})

	client := testClient(t, handler)
	opts := &SummaryReportOptions{
		ReportFilters: ReportFilters{
			StartDate:  String("2024-01-01"),
			EndDate:    String("2024-01-31"),
			Billable:   Bool(true),
		},
		Grouping:    String("projects"),
		SubGrouping: String("users"),
	}
	if _, _, err := client.Reports.SummaryReport(context.Background(), 100, opts); err != nil {
		t.Fatalf("SummaryReport() unexpected error: %v", err)
	}
}

func TestReportsService_SummaryReport_Fields(t *testing.T) {
	client := testClient(t, reportHandler(t,
		"/reports/api/v3/workspace/100/summary/time_entries",
		http.StatusOK, summaryReportJSON,
	))

	data, _, err := client.Reports.SummaryReport(context.Background(), 100, nil)
	if err != nil {
		t.Fatalf("SummaryReport() error = %v", err)
	}
	if len(data.Groups) != 1 {
		t.Fatalf("Groups count = %d, want 1", len(data.Groups))
	}

	g := data.Groups[0]
	if g.ID == nil || *g.ID != 1 {
		t.Errorf("Groups[0].ID = %v, want 1", g.ID)
	}
	if g.Seconds != 3600 {
		t.Errorf("Groups[0].Seconds = %d, want 3600", g.Seconds)
	}
	if len(g.Rates) != 1 {
		t.Fatalf("Groups[0].Rates count = %d, want 1", len(g.Rates))
	}
	if g.Rates[0].Currency != "USD" {
		t.Errorf("Groups[0].Rates[0].Currency = %q, want %q", g.Rates[0].Currency, "USD")
	}
	if g.Rates[0].HourlyRateCents != 10000 {
		t.Errorf("Groups[0].Rates[0].HourlyRateCents = %d, want 10000", g.Rates[0].HourlyRateCents)
	}
	if len(g.SubGroups) != 1 {
		t.Fatalf("Groups[0].SubGroups count = %d, want 1", len(g.SubGroups))
	}

	sg := g.SubGroups[0]
	if sg.ID == nil || *sg.ID != 2 {
		t.Errorf("SubGroups[0].ID = %v, want 2", sg.ID)
	}
	if sg.Seconds != 1800 {
		t.Errorf("SubGroups[0].Seconds = %d, want 1800", sg.Seconds)
	}
	if sg.Title == nil || *sg.Title != "My Project" {
		t.Errorf("SubGroups[0].Title = %v, want %q", sg.Title, "My Project")
	}
}

// ── Summary Exports ───────────────────────────────────────────────────────────

func TestReportsService_ExportSummaryPDF(t *testing.T) {
	wantPath := "/reports/api/v3/workspace/100/summary/time_entries.pdf"
	pdfBytes := []byte("%PDF-1.4 fake pdf content")

	tests := []struct {
		name       string
		opts       *SummaryExportOptions
		statusCode int
		response   []byte
		wantErr    bool
	}{
		{
			name:       "success — nil opts",
			opts:       nil,
			statusCode: http.StatusOK,
			response:   pdfBytes,
			wantErr:    false,
		},
		{
			name: "success — with opts",
			opts: &SummaryExportOptions{
				DateFormat:     String("YYYY-MM-DD"),
				DurationFormat: String("decimal"),
				OrderBy:        String("duration"),
				HideRates:      Bool(true),
			},
			statusCode: http.StatusOK,
			response:   pdfBytes,
			wantErr:    false,
		},
		{
			name:       "unauthorized",
			opts:       nil,
			statusCode: http.StatusUnauthorized,
			response:   []byte(`{"error": "unauthorized"}`),
			wantErr:    true,
		},
		{
			name:       "server error",
			opts:       nil,
			statusCode: http.StatusInternalServerError,
			response:   []byte(`{"error": "internal server error"}`),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.URL.Path != wantPath {
					t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if tt.opts != nil {
					body := assertBody(t, r)
					if tt.opts.DateFormat != nil {
						if got, ok := body["date_format"].(string); !ok || got != *tt.opts.DateFormat {
							t.Errorf("date_format = %v, want %q", body["date_format"], *tt.opts.DateFormat)
						}
					}
					if tt.opts.DurationFormat != nil {
						if got, ok := body["duration_format"].(string); !ok || got != *tt.opts.DurationFormat {
							t.Errorf("duration_format = %v, want %q", body["duration_format"], *tt.opts.DurationFormat)
						}
					}
					if tt.opts.OrderBy != nil {
						if got, ok := body["order_by"].(string); !ok || got != *tt.opts.OrderBy {
							t.Errorf("order_by = %v, want %q", body["order_by"], *tt.opts.OrderBy)
						}
					}
					if tt.opts.HideRates != nil {
						if got, ok := body["hide_rates"].(bool); !ok || got != *tt.opts.HideRates {
							t.Errorf("hide_rates = %v, want %v", body["hide_rates"], *tt.opts.HideRates)
						}
					}
				}
				w.WriteHeader(tt.statusCode)
				w.Write(tt.response)
			})

			client := testClient(t, handler)
			data, _, err := client.Reports.ExportSummaryPDF(context.Background(), 100, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ExportSummaryPDF() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(data) == 0 {
				t.Errorf("ExportSummaryPDF() returned empty bytes")
			}
		})
	}
}

func TestReportsService_ExportSummaryCSV(t *testing.T) {
	wantPath := "/reports/api/v3/workspace/100/summary/time_entries.csv"
	csvBytes := []byte("project,duration\nWebsite Redesign,3600\n")

	tests := []struct {
		name       string
		statusCode int
		response   []byte
		wantErr    bool
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			response:   csvBytes,
			wantErr:    false,
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			response:   []byte(`{"error": "unauthorized"}`),
			wantErr:    true,
		},
		{
			name:       "server error",
			statusCode: http.StatusInternalServerError,
			response:   []byte(`{"error": "internal server error"}`),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.URL.Path != wantPath {
					t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(tt.statusCode)
				w.Write(tt.response)
			}))

			data, _, err := client.Reports.ExportSummaryCSV(context.Background(), 100, nil)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ExportSummaryCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(data) == 0 {
				t.Errorf("ExportSummaryCSV() returned empty bytes")
			}
		})
	}
}

// ── Detailed Report ───────────────────────────────────────────────────────────

func TestReportsService_DetailedReport(t *testing.T) {
	wantPath := "/reports/api/v3/workspace/100/search/time_entries"

	tests := []struct {
		name       string
		opts       *DetailedReportOptions
		statusCode int
		response   string
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "success — nil opts",
			opts:       nil,
			statusCode: http.StatusOK,
			response:   "[" + detailedTimeEntryJSON + "]",
			wantCount:  1,
			wantErr:    false,
		},
		{
			name:       "success — with opts",
			opts:       &DetailedReportOptions{OrderBy: String("duration"), PageSize: Int(100)},
			statusCode: http.StatusOK,
			response:   "[" + detailedTimeEntryJSON + "," + detailedTimeEntryJSON + "]",
			wantCount:  2,
			wantErr:    false,
		},
		{
			name:       "empty result",
			opts:       nil,
			statusCode: http.StatusOK,
			response:   "[]",
			wantCount:  0,
			wantErr:    false,
		},
		{
			name:       "unauthorized",
			opts:       nil,
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			opts:       nil,
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.URL.Path != wantPath {
					t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if tt.opts != nil {
					body := assertBody(t, r)
					if tt.opts.OrderBy != nil {
						if got, ok := body["order_by"].(string); !ok || got != *tt.opts.OrderBy {
							t.Errorf("order_by = %v, want %q", body["order_by"], *tt.opts.OrderBy)
						}
					}
					if tt.opts.PageSize != nil {
						if got, ok := body["page_size"].(float64); !ok || int(got) != *tt.opts.PageSize {
							t.Errorf("page_size = %v, want %d", body["page_size"], *tt.opts.PageSize)
						}
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			client := testClient(t, handler)
			entries, _, err := client.Reports.DetailedReport(context.Background(), 100, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("DetailedReport() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(entries) != tt.wantCount {
				t.Errorf("DetailedReport() count = %d, want %d", len(entries), tt.wantCount)
			}
		})
	}
}

func TestReportsService_DetailedReport_Fields(t *testing.T) {
	t.Run("populated fields", func(t *testing.T) {
		client := testClient(t, reportHandler(t,
			"/reports/api/v3/workspace/100/search/time_entries",
			http.StatusOK, "["+detailedTimeEntryJSON+"]",
		))

		entries, _, err := client.Reports.DetailedReport(context.Background(), 100, nil)
		if err != nil {
			t.Fatalf("DetailedReport() error = %v", err)
		}
		if len(entries) != 1 {
			t.Fatalf("DetailedReport() count = %d, want 1", len(entries))
		}

		e := entries[0]
		if !e.Billable {
			t.Errorf("Billable = false, want true")
		}
		if e.BillableAmountCents == nil || *e.BillableAmountCents != 5000 {
			t.Errorf("BillableAmountCents = %v, want 5000", e.BillableAmountCents)
		}
		if e.ClientName == nil || *e.ClientName != "Acme Corp" {
			t.Errorf("ClientName = %v, want %q", e.ClientName, "Acme Corp")
		}
		if e.Currency == nil || *e.Currency != "USD" {
			t.Errorf("Currency = %v, want %q", e.Currency, "USD")
		}
		if e.Description == nil || *e.Description != "Development work" {
			t.Errorf("Description = %v, want %q", e.Description, "Development work")
		}
		if e.ProjectID == nil || *e.ProjectID != 42 {
			t.Errorf("ProjectID = %v, want 42", e.ProjectID)
		}
		if e.ProjectName == nil || *e.ProjectName != "Website Redesign" {
			t.Errorf("ProjectName = %v, want %q", e.ProjectName, "Website Redesign")
		}
		if e.RowNumber != 1 {
			t.Errorf("RowNumber = %d, want 1", e.RowNumber)
		}
		if len(e.TagIDs) != 2 {
			t.Errorf("TagIDs count = %d, want 2", len(e.TagIDs))
		}
		if len(e.TagNames) != 2 {
			t.Errorf("TagNames count = %d, want 2", len(e.TagNames))
		}
		if e.UserID == nil || *e.UserID != 7 {
			t.Errorf("UserID = %v, want 7", e.UserID)
		}
		if e.Username == nil || *e.Username != "alice@example.com" {
			t.Errorf("Username = %v, want %q", e.Username, "alice@example.com")
		}
	})

	t.Run("null optional fields", func(t *testing.T) {
		client := testClient(t, reportHandler(t,
			"/reports/api/v3/workspace/100/search/time_entries",
			http.StatusOK, "["+detailedTimeEntryNullsJSON+"]",
		))

		entries, _, err := client.Reports.DetailedReport(context.Background(), 100, nil)
		if err != nil {
			t.Fatalf("DetailedReport() error = %v", err)
		}
		if len(entries) != 1 {
			t.Fatalf("DetailedReport() count = %d, want 1", len(entries))
		}

		e := entries[0]
		if e.Billable {
			t.Errorf("Billable = true, want false")
		}
		if e.BillableAmountCents != nil {
			t.Errorf("BillableAmountCents = %v, want nil", e.BillableAmountCents)
		}
		if e.ClientName != nil {
			t.Errorf("ClientName = %v, want nil", e.ClientName)
		}
		if e.ProjectID != nil {
			t.Errorf("ProjectID = %v, want nil", e.ProjectID)
		}
		if e.UserID != nil {
			t.Errorf("UserID = %v, want nil", e.UserID)
		}
		if e.RowNumber != 2 {
			t.Errorf("RowNumber = %d, want 2", e.RowNumber)
		}
	})
}

// ── Detailed Report Totals ────────────────────────────────────────────────────

func TestReportsService_DetailedReportTotals(t *testing.T) {
	wantPath := "/reports/api/v3/workspace/100/search/time_entries/totals"

	tests := []struct {
		name        string
		opts        *DetailedReportOptions
		statusCode  int
		response    string
		wantSeconds int
		wantErr     bool
	}{
		{
			name:        "success — nil opts",
			opts:        nil,
			statusCode:  http.StatusOK,
			response:    totalsReportJSON,
			wantSeconds: 7200,
			wantErr:     false,
		},
		{
			name: "success — with filters",
			opts: &DetailedReportOptions{
				ReportFilters: ReportFilters{
					StartDate: String("2024-01-01"),
					EndDate:   String("2024-01-31"),
				},
			},
			statusCode:  http.StatusOK,
			response:    totalsReportJSON,
			wantSeconds: 7200,
			wantErr:     false,
		},
		{
			name:       "unauthorized",
			opts:       nil,
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			opts:       nil,
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.URL.Path != wantPath {
					t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if tt.opts != nil {
					body := assertBody(t, r)
					if tt.opts.StartDate != nil {
						if got, ok := body["start_date"].(string); !ok || got != *tt.opts.StartDate {
							t.Errorf("start_date = %v, want %q", body["start_date"], *tt.opts.StartDate)
						}
					}
					if tt.opts.EndDate != nil {
						if got, ok := body["end_date"].(string); !ok || got != *tt.opts.EndDate {
							t.Errorf("end_date = %v, want %q", body["end_date"], *tt.opts.EndDate)
						}
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			client := testClient(t, handler)
			totals, _, err := client.Reports.DetailedReportTotals(context.Background(), 100, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("DetailedReportTotals() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && totals.Seconds != tt.wantSeconds {
				t.Errorf("DetailedReportTotals() seconds = %d, want %d", totals.Seconds, tt.wantSeconds)
			}
		})
	}
}

func TestReportsService_DetailedReportTotals_Fields(t *testing.T) {
	client := testClient(t, reportHandler(t,
		"/reports/api/v3/workspace/100/search/time_entries/totals",
		http.StatusOK, totalsReportJSON,
	))

	totals, _, err := client.Reports.DetailedReportTotals(context.Background(), 100, nil)
	if err != nil {
		t.Fatalf("DetailedReportTotals() error = %v", err)
	}

	if totals.BillableAmountCents != 10000 {
		t.Errorf("BillableAmountCents = %d, want 10000", totals.BillableAmountCents)
	}
	if totals.LabourCostCents != 5000 {
		t.Errorf("LabourCostCents = %d, want 5000", totals.LabourCostCents)
	}
	if totals.Seconds != 7200 {
		t.Errorf("Seconds = %d, want 7200", totals.Seconds)
	}
	if totals.TrackedDays != 2 {
		t.Errorf("TrackedDays = %d, want 2", totals.TrackedDays)
	}
	if totals.Resolution != "day" {
		t.Errorf("Resolution = %q, want %q", totals.Resolution, "day")
	}
	if len(totals.Graph) != 1 {
		t.Fatalf("Graph count = %d, want 1", len(totals.Graph))
	}
	if totals.Graph[0].Seconds != 3600 {
		t.Errorf("Graph[0].Seconds = %d, want 3600", totals.Graph[0].Seconds)
	}
	if len(totals.Rates) != 1 {
		t.Fatalf("Rates count = %d, want 1", len(totals.Rates))
	}
	if totals.Rates[0].Currency != "USD" {
		t.Errorf("Rates[0].Currency = %q, want %q", totals.Rates[0].Currency, "USD")
	}
	if totals.Rates[0].HourlyRateCents != 10000 {
		t.Errorf("Rates[0].HourlyRateCents = %d, want 10000", totals.Rates[0].HourlyRateCents)
	}
}

// ── Detailed Exports ──────────────────────────────────────────────────────────

func TestReportsService_ExportDetailedPDF(t *testing.T) {
	wantPath := "/reports/api/v3/workspace/100/search/time_entries.pdf"
	pdfBytes := []byte("%PDF-1.4 fake content")

	tests := []struct {
		name       string
		opts       *DetailedExportOptions
		statusCode int
		response   []byte
		wantErr    bool
	}{
		{
			name:       "success — nil opts",
			opts:       nil,
			statusCode: http.StatusOK,
			response:   pdfBytes,
			wantErr:    false,
		},
		{
			name: "success — with opts",
			opts: &DetailedExportOptions{
				DateFormat:     String("YYYY-MM-DD"),
				DurationFormat: String("decimal"),
				DisplayMode:    String("date_and_time"),
			},
			statusCode: http.StatusOK,
			response:   pdfBytes,
			wantErr:    false,
		},
		{
			name:       "unauthorized",
			opts:       nil,
			statusCode: http.StatusUnauthorized,
			response:   []byte(`{"error": "unauthorized"}`),
			wantErr:    true,
		},
		{
			name:       "server error",
			opts:       nil,
			statusCode: http.StatusInternalServerError,
			response:   []byte(`{"error": "internal server error"}`),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.URL.Path != wantPath {
					t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if tt.opts != nil {
					body := assertBody(t, r)
					if tt.opts.DateFormat != nil {
						if got, ok := body["date_format"].(string); !ok || got != *tt.opts.DateFormat {
							t.Errorf("date_format = %v, want %q", body["date_format"], *tt.opts.DateFormat)
						}
					}
					if tt.opts.DurationFormat != nil {
						if got, ok := body["duration_format"].(string); !ok || got != *tt.opts.DurationFormat {
							t.Errorf("duration_format = %v, want %q", body["duration_format"], *tt.opts.DurationFormat)
						}
					}
					if tt.opts.DisplayMode != nil {
						if got, ok := body["display_mode"].(string); !ok || got != *tt.opts.DisplayMode {
							t.Errorf("display_mode = %v, want %q", body["display_mode"], *tt.opts.DisplayMode)
						}
					}
				}
				w.WriteHeader(tt.statusCode)
				w.Write(tt.response)
			})

			client := testClient(t, handler)
			data, _, err := client.Reports.ExportDetailedPDF(context.Background(), 100, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ExportDetailedPDF() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(data) == 0 {
				t.Errorf("ExportDetailedPDF() returned empty bytes")
			}
		})
	}
}

func TestReportsService_ExportDetailedCSV(t *testing.T) {
	wantPath := "/reports/api/v3/workspace/100/search/time_entries.csv"
	csvBytes := []byte("date,description,duration\n2024-01-15,Development work,3600\n")

	tests := []struct {
		name       string
		statusCode int
		response   []byte
		wantErr    bool
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			response:   csvBytes,
			wantErr:    false,
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			response:   []byte(`{"error": "unauthorized"}`),
			wantErr:    true,
		},
		{
			name:       "server error",
			statusCode: http.StatusInternalServerError,
			response:   []byte(`{"error": "internal server error"}`),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.URL.Path != wantPath {
					t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(tt.statusCode)
				w.Write(tt.response)
			}))

			data, _, err := client.Reports.ExportDetailedCSV(context.Background(), 100, nil)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ExportDetailedCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(data) == 0 {
				t.Errorf("ExportDetailedCSV() returned empty bytes")
			}
		})
	}
}

// ── Weekly Report ─────────────────────────────────────────────────────────────

func TestReportsService_WeeklyReport(t *testing.T) {
	wantPath := "/reports/api/v3/workspace/100/weekly/time_entries"

	tests := []struct {
		name       string
		opts       *WeeklyReportOptions
		statusCode int
		response   string
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "success — nil opts",
			opts:       nil,
			statusCode: http.StatusOK,
			response:   "[" + weeklyReportEntryJSON + "]",
			wantCount:  1,
			wantErr:    false,
		},
		{
			name: "success — with filters",
			opts: &WeeklyReportOptions{
				ReportFilters: ReportFilters{
					StartDate: String("2024-01-01"),
					EndDate:   String("2024-01-07"),
				},
			},
			statusCode: http.StatusOK,
			response:   "[" + weeklyReportEntryJSON + "]",
			wantCount:  1,
			wantErr:    false,
		},
		{
			name:       "empty result",
			opts:       nil,
			statusCode: http.StatusOK,
			response:   "[]",
			wantCount:  0,
			wantErr:    false,
		},
		{
			name:       "unauthorized",
			opts:       nil,
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			opts:       nil,
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.URL.Path != wantPath {
					t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if tt.opts != nil {
					body := assertBody(t, r)
					if tt.opts.StartDate != nil {
						if got, ok := body["start_date"].(string); !ok || got != *tt.opts.StartDate {
							t.Errorf("start_date = %v, want %q", body["start_date"], *tt.opts.StartDate)
						}
					}
					if tt.opts.EndDate != nil {
						if got, ok := body["end_date"].(string); !ok || got != *tt.opts.EndDate {
							t.Errorf("end_date = %v, want %q", body["end_date"], *tt.opts.EndDate)
						}
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			client := testClient(t, handler)
			entries, _, err := client.Reports.WeeklyReport(context.Background(), 100, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("WeeklyReport() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(entries) != tt.wantCount {
				t.Errorf("WeeklyReport() count = %d, want %d", len(entries), tt.wantCount)
			}
		})
	}
}

func TestReportsService_WeeklyReport_Fields(t *testing.T) {
	t.Run("populated fields", func(t *testing.T) {
		client := testClient(t, reportHandler(t,
			"/reports/api/v3/workspace/100/weekly/time_entries",
			http.StatusOK, "["+weeklyReportEntryJSON+"]",
		))

		entries, _, err := client.Reports.WeeklyReport(context.Background(), 100, nil)
		if err != nil {
			t.Fatalf("WeeklyReport() error = %v", err)
		}
		if len(entries) != 1 {
			t.Fatalf("WeeklyReport() count = %d, want 1", len(entries))
		}

		e := entries[0]
		if e.ProjectID == nil || *e.ProjectID != 42 {
			t.Errorf("ProjectID = %v, want 42", e.ProjectID)
		}
		if e.ClientID == nil || *e.ClientID != 10 {
			t.Errorf("ClientID = %v, want 10", e.ClientID)
		}
		if e.UserID == nil || *e.UserID != 7 {
			t.Errorf("UserID = %v, want 7", e.UserID)
		}
		if e.Title == nil || *e.Title != "Website Redesign" {
			t.Errorf("Title = %v, want %q", e.Title, "Website Redesign")
		}
		if len(e.Seconds) != 7 {
			t.Errorf("Seconds count = %d, want 7", len(e.Seconds))
		}
		if e.Seconds[0] != 3600 {
			t.Errorf("Seconds[0] = %d, want 3600", e.Seconds[0])
		}
	})

	t.Run("null optional fields", func(t *testing.T) {
		client := testClient(t, reportHandler(t,
			"/reports/api/v3/workspace/100/weekly/time_entries",
			http.StatusOK, "["+weeklyReportEntryNullsJSON+"]",
		))

		entries, _, err := client.Reports.WeeklyReport(context.Background(), 100, nil)
		if err != nil {
			t.Fatalf("WeeklyReport() error = %v", err)
		}
		if len(entries) != 1 {
			t.Fatalf("WeeklyReport() count = %d, want 1", len(entries))
		}

		e := entries[0]
		if e.ProjectID != nil {
			t.Errorf("ProjectID = %v, want nil", e.ProjectID)
		}
		if e.ClientID != nil {
			t.Errorf("ClientID = %v, want nil", e.ClientID)
		}
		if e.UserID != nil {
			t.Errorf("UserID = %v, want nil", e.UserID)
		}
		if e.Title != nil {
			t.Errorf("Title = %v, want nil", e.Title)
		}
	})
}

// ── Weekly Exports ────────────────────────────────────────────────────────────

func TestReportsService_ExportWeeklyPDF(t *testing.T) {
	wantPath := "/reports/api/v3/workspace/100/weekly/time_entries.pdf"

	tests := []struct {
		name       string
		opts       *WeeklyExportOptions
		statusCode int
		response   []byte
		wantErr    bool
	}{
		{
			name:       "success — nil opts",
			opts:       nil,
			statusCode: http.StatusOK,
			response:   []byte("%PDF-1.4 fake content"),
			wantErr:    false,
		},
		{
			name: "success — with opts",
			opts: &WeeklyExportOptions{
				GroupByTask:    Bool(true),
				DurationFormat: String("decimal"),
				DateFormat:     String("YYYY-MM-DD"),
			},
			statusCode: http.StatusOK,
			response:   []byte("%PDF-1.4 fake content"),
			wantErr:    false,
		},
		{
			name:       "unauthorized",
			opts:       nil,
			statusCode: http.StatusUnauthorized,
			response:   []byte(`{"error": "unauthorized"}`),
			wantErr:    true,
		},
		{
			name:       "server error",
			opts:       nil,
			statusCode: http.StatusInternalServerError,
			response:   []byte(`{"error": "internal server error"}`),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.URL.Path != wantPath {
					t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if tt.opts != nil {
					body := assertBody(t, r)
					if tt.opts.GroupByTask != nil {
						if got, ok := body["group_by_task"].(bool); !ok || got != *tt.opts.GroupByTask {
							t.Errorf("group_by_task = %v, want %v", body["group_by_task"], *tt.opts.GroupByTask)
						}
					}
					if tt.opts.DurationFormat != nil {
						if got, ok := body["duration_format"].(string); !ok || got != *tt.opts.DurationFormat {
							t.Errorf("duration_format = %v, want %q", body["duration_format"], *tt.opts.DurationFormat)
						}
					}
					if tt.opts.DateFormat != nil {
						if got, ok := body["date_format"].(string); !ok || got != *tt.opts.DateFormat {
							t.Errorf("date_format = %v, want %q", body["date_format"], *tt.opts.DateFormat)
						}
					}
				}
				w.WriteHeader(tt.statusCode)
				w.Write(tt.response)
			})

			client := testClient(t, handler)
			data, _, err := client.Reports.ExportWeeklyPDF(context.Background(), 100, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ExportWeeklyPDF() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(data) == 0 {
				t.Errorf("ExportWeeklyPDF() returned empty bytes")
			}
		})
	}
}

func TestReportsService_ExportWeeklyCSV(t *testing.T) {
	wantPath := "/reports/api/v3/workspace/100/weekly/time_entries.csv"
	csvBytes := []byte("project,mon,tue,wed,thu,fri,sat,sun,total\nWebsite Redesign,3600,0,1800,0,7200,0,0,12600\n")

	tests := []struct {
		name       string
		statusCode int
		response   []byte
		wantErr    bool
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			response:   csvBytes,
			wantErr:    false,
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			response:   []byte(`{"error": "unauthorized"}`),
			wantErr:    true,
		},
		{
			name:       "server error",
			statusCode: http.StatusInternalServerError,
			response:   []byte(`{"error": "internal server error"}`),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.URL.Path != wantPath {
					t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(tt.statusCode)
				w.Write(tt.response)
			}))

			data, _, err := client.Reports.ExportWeeklyCSV(context.Background(), 100, nil)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ExportWeeklyCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(data) == 0 {
				t.Errorf("ExportWeeklyCSV() returned empty bytes")
			}
		})
	}
}
