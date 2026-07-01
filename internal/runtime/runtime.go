// Package runtime holds process-wide CLI dependencies shared across commands.
package runtime

import (
	"context"
	"io"
	"os"

	"github.com/liuguancheng/musicbrainz-cli/internal/client"
	"github.com/liuguancheng/musicbrainz-cli/internal/output"
)

var (
	// Client is the active MusicBrainz API client for the current command.
	Client client.Interface
	// OutputMode selects simplified or full JSON serialization.
	OutputMode output.Mode = output.ModeSimple
	// Stdout receives successful JSON responses.
	Stdout io.Writer = os.Stdout
	// Stderr receives JSON error responses.
	Stderr io.Writer = os.Stderr
)

// Context returns the root context for API calls in this process.
func Context() context.Context {
	return context.Background()
}

// ResetIO restores default stdout and stderr writers.
func ResetIO() {
	Stdout = os.Stdout
	Stderr = os.Stderr
}

// ResetForTest clears injected dependencies and returns defaults for tests.
func ResetForTest() {
	ResetIO()
	Client = nil
	OutputMode = output.ModeSimple
}
