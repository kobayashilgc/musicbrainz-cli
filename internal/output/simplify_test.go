package output

import (
	"encoding/json"
	"strings"
	"testing"

	"go.uploadedlobster.com/mbtypes"
	"go.uploadedlobster.com/musicbrainzws2"
)

func TestParseMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input   string
		want    Mode
		wantErr bool
	}{
		{input: "simple", want: ModeSimple},
		{input: "", want: ModeSimple},
		{input: "full", want: ModeFull},
		{input: "FULL", want: ModeFull},
		{input: "invalid", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got, err := ParseMode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseMode(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("ParseMode(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSimplifyArtistOmitsEmptyFields(t *testing.T) {
	t.Parallel()

	item := SimplifyArtist(musicbrainzws2.Artist{
		ID:    mbtypes.MBID("b10bbbfc-cf9e-42e6-888b-88b6b374d5d4"),
		Name:  "The Beatles",
		Score: 100,
		Type:  "Group",
	})

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("marshal error = %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal error = %v", err)
	}

	for _, key := range []string{"release", "format", "barcode", "primary_alias"} {
		if _, ok := raw[key]; ok {
			t.Fatalf("unexpected field %q in simplified artist: %s", key, string(data))
		}
	}
	if raw["artist"] != "The Beatles" {
		t.Fatalf("artist = %v", raw["artist"])
	}
}

func TestSimplifyReleaseFields(t *testing.T) {
	t.Parallel()

	release := musicbrainzws2.Release{
		ID:    mbtypes.MBID("464a321e-97a0-4654-8a7a-d1d88e8496e0"),
		Title: "Abbey Road",
		Score: 95,
		ArtistCredit: musicbrainzws2.ArtistCredit{
			{Name: "The Beatles", Artist: musicbrainzws2.Artist{Name: "The Beatles"}},
		},
		ReleaseGroup: &musicbrainzws2.ReleaseGroup{PrimaryType: "Album"},
		Barcode:      mbtypes.Barcode("602527306377"),
		Media: musicbrainzws2.MediaList{
			{Format: "CD"},
			{Format: "CD"},
			{Format: "Vinyl"},
		},
		Aliases: []musicbrainzws2.Alias{
			{Name: "别名", IsPrimary: true},
			{Name: "Other"},
		},
		Tags: []musicbrainzws2.Tag{{Name: "rock"}},
	}
	release.Date.Parse("1969-09-26")

	item := SimplifyRelease(release)
	if item.Release != "Abbey Road" {
		t.Fatalf("release = %q", item.Release)
	}
	if item.Artist != "The Beatles" {
		t.Fatalf("artist = %q", item.Artist)
	}
	if item.Type != "Album" {
		t.Fatalf("type = %q", item.Type)
	}
	if item.PrimaryAlias != "别名" {
		t.Fatalf("primary_alias = %q", item.PrimaryAlias)
	}
	if len(item.Format) != 2 {
		t.Fatalf("format = %#v, want 2 unique values", item.Format)
	}
	if item.Barcode != "602527306377" {
		t.Fatalf("barcode = %q", item.Barcode)
	}
	if len(item.Tag) != 1 || item.Tag[0] != "rock" {
		t.Fatalf("tag = %#v", item.Tag)
	}
}

func TestSimplifyReleaseGroupFields(t *testing.T) {
	t.Parallel()

	releaseGroup := musicbrainzws2.ReleaseGroup{
		ID:          mbtypes.MBID("abbc4905-c25f-4c67-8e2d-19329ec48b1f"),
		Title:       "Abbey Road",
		Score:       95,
		PrimaryType: "Album",
		ArtistCredit: musicbrainzws2.ArtistCredit{
			{Name: "The Beatles", Artist: musicbrainzws2.Artist{Name: "The Beatles"}},
		},
		Aliases: []musicbrainzws2.Alias{
			{Name: "别名", IsPrimary: true},
		},
		Tags: []musicbrainzws2.Tag{{Name: "rock"}},
	}
	releaseGroup.FirstReleaseDate.Parse("1969-09-26")

	item := SimplifyReleaseGroup(releaseGroup)
	if item.ReleaseGroup != "Abbey Road" {
		t.Fatalf("releasegroup = %q", item.ReleaseGroup)
	}
	if item.Release != "" {
		t.Fatalf("release should be omitted, got %q", item.Release)
	}
	if item.Artist != "The Beatles" {
		t.Fatalf("artist = %q", item.Artist)
	}
	if item.Type != "Album" {
		t.Fatalf("type = %q", item.Type)
	}
	if item.Date != "1969-09-26" {
		t.Fatalf("date = %q", item.Date)
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("marshal error = %v", err)
	}
	if !strings.Contains(string(data), `"releasegroup":"Abbey Road"`) {
		t.Fatalf("unexpected json: %s", string(data))
	}
}
