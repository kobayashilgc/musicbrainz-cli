package output

import (
	"go.uploadedlobster.com/musicbrainzws2"
)

// SimplifiedItem is the compact JSON shape used in simple output mode.
// Fields are omitted when the source entity has no corresponding data.
type SimplifiedItem struct {
	MBID         string   `json:"mbid,omitempty"`
	Score        int      `json:"score,omitempty"`
	Artist       string   `json:"artist,omitempty"`
	Release      string   `json:"release,omitempty"`
	Type         string   `json:"type,omitempty"`
	Country      string   `json:"country,omitempty"`
	Date         string   `json:"date,omitempty"`
	Format       []string `json:"format,omitempty"`
	Barcode      string   `json:"barcode,omitempty"`
	Alias        []string `json:"alias,omitempty"`
	PrimaryAlias string   `json:"primary_alias,omitempty"`
	Tag          []string `json:"tag,omitempty"`
}

// SimplifyArtist maps an artist entity to the compact result schema.
func SimplifyArtist(artist musicbrainzws2.Artist) SimplifiedItem {
	item := SimplifiedItem{
		MBID:         string(artist.ID),
		Score:        artist.Score,
		Artist:       artist.Name,
		Type:         artist.Type,
		Country:      string(artist.CountryCode),
		Alias:        aliasNames(artist.Aliases),
		PrimaryAlias: primaryAlias(artist.Aliases),
		Tag:          tagNames(artist.Tags),
	}
	if date := artist.LifeSpan.Begin.String(); date != "" {
		item.Date = date
	}
	return item
}

// SimplifyRelease maps a release entity to the compact result schema.
func SimplifyRelease(release musicbrainzws2.Release) SimplifiedItem {
	item := SimplifiedItem{
		MBID:         string(release.ID),
		Score:        release.Score,
		Release:      release.Title,
		Artist:       artistCreditString(release.ArtistCredit),
		Country:      string(release.CountryCode),
		Alias:        aliasNames(release.Aliases),
		PrimaryAlias: primaryAlias(release.Aliases),
		Tag:          tagNames(release.Tags),
		Format:       mediaFormats(release.Media),
		Barcode:      string(release.Barcode),
	}
	if release.ReleaseGroup != nil {
		item.Type = release.ReleaseGroup.PrimaryType
	}
	if date := release.Date.String(); date != "" {
		item.Date = date
	}
	return item
}

// SimplifyArtists converts a slice of artists to simplified result items.
func SimplifyArtists(artists []musicbrainzws2.Artist) []SimplifiedItem {
	items := make([]SimplifiedItem, len(artists))
	for i, artist := range artists {
		items[i] = SimplifyArtist(artist)
	}
	return items
}

// SimplifyReleases converts a slice of releases to simplified result items.
func SimplifyReleases(releases []musicbrainzws2.Release) []SimplifiedItem {
	items := make([]SimplifiedItem, len(releases))
	for i, release := range releases {
		items[i] = SimplifyRelease(release)
	}
	return items
}

func aliasNames(aliases []musicbrainzws2.Alias) []string {
	if len(aliases) == 0 {
		return nil
	}
	names := make([]string, 0, len(aliases))
	for _, alias := range aliases {
		if alias.Name != "" {
			names = append(names, alias.Name)
		}
	}
	if len(names) == 0 {
		return nil
	}
	return names
}

// primaryAlias returns the first alias marked as primary in MusicBrainz metadata.
func primaryAlias(aliases []musicbrainzws2.Alias) string {
	for _, alias := range aliases {
		if alias.IsPrimary && alias.Name != "" {
			return alias.Name
		}
	}
	return ""
}

func tagNames(tags []musicbrainzws2.Tag) []string {
	if len(tags) == 0 {
		return nil
	}
	names := make([]string, 0, len(tags))
	for _, tag := range tags {
		if tag.Name != "" {
			names = append(names, tag.Name)
		}
	}
	if len(names) == 0 {
		return nil
	}
	return names
}

func artistCreditString(credit musicbrainzws2.ArtistCredit) string {
	return credit.String()
}

// mediaFormats collects unique medium format strings from a release.
func mediaFormats(media musicbrainzws2.MediaList) []string {
	if len(media) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(media))
	formats := make([]string, 0, len(media))
	for _, medium := range media {
		if medium.Format == "" {
			continue
		}
		if _, ok := seen[medium.Format]; ok {
			continue
		}
		seen[medium.Format] = struct{}{}
		formats = append(formats, medium.Format)
	}
	if len(formats) == 0 {
		return nil
	}
	return formats
}
