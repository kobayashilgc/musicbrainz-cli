package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"go.uploadedlobster.com/mbtypes"
	"go.uploadedlobster.com/musicbrainzws2"
)

func TestArtistSearch(t *testing.T) {
	t.Parallel()

	created := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	result := musicbrainzws2.SearchArtistsResult{
		SearchResult: musicbrainzws2.SearchResult{
			Count:   42,
			Created: created,
		},
		Artists: []musicbrainzws2.Artist{
			{ID: mbtypes.MBID("b10bbbfc-cf9e-42e6-888b-88b6b374d5d4"), Name: "The Beatles", Score: 100},
		},
	}

	resp, err := ArtistSearch(ModeSimple, `artist:"The Beatles"`, 25, 1, result)
	if err != nil {
		t.Fatalf("ArtistSearch() error = %v", err)
	}
	if resp.Type != TypeArtistSearch {
		t.Fatalf("type = %q, want %q", resp.Type, TypeArtistSearch)
	}
	if resp.Output != string(ModeSimple) {
		t.Fatalf("output = %q", resp.Output)
	}
	if resp.Count != 42 {
		t.Fatalf("count = %d, want 42", resp.Count)
	}
	if resp.CurrentCount != 1 {
		t.Fatalf("current_count = %d, want 1", resp.CurrentCount)
	}
	if !resp.HasData {
		t.Fatal("has_data = false, want true")
	}
	if resp.Scores != nil {
		t.Fatalf("scores should be omitted in simple mode")
	}
}

func TestArtistSearchFullIncludesScores(t *testing.T) {
	t.Parallel()

	result := musicbrainzws2.SearchArtistsResult{
		Artists: []musicbrainzws2.Artist{
			{ID: mbtypes.MBID("b10bbbfc-cf9e-42e6-888b-88b6b374d5d4"), Name: "The Beatles", Score: 100},
		},
	}

	resp, err := ArtistSearch(ModeFull, "Beatles", 25, 1, result)
	if err != nil {
		t.Fatalf("ArtistSearch() error = %v", err)
	}
	if resp.Output != string(ModeFull) {
		t.Fatalf("output = %q", resp.Output)
	}
	if resp.Scores["b10bbbfc-cf9e-42e6-888b-88b6b374d5d4"] != 100 {
		t.Fatalf("unexpected scores: %#v", resp.Scores)
	}
}

func TestReleaseSearchUsesPageCount(t *testing.T) {
	t.Parallel()

	result := musicbrainzws2.SearchReleasesResult{
		SearchResult: musicbrainzws2.SearchResult{Count: 100},
		Releases: []musicbrainzws2.Release{
			{ID: mbtypes.MBID("464a321e-97a0-4654-8a7a-d1d88e8496e0"), Title: "Abbey Road", Score: 95},
			{ID: mbtypes.MBID("00000000-0000-0000-0000-000000000001"), Title: "Other", Score: 50},
		},
	}

	resp, err := ReleaseSearch(ModeSimple, "Abbey Road", 10, 1, result)
	if err != nil {
		t.Fatalf("ReleaseSearch() error = %v", err)
	}
	if resp.Count != 100 {
		t.Fatalf("count = %d, want 100", resp.Count)
	}
	if resp.CurrentCount != 2 {
		t.Fatalf("current_count = %d, want 2", resp.CurrentCount)
	}
	if !resp.HasData {
		t.Fatal("has_data = false, want true")
	}
	if resp.MinScore != MinSearchScore {
		t.Fatalf("min_score = %d, want %d", resp.MinScore, MinSearchScore)
	}
	if resp.PrimaryType != RequiredReleasePrimaryType {
		t.Fatalf("primary_type = %q, want %q", resp.PrimaryType, RequiredReleasePrimaryType)
	}
	if resp.PageNo != 1 || resp.Limit != 10 {
		t.Fatalf("pageno/limit = %d/%d, want 1/10", resp.PageNo, resp.Limit)
	}
}

func TestSearchScoreFilter(t *testing.T) {
	t.Parallel()

	result := musicbrainzws2.SearchArtistsResult{
		SearchResult: musicbrainzws2.SearchResult{Count: 3},
		Artists: []musicbrainzws2.Artist{
			{ID: mbtypes.MBID("b10bbbfc-cf9e-42e6-888b-88b6b374d5d4"), Name: "High", Score: 100},
			{ID: mbtypes.MBID("00000000-0000-0000-0000-000000000001"), Name: "Boundary", Score: 50},
			{ID: mbtypes.MBID("00000000-0000-0000-0000-000000000002"), Name: "Low", Score: 49},
		},
	}

	resp, err := ArtistSearch(ModeSimple, "test", 25, 1, result)
	if err != nil {
		t.Fatalf("ArtistSearch() error = %v", err)
	}
	if resp.Count != 3 {
		t.Fatalf("count = %d, want 3", resp.Count)
	}
	if resp.CurrentCount != 2 {
		t.Fatalf("current_count = %d, want 2", resp.CurrentCount)
	}

	var results []map[string]any
	if err := json.Unmarshal(resp.Results, &results); err != nil {
		t.Fatalf("unmarshal results error = %v", err)
	}
	for _, item := range results {
		if mbid, _ := item["mbid"].(string); mbid == "00000000-0000-0000-0000-000000000002" {
			t.Fatalf("score 49 result should be filtered out")
		}
	}
}

