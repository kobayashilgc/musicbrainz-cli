// Package pagination validates CLI pagination flags against MusicBrainz limits.
package pagination

import (
	"fmt"

	"go.uploadedlobster.com/musicbrainzws2"
)

const (
	// DefaultLimit matches the MusicBrainz search API default page size.
	DefaultLimit = musicbrainzws2.DefaultLimit
	// DefaultOffset is the first page starting index.
	DefaultOffset = 0
	// MinLimit is the smallest allowed page size.
	MinLimit = 1
	// MaxLimit is the largest page size allowed by the MusicBrainz search API.
	MaxLimit = musicbrainzws2.MaxLimit
)

// Validate checks limit and offset against MusicBrainz search constraints.
func Validate(limit, offset int) error {
	if limit < MinLimit || limit > MaxLimit {
		return fmt.Errorf("limit must be between %d and %d", MinLimit, MaxLimit)
	}
	if offset < 0 {
		return fmt.Errorf("offset must be >= 0")
	}
	return nil
}

// NewPaginator builds a musicbrainzws2 paginator from validated CLI values.
func NewPaginator(limit, offset int) musicbrainzws2.Paginator {
	return musicbrainzws2.Paginator{
		Limit:  limit,
		Offset: offset,
	}
}
