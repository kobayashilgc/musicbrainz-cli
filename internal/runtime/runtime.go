// Package runtime holds process-wide CLI dependencies shared across commands.
package runtime

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/liuguancheng/musicbrainz-cli/internal/client"
	"github.com/liuguancheng/musicbrainz-cli/internal/output"
)

const (
	// DefaultCommandTimeout is the maximum duration for a single CLI command,
	// including all MusicBrainz API calls made during that invocation.
	DefaultCommandTimeout = 60 * time.Second
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

	commandCtx    = context.Background()
	cancelCommand context.CancelFunc
)

// Context returns the root context for API calls in the current command.
func Context() context.Context {
	return commandCtx
}

// StartCommandContext installs a cancelable context for the current command.
// It listens for SIGINT/SIGTERM and applies timeout when timeout > 0.
func StartCommandContext(timeout time.Duration) {
	EndCommandContext()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	if timeout > 0 {
		var cancelTimeout context.CancelFunc
		ctx, cancelTimeout = context.WithTimeout(ctx, timeout)
		prevStop := stop
		stop = func() {
			cancelTimeout()
			prevStop()
		}
	}

	commandCtx = ctx
	cancelCommand = stop
}

// EndCommandContext cancels and clears the active command context.
func EndCommandContext() {
	if cancelCommand != nil {
		cancelCommand()
		cancelCommand = nil
	}
	commandCtx = context.Background()
}

// CloseClient closes the active API client and clears the reference.
func CloseClient() error {
	if Client == nil {
		return nil
	}
	err := Client.Close()
	Client = nil
	return err
}

// ResetIO restores default stdout and stderr writers.
func ResetIO() {
	Stdout = os.Stdout
	Stderr = os.Stderr
}

// ResetForTest clears injected dependencies and returns defaults for tests.
func ResetForTest() {
	EndCommandContext()
	ResetIO()
	Client = nil
	OutputMode = output.ModeSimple
}