func TestArtistSearchBeyondPage(t *testing.T) {
	t.Parallel()

	result := musicbrainzws2.SearchArtistsResult{
		SearchResult: musicbrainzws2.SearchResult{Count: 10},
		Artists:      nil,
	}

	resp, err := ArtistSearch(ModeSimple, "Beatles", 25, 2, result)
	if err != nil {
		t.Fatalf("ArtistSearch() error = %v", err)
	}
	if resp.Count != 10 {
		t.Fatalf("count = %d, want 10", resp.Count)
	}
	if resp.CurrentCount != 0 {
		t.Fatalf("current_count = %d, want 0", resp.CurrentCount)
	}
	if resp.HasData {
		t.Fatal("has_data = true, want false")
	}

	var results []map[string]any
	if err := json.Unmarshal(resp.Results, &results); err != nil {
		t.Fatalf("unmarshal results error = %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("results length = %d, want 0", len(results))
	}
}

func TestArtistSearchInRangeButFilteredEmpty(t *testing.T) {
	t.Parallel()

	result := musicbrainzws2.SearchArtistsResult{
		SearchResult: musicbrainzws2.SearchResult{Count: 10},
		Artists: []musicbrainzws2.Artist{
			{ID: mbtypes.MBID("00000000-0000-0000-0000-000000000002"), Name: "Low", Score: 49},
		},
	}

	resp, err := ArtistSearch(ModeSimple, "test", 25, 1, result)
	if err != nil {
		t.Fatalf("ArtistSearch() error = %v", err)
	}
	if resp.Count != 10 {
		t.Fatalf("count = %d, want 10", resp.Count)
	}
	if resp.CurrentCount != 0 {
		t.Fatalf("current_count = %d, want 0", resp.CurrentCount)
	}
	if !resp.HasData {
		t.Fatal("has_data = false, want true")
	}
}

func TestReleaseSearchBeyondPage(t *testing.T) {
	t.Parallel()

	result := musicbrainzws2.SearchReleasesResult{
		SearchResult: musicbrainzws2.SearchResult{Count: 42},
		Releases:     nil,
	}

	resp, err := ReleaseSearch(ModeSimple, "Abbey Road", 25, 3, result)
	if err != nil {
		t.Fatalf("ReleaseSearch() error = %v", err)
	}
	if resp.Count != 42 {
		t.Fatalf("count = %d, want 42", resp.Count)
	}
	if resp.CurrentCount != 0 {
		t.Fatalf("current_count = %d, want 0", resp.CurrentCount)
	}
	if resp.HasData {
		t.Fatal("has_data = true, want false")
	}
}

func TestWriteJSON(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	err := WriteJSON(&buf, map[string]string{"hello": "world"})
	if err != nil {
		t.Fatalf("WriteJSON() error = %v", err)
	}
	if !strings.Contains(buf.String(), `"hello"`) {
		t.Fatalf("unexpected output: %s", buf.String())
	}
}

func TestWriteError(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := WriteError(&buf, "bad input", CodeInvalidArgument, 0); err != nil {
		t.Fatalf("WriteError() error = %v", err)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error = %v", err)
	}
	if resp.Code != CodeInvalidArgument || resp.Error != "bad input" {
		t.Fatalf("unexpected response: %#v", resp)
	}
}

func TestArtistLookupSimple(t *testing.T) {
	t.Parallel()

	artist := musicbrainzws2.Artist{
		ID:   mbtypes.MBID("b10bbbfc-cf9e-42e6-888b-88b6b374d5d4"),
		Name: "The Beatles",
	}
	resp, err := ArtistLookup(ModeSimple, "b10bbbfc-cf9e-42e6-888b-88b6b374d5d4", artist)
	if err != nil {
		t.Fatalf("ArtistLookup() error = %v", err)
	}
	if resp.Type != TypeArtistLookup {
		t.Fatalf("type = %q, want %q", resp.Type, TypeArtistLookup)
	}
	if resp.Output != string(ModeSimple) {
		t.Fatalf("output = %q", resp.Output)
	}

	var result map[string]any
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		t.Fatalf("unmarshal result error = %v", err)
	}
	if result["artist"] != "The Beatles" {
		t.Fatalf("result = %#v", result)
	}
}
