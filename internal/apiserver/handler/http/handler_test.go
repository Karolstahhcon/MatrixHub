package backend_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matrixhub-ai/hfd/pkg/permission"
	"github.com/matrixhub-ai/hfd/pkg/storage"

	backendhttp "github.com/matrixhub-ai/matrixhub/internal/apiserver/handler/http"
)

// runGitCmd runs a git command in the specified directory.
func runGitCmd(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := Command(t.Context(), "git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Git command failed: git %s\nError: %v\nOutput: %s", strings.Join(args, " "), err, output)
	}
	return string(output)
}

// Command creates an exec.Cmd with the given name and arguments, and logs the command being executed.
func Command(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stderr = os.Stderr
	return cmd
}

func TestHTTPHandler(t *testing.T) {
	// Create a temporary directory for the upstream server
	upstreamDir, err := os.MkdirTemp("", "http-test-upstream")
	if err != nil {
		t.Fatalf("Failed to create temp upstream dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(upstreamDir)
	}()

	// Create a temporary directory for client operations
	clientDir, err := os.MkdirTemp("", "http-test-client")
