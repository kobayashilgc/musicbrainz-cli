package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"go.uploadedlobster.com/mbtypes"
	"go.uploadedlobster.com/musicbrainzws2"

	"github.com/liuguancheng/musicbrainz-cli/internal/apperr"
)

type mockClient struct {
	mu                 sync.Mutex
	lastReleaseQuery   string
	lastReleaseLimit   int
	lastReleaseOffset  int
}

func (m *mockClient) SearchArtists(_ context.Context, _ string, limit, offset int) (musicbrainzws2.SearchArtistsResult, error) {
	return musicbrainzws2.SearchArtistsResult{
		SearchResult: musicbrainzws2.SearchResult{Count: 1, Offset: offset},
		Artists: []musicbrainzws2.Artist{
			{ID: mbtypes.MBID("b10bbbfc-cf9e-42e6-888b-88b6b374d5d4"), Name: "The Beatles", Score: 100},
		},
	}, nil
}

func (m *mockClient) SearchReleases(_ context.Context, query string, limit, offset int) (musicbrainzws2.SearchReleasesResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastReleaseQuery = query
	m.lastReleaseLimit = limit
	m.lastReleaseOffset = offset
	return musicbrainzws2.SearchReleasesResult{
		SearchResult: musicbrainzws2.SearchResult{Count: 1, Offset: offset},
		Releases: []musicbrainzws2.Release{
			{ID: mbtypes.MBID("464a321e-97a0-4654-8a7a-d1d88e8496e0"), Title: "Abbey Road", Score: 100},
		},
	}, nil
}

func (m *mockClient) releaseQuery() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.lastReleaseQuery
}

func (m *mockClient) LookupArtist(_ context.Context, mbid mbtypes.MBID, _ []string) (musicbrainzws2.Artist, error) {
	return musicbrainzws2.Artist{ID: mbid, Name: "The Beatles"}, nil
}

func (m *mockClient) LookupRelease(_ context.Context, _ mbtypes.MBID, _ []string) (musicbrainzws2.Release, error) {
	return musicbrainzws2.Release{}, nil
}

func (m *mockClient) Close() error { return nil }

func executeWithArgs(t *testing.T, args []string) (int, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	return executeWithClient(t, &mockClient{}, args)
}

func executeWithClient(t *testing.T, client *mockClient, args []string) (int, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()

	var stdout, stderr bytes.Buffer
	SetIO(&stdout, &stderr)
	SetClient(client)

	cmd := RootCmd()
	cmd.SetArgs(args)
	err := cmd.Execute()
	if err != nil {
		return apperr.WriteAndExitCode(&stderr, err), &stdout, &stderr
	}
	return apperr.ExitSuccess, &stdout, &stderr
}

func TestSearchArtistSuccess(t *testing.T) {
	t.Cleanup(ResetForTest)

	code, stdout, stderr := executeWithArgs(t, []string{"search", "artist", "Beatles"})
	if code != apperr.ExitSuccess {
		t.Fatalf("exit code = %d, stderr = %s", code, stderr.String())
	}

	var resp struct {
		Type         string `json:"type"`
		Output       string `json:"output"`
		Query        string `json:"query"`
		Count        int    `json:"count"`
		CurrentCount int    `json:"current_count"`
		HasData      bool   `json:"has_data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal stdout error = %v, stdout = %s", err, stdout.String())
	}
	if resp.Type != "artist_search" || resp.Output != "simple" || resp.Query != "Beatles" {
		t.Fatalf("unexpected response: %#v", resp)
	}
	if resp.Count != 1 || resp.CurrentCount != 1 || !resp.HasData {
		t.Fatalf("unexpected counts: %#v", resp)
	}
}

func TestSearchArtistFullOutput(t *testing.T) {
	t.Cleanup(ResetForTest)

	code, stdout, stderr := executeWithArgs(t, []string{"search", "artist", "Beatles", "--output", "full"})
	if code != apperr.ExitSuccess {
		t.Fatalf("exit code = %d, stderr = %s", code, stderr.String())
	}

	var resp struct {
		Output string         `json:"output"`
		Scores map[string]int `json:"scores"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error = %v", err)
	}
	if resp.Output != "full" {
		t.Fatalf("output = %q", resp.Output)
	}
	if resp.Scores["b10bbbfc-cf9e-42e6-888b-88b6b374d5d4"] != 100 {
		t.Fatalf("scores = %#v", resp.Scores)
	}
}

func TestInvalidOutputMode(t *testing.T) {
	t.Cleanup(ResetForTest)

	code, _, stderr := executeWithArgs(t, []string{"search", "artist", "Beatles", "--output", "table"})
	if code != apperr.ExitInvalidArgument {
		t.Fatalf("exit code = %d, want %d", code, apperr.ExitInvalidArgument)
	}
	if !strings.Contains(stderr.String(), "output must be simple or full") {
		t.Fatalf("stderr = %s", stderr.String())
	}
}

