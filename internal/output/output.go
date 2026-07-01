package output

import (
	"encoding/json"
	"io"
	"time"

	"go.uploadedlobster.com/musicbrainzws2"
)

const (
	TypeArtistSearch  = "artist_search"
	TypeReleaseSearch = "release_search"
	TypeArtistLookup  = "artist_lookup"
	TypeReleaseLookup = "release_lookup"
)

// MinSearchScore is the minimum MusicBrainz search score included in results.
const MinSearchScore = 50

const (
	CodeInvalidArgument = "INVALID_ARGUMENT"
	CodeNotFound        = "NOT_FOUND"
	CodeAPIError        = "API_ERROR"
	CodeInternal        = "INTERNAL"
)

// ErrorResponse is the JSON envelope written to stderr on failure.
type ErrorResponse struct {
	Error      string `json:"error"`
	Code       string `json:"code"`
	StatusCode int    `json:"status_code,omitempty"`
}

// SearchResponse is the JSON envelope for artist and release search results.
type SearchResponse struct {
	Type     string          `json:"type"`
	Output   string          `json:"output"`
	Query    string          `json:"query"`
	Offset   int             `json:"offset"`
	Limit    int             `json:"limit"`
	MinScore int             `json:"min_score"`
	Count    int             `json:"count"`
	Created  time.Time       `json:"created"`
	Results  json.RawMessage `json:"results"`
	Scores   map[string]int  `json:"scores,omitempty"`
}

// LookupResponse is the JSON envelope for artist and release lookup results.
type LookupResponse struct {
	Type   string          `json:"type"`
	Output string          `json:"output"`
	ID     string          `json:"id"`
	Result json.RawMessage `json:"result"`
}

// WriteJSON writes v as indented JSON to w.
func WriteJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// WriteError writes a structured error object to w.
func WriteError(w io.Writer, message, code string, statusCode int) error {
	resp := ErrorResponse{
		Error: message,
		Code:  code,
	}
	if statusCode != 0 {
		resp.StatusCode = statusCode
	}
	return WriteJSON(w, resp)
}

// ArtistSearch builds a search response after score filtering and mode-specific serialization.
func ArtistSearch(mode Mode, query string, limit, offset int, result musicbrainzws2.SearchArtistsResult) (SearchResponse, error) {
	artists := filterArtists(result.Artists)
	created := result.Created
	if created.IsZero() {
		created = time.Now().UTC()
	}

	resp := SearchResponse{
		Type:     TypeArtistSearch,
		Output:   mode.String(),
		Query:    query,
		Offset:   offset,
		Limit:    limit,
		MinScore: MinSearchScore,
		Count:    len(artists),
		Created:  created,
	}

	if mode == ModeFull {
		results, err := json.Marshal(artists)
		if err != nil {
			return SearchResponse{}, err
		}
		scores := make(map[string]int, len(artists))
		for _, artist := range artists {
			scores[string(artist.ID)] = artist.Score
		}
		resp.Results = results
		resp.Scores = scores
		return resp, nil
	}

	simplified := SimplifyArtists(artists)
	results, err := json.Marshal(simplified)
	if err != nil {
		return SearchResponse{}, err
	}
	resp.Results = results
	return resp, nil
}

// ReleaseSearch builds a search response after score filtering and mode-specific serialization.
func ReleaseSearch(mode Mode, query string, limit, offset int, result musicbrainzws2.SearchReleasesResult) (SearchResponse, error) {
	releases := filterReleases(result.Releases)
	created := result.Created
	if created.IsZero() {
		created = time.Now().UTC()
	}

	resp := SearchResponse{
		Type:     TypeReleaseSearch,
		Output:   mode.String(),
		Query:    query,
		Offset:   offset,
		Limit:    limit,
		MinScore: MinSearchScore,
		Count:    len(releases),
		Created:  created,
	}

	if mode == ModeFull {
		results, err := json.Marshal(releases)
		if err != nil {
			return SearchResponse{}, err
		}
		scores := make(map[string]int, len(releases))
		for _, release := range releases {
			scores[string(release.ID)] = release.Score
		}
		resp.Results = results
		resp.Scores = scores
		return resp, nil
	}

	simplified := SimplifyReleases(releases)
	results, err := json.Marshal(simplified)
	if err != nil {
		return SearchResponse{}, err
	}
	resp.Results = results
	return resp, nil
}

// filterArtists drops search hits below MinSearchScore before serialization.
func filterArtists(artists []musicbrainzws2.Artist) []musicbrainzws2.Artist {
	filtered := make([]musicbrainzws2.Artist, 0, len(artists))
	for _, artist := range artists {
		if artist.Score >= MinSearchScore {
			filtered = append(filtered, artist)
		}
	}
	return filtered
}

// filterReleases drops search hits below MinSearchScore before serialization.
func filterReleases(releases []musicbrainzws2.Release) []musicbrainzws2.Release {
	filtered := make([]musicbrainzws2.Release, 0, len(releases))
	for _, release := range releases {
		if release.Score >= MinSearchScore {
			filtered = append(filtered, release)
		}
	}
	return filtered
}

// ArtistLookup builds a lookup response in the selected output mode.
func ArtistLookup(mode Mode, id string, artist musicbrainzws2.Artist) (LookupResponse, error) {
	resp := LookupResponse{
		Type:   TypeArtistLookup,
		Output: mode.String(),
		ID:     id,
	}

	if mode == ModeFull {
		result, err := json.Marshal(artist)
		if err != nil {
			return LookupResponse{}, err
		}
		resp.Result = result
		return resp, nil
	}

	simplified := SimplifyArtist(artist)
	result, err := json.Marshal(simplified)
	if err != nil {
		return LookupResponse{}, err
	}
	resp.Result = result
	return resp, nil
}

// ReleaseLookup builds a lookup response in the selected output mode.
func ReleaseLookup(mode Mode, id string, release musicbrainzws2.Release) (LookupResponse, error) {
	resp := LookupResponse{
		Type:   TypeReleaseLookup,
		Output: mode.String(),
		ID:     id,
	}

	if mode == ModeFull {
		result, err := json.Marshal(release)
		if err != nil {
			return LookupResponse{}, err
		}
		resp.Result = result
		return resp, nil
	}

	simplified := SimplifyRelease(release)
	result, err := json.Marshal(simplified)
	if err != nil {
		return LookupResponse{}, err
	}
	resp.Result = result
	return resp, nil
}
