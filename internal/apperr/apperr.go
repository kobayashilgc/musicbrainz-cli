// Package apperr maps errors to JSON stderr output and process exit codes.
package apperr

import (
	"errors"
	"io"
	"net/http"

	"go.uploadedlobster.com/mbtypes"
	"go.uploadedlobster.com/musicbrainzws2"

	"github.com/liuguancheng/musicbrainz-cli/internal/output"
)

const (
	// ExitSuccess indicates the command completed without error.
	ExitSuccess = 0
	// ExitError indicates a runtime or API failure.
	ExitError = 1
	// ExitInvalidArgument indicates invalid CLI input.
	ExitInvalidArgument = 2
)

// WriteAndExitCode serializes err to w as JSON and returns the matching exit code.
func WriteAndExitCode(w io.Writer, err error) int {
	if err == nil {
		return ExitSuccess
	}

	var invalid InvalidArgumentError
	if errors.As(err, &invalid) {
		_ = output.WriteError(w, invalid.Error(), output.CodeInvalidArgument, 0)
		return ExitInvalidArgument
	}

	var clientErr *musicbrainzws2.ClientError
	if errors.As(err, &clientErr) {
		code := output.CodeAPIError
		if clientErr.StatusCode == http.StatusNotFound {
			code = output.CodeNotFound
		}
		_ = output.WriteError(w, clientErr.Error(), code, clientErr.StatusCode)
		return ExitError
	}

	_ = output.WriteError(w, err.Error(), output.CodeInternal, 0)
	return ExitError
}

// InvalidArgumentError marks CLI validation failures handled as exit code 2.
type InvalidArgumentError struct {
	msg string
}

func (e InvalidArgumentError) Error() string {
	return e.msg
}

// InvalidArgument wraps a validation message as InvalidArgumentError.
func InvalidArgument(msg string) error {
	return InvalidArgumentError{msg: msg}
}

// ValidateMBID ensures the argument is a non-empty, well-formed MusicBrainz ID.
func ValidateMBID(id string) error {
	if id == "" {
		return InvalidArgument("mbid is required")
	}
	if !mbtypes.MBID(id).IsValid() {
		return InvalidArgument("invalid mbid format")
	}
	return nil
}
