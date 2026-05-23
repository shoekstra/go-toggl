//go:build integration

package toggl_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	toggl "github.com/shoekstra/go-toggl"
)

var integrationThrottleMu sync.Mutex

// integrationClient returns a real API client, skipping the test if
// TOGGL_API_TOKEN is not set.
func integrationClient(t *testing.T) *toggl.Client {
	t.Helper()
	token := os.Getenv("TOGGL_API_TOKEN")
	if token == "" {
		t.Skip("TOGGL_API_TOKEN not set")
	}
	client, err := toggl.NewClient(token)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}

// integrationWorkspaceID returns the workspace ID from TOGGL_WORKSPACE_ID,
// skipping the test if the variable is not set.
func integrationWorkspaceID(t *testing.T) int {
	t.Helper()
	s := os.Getenv("TOGGL_WORKSPACE_ID")
	if s == "" {
		t.Skip("TOGGL_WORKSPACE_ID not set")
	}
	id, err := strconv.Atoi(s)
	if err != nil {
		t.Fatalf("TOGGL_WORKSPACE_ID is not a valid integer: %v", err)
	}
	return id
}

// uniqueName returns a name prefixed with "go-toggl-test-" and a nanosecond
// timestamp so parallel test runs and retries don't collide.
func uniqueName(suffix string) string {
	return fmt.Sprintf("go-toggl-test-%s-%d", suffix, time.Now().UnixNano())
}

// integrationCtx returns a context with a per-test deadline. It also sleeps
// briefly to stay within Toggl's API rate limits when tests run sequentially.
func integrationCtx(t *testing.T) context.Context {
	t.Helper()
	integrationThrottleMu.Lock()
	t.Cleanup(integrationThrottleMu.Unlock)
	time.Sleep(3 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func integrationRequireNoError(t *testing.T, op string, err error) {
	t.Helper()
	if err == nil {
		return
	}
	if integrationIsHourlyQuotaError(err) {
		t.Skipf("%s: skipping due to Toggl API hourly quota: %v", op, err)
	}
	t.Fatalf("%s: %v", op, err)
}

func integrationCleanupError(t *testing.T, op string, err error) {
	t.Helper()
	if err == nil {
		return
	}
	if integrationIsHourlyQuotaError(err) {
		t.Logf("%s: cleanup deferred due to Toggl API hourly quota: %v", op, err)
		return
	}
	t.Errorf("%s: %v", op, err)
}

func integrationIsHourlyQuotaError(err error) bool {
	var apiErr *toggl.ErrorResponse
	if !errors.As(err, &apiErr) {
		return false
	}
	if apiErr.StatusCode != http.StatusPaymentRequired && apiErr.StatusCode != http.StatusTooManyRequests {
		return false
	}
	msg := strings.ToLower(apiErr.Message)
	return strings.Contains(msg, "hourly limit") || strings.Contains(msg, "quota will reset")
}
