// Package search builds MusicBrainz Lucene queries for CLI search commands.
package search

import (
	"strings"

	"github.com/liuguancheng/musicbrainz-cli/internal/apperr"
)

const releasePrimaryTypeAlbum = "primarytype:album"

// BuildReleaseQuery composes a release search Lucene query from optional text and artist MBID.
// The query always includes primarytype:album so the API filters to album releases.
// At least one of textQuery or artistMBID must be non-empty.
func BuildReleaseQuery(textQuery, artistMBID string) (string, error) {
	textQuery = strings.TrimSpace(textQuery)
	artistMBID = strings.TrimSpace(artistMBID)

	if textQuery == "" && artistMBID == "" {
		return "", apperr.InvalidArgument("query or --artist-mbid is required")
	}

	if artistMBID != "" {
		if err := apperr.ValidateMBID(artistMBID); err != nil {
			return "", err
		}
	}

	clauses := make([]string, 0, 3)
	if textQuery != "" {
		clauses = append(clauses, "("+textQuery+")")
	}
	if artistMBID != "" {
		clauses = append(clauses, "arid:"+artistMBID)
	}
	clauses = append(clauses, releasePrimaryTypeAlbum)

	return strings.Join(clauses, " AND "), nil
}
