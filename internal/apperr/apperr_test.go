package apperr

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"go.uploadedlobster.com/musicbrainzws2"
)

func TestWriteAndExitCodeInvalidArgument(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	code := WriteAndExitCode(&buf, InvalidArgument("limit must be between 1 and 100"))
	if code != ExitInvalidArgument {
		t.Fatalf("exit code = %d, want %d", code, ExitInvalidArgument)
	}

	var resp struct {
		Code  string `json:"code"`
		Error string `json:"error"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error = %v", err)
	}
	if resp.Code != "INVALID_ARGUMENT" {
		t.Fatalf("code = %q", resp.Code)
	}
}

func TestWriteAndExitCodeClientError(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	err := &musicbrainzws2.ClientError{
		StatusCode: http.StatusNotFound,
		Message:    "not found",
	}
	code := WriteAndExitCode(&buf, err)
	if code != ExitError {
		t.Fatalf("exit code = %d, want %d", code, ExitError)
	}

	var resp struct {
		Code       string `json:"code"`
		StatusCode int    `json:"status_code"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error = %v", err)
	}
	if resp.Code != "NOT_FOUND" || resp.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected response: %#v", resp)
	}
}

func TestValidateMBID(t *testing.T) {
	t.Parallel()

	if err := ValidateMBID(""); err == nil {
		t.Fatal("expected error for empty mbid")
	}
	if err := ValidateMBID("not-a-mbid"); err == nil {
		t.Fatal("expected error for invalid mbid")
	}
	if err := ValidateMBID("b10bbbfc-cf9e-42e6-888b-88b6b374d5d4"); err != nil {
		t.Fatalf("unexpected error for valid mbid: %v", err)
	}
}
