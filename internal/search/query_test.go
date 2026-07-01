package search

import (
	"strings"
	"testing"
)

func TestBuildReleaseQueryTextOnly(t *testing.T) {
	t.Parallel()

	query, err := BuildReleaseQuery("Abbey Road", "")
	if err != nil {
		t.Fatalf("BuildReleaseQuery() error = %v", err)
	}
	if query != "(Abbey Road) AND primarytype:album" {
		t.Fatalf("query = %q, want %q", query, "(Abbey Road) AND primarytype:album")
	}
}

func TestBuildReleaseQueryArtistMBIDOnly(t *testing.T) {
	t.Parallel()

	const mbid = "b10bbbfc-cf9e-42e6-888b-88b6b374d5d4"
	query, err := BuildReleaseQuery("", mbid)
	if err != nil {
		t.Fatalf("BuildReleaseQuery() error = %v", err)
	}
	want := "arid:" + mbid + " AND primarytype:album"
	if query != want {
		t.Fatalf("query = %q, want %q", query, want)
	}
}

func TestBuildReleaseQueryCombined(t *testing.T) {
	t.Parallel()

	const mbid = "b10bbbfc-cf9e-42e6-888b-88b6b374d5d4"
	query, err := BuildReleaseQuery(`release:"Abbey Road"`, mbid)
	if err != nil {
		t.Fatalf("BuildReleaseQuery() error = %v", err)
	}
	want := `(release:"Abbey Road") AND arid:` + mbid + " AND primarytype:album"
	if query != want {
		t.Fatalf("query = %q, want %q", query, want)
	}
}

func TestBuildReleaseQueryEmpty(t *testing.T) {
	t.Parallel()

	_, err := BuildReleaseQuery("", "")
	if err == nil {
		t.Fatal("expected error for empty query and artist MBID")
	}
	if !strings.Contains(err.Error(), "query or --artist-mbid is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestBuildReleaseQueryInvalidMBID(t *testing.T) {
	t.Parallel()

	_, err := BuildReleaseQuery("Abbey Road", "not-a-mbid")
	if err == nil {
		t.Fatal("expected error for invalid MBID")
	}
	if !strings.Contains(err.Error(), "invalid mbid") {
		t.Fatalf("error = %v", err)
	}
}

func TestBuildReleaseQueryTrimsWhitespace(t *testing.T) {
	t.Parallel()

	const mbid = "b10bbbfc-cf9e-42e6-888b-88b6b374d5d4"
	query, err := BuildReleaseQuery("  Abbey Road  ", "  "+mbid+"  ")
	if err != nil {
		t.Fatalf("BuildReleaseQuery() error = %v", err)
	}
	want := `(Abbey Road) AND arid:` + mbid + " AND primarytype:album"
	if query != want {
		t.Fatalf("query = %q, want %q", query, want)
	}
}

func TestBuildReleaseGroupQueryMatchesRelease(t *testing.T) {
	t.Parallel()

	const mbid = "b10bbbfc-cf9e-42e6-888b-88b6b374d5d4"
	releaseQuery, err := BuildReleaseQuery("Abbey Road", mbid)
	if err != nil {
		t.Fatalf("BuildReleaseQuery() error = %v", err)
	}
	groupQuery, err := BuildReleaseGroupQuery("Abbey Road", mbid)
	if err != nil {
		t.Fatalf("BuildReleaseGroupQuery() error = %v", err)
	}
	if releaseQuery != groupQuery {
		t.Fatalf("queries differ: release=%q releasegroup=%q", releaseQuery, groupQuery)
	}
}
