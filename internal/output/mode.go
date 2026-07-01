// Package output builds JSON responses and maps MusicBrainz entities to CLI output.
package output

import (
	"fmt"
	"strings"
)

// Mode selects simplified or full JSON serialization.
type Mode string

const (
	// ModeSimple extracts key fields into compact result objects.
	ModeSimple Mode = "simple"
	// ModeFull emits raw musicbrainzws2 entity structures.
	ModeFull Mode = "full"
)

// ParseMode normalizes and validates a --output flag value.
func ParseMode(value string) (Mode, error) {
	switch Mode(strings.ToLower(strings.TrimSpace(value))) {
	case ModeSimple, "":
		return ModeSimple, nil
	case ModeFull:
		return ModeFull, nil
	default:
		return "", fmt.Errorf("output must be simple or full")
	}
}

func (m Mode) String() string {
	return string(m)
}