func TestSearchArtistInvalidLimit(t *testing.T) {
	t.Cleanup(ResetForTest)

	code, _, stderr := executeWithArgs(t, []string{"search", "artist", "Beatles", "--limit", "200"})
	if code != apperr.ExitInvalidArgument {
		t.Fatalf("exit code = %d, want %d", code, apperr.ExitInvalidArgument)
	}
	if !strings.Contains(stderr.String(), "INVALID_ARGUMENT") {
		t.Fatalf("stderr = %s", stderr.String())
	}
}

func TestSearchReleaseByArtistMBIDOnly(t *testing.T) {
	t.Cleanup(ResetForTest)

	const mbid = "b10bbbfc-cf9e-42e6-888b-88b6b374d5d4"
	client := &mockClient{}
	code, stdout, stderr := executeWithClient(t, client, []string{"search", "release", "--artist-mbid", mbid})
	if code != apperr.ExitSuccess {
		t.Fatalf("exit code = %d, stderr = %s", code, stderr.String())
	}
	want := "arid:" + mbid + " AND primarytype:album"
	if got := client.releaseQuery(); got != want {
		t.Fatalf("release query = %q, want %q", got, want)
	}

	var resp struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error = %v", err)
	}
	if resp.Query != want {
		t.Fatalf("response query = %q", resp.Query)
	}
}

func TestSearchReleaseWithArtistMBIDAndQuery(t *testing.T) {
	t.Cleanup(ResetForTest)

	const mbid = "b10bbbfc-cf9e-42e6-888b-88b6b374d5d4"
	client := &mockClient{}
	code, stdout, stderr := executeWithClient(t, client, []string{
		"search", "release", "Abbey Road", "--artist-mbid", mbid,
	})
	if code != apperr.ExitSuccess {
		t.Fatalf("exit code = %d, stderr = %s", code, stderr.String())
	}
	want := "(Abbey Road) AND arid:" + mbid + " AND primarytype:album"
	if got := client.releaseQuery(); got != want {
		t.Fatalf("release query = %q, want %q", got, want)
	}

	var resp struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error = %v", err)
	}
	if resp.Query != want {
		t.Fatalf("response query = %q", resp.Query)
	}
}

func TestSearchReleaseRequiresQueryOrArtistMBID(t *testing.T) {
	t.Cleanup(ResetForTest)

	code, _, stderr := executeWithArgs(t, []string{"search", "release"})
	if code != apperr.ExitInvalidArgument {
		t.Fatalf("exit code = %d, want %d", code, apperr.ExitInvalidArgument)
	}
	if !strings.Contains(stderr.String(), "query or --artist-mbid is required") {
		t.Fatalf("stderr = %s", stderr.String())
	}
}

func TestLookupArtistInvalidMBID(t *testing.T) {
	t.Cleanup(ResetForTest)

	code, _, stderr := executeWithArgs(t, []string{"lookup", "artist", "bad-id"})
	if code != apperr.ExitInvalidArgument {
		t.Fatalf("exit code = %d, want %d", code, apperr.ExitInvalidArgument)
	}
	if !strings.Contains(stderr.String(), "invalid mbid") {
		t.Fatalf("stderr = %s", stderr.String())
	}
}

func TestLookupArtistSuccess(t *testing.T) {
	t.Cleanup(ResetForTest)

	code, stdout, stderr := executeWithArgs(t, []string{"lookup", "artist", "b10bbbfc-cf9e-42e6-888b-88b6b374d5d4"})
	if code != apperr.ExitSuccess {
		t.Fatalf("exit code = %d, stderr = %s", code, stderr.String())
	}

	var resp struct {
		Type   string `json:"type"`
		Output string `json:"output"`
		ID     string `json:"id"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error = %v", err)
	}
	if resp.Type != "artist_lookup" || resp.Output != "simple" || resp.ID != "b10bbbfc-cf9e-42e6-888b-88b6b374d5d4" {
		t.Fatalf("unexpected response: %#v", resp)
	}
}

func TestExecuteUsesRuntimeStderr(t *testing.T) {
	t.Cleanup(ResetForTest)

	code, _, stderr := executeWithArgs(t, []string{"search", "artist", "Beatles", "--limit", "0"})
	if code != apperr.ExitInvalidArgument {
		t.Fatalf("exit code = %d, want %d", code, apperr.ExitInvalidArgument)
	}
	if !strings.Contains(stderr.String(), "INVALID_ARGUMENT") {
		t.Fatalf("stderr = %s", stderr.String())
	}
}
