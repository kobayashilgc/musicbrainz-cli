// Package pagination validates CLI pagination flags against MusicBrainz limits.
package pagination

import (
	"fmt"

	"go.uploadedlobster.com/musicbrainzws2"
)

const (
	// DefaultLimit matches the MusicBrainz search API default page size.
	DefaultLimit = musicbrainzws2.DefaultLimit
	// DefaultPageNo is the first page number (1-based).
	DefaultPageNo = 1
	// MinPageNo is the smallest allowed page number.
	MinPageNo = 1
	// MinLimit is the smallest allowed page size.
	MinLimit = 1
	// MaxLimit is the largest page size allowed by the MusicBrainz search API.
	MaxLimit = musicbrainzws2.MaxLimit
)

// Validate checks limit and page number against MusicBrainz search constraints.
func Validate(limit, pageNo int) error {
	if limit < MinLimit || limit > MaxLimit {
		return fmt.Errorf("limit must be between %d and %d", MinLimit, MaxLimit)
	}
	if pageNo < MinPageNo {
		return fmt.Errorf("pageno must be >= %d", MinPageNo)
	}
	return nil
}

// Offset converts a 1-based page number and page size to a MusicBrainz API offset.
func Offset(limit, pageNo int) int {
	return (pageNo - 1) * limit
}

// NewPaginator builds a musicbrainzws2 paginator from validated CLI page number and limit.
func NewPaginator(limit, pageNo int) musicbrainzws2.Paginator {
	return musicbrainzws2.Paginator{
		Limit:  limit,
		Offset: Offset(limit, pageNo),
	}
}
