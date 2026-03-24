//go:build integration

package toggl_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	toggl "github.com/shoekstra/go-toggl"
)

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
	time.Sleep(2 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	t.Cleanup(cancel)
	return ctx
}
