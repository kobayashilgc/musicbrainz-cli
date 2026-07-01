// Package client wraps musicbrainzws2 with CLI-friendly configuration.
package client

import (
	"context"

	"go.uploadedlobster.com/mbtypes"
	"go.uploadedlobster.com/musicbrainzws2"
)

const (
	// DefaultAppName is used when building the MusicBrainz User-Agent.
	DefaultAppName = "musicbrainz-cli"
	// DefaultContactURL is embedded in the User-Agent when --contact is unset.
	DefaultContactURL = "https://github.com/liuguancheng/musicbrainz-cli"
	// DefaultAPIURL is the MusicBrainz Web Service v2 base URL.
	DefaultAPIURL = "https://musicbrainz.org/ws/2/"
)

// Interface abstracts MusicBrainz operations for production and test doubles.
type Interface interface {
	SearchArtists(ctx context.Context, query string, limit, offset int) (musicbrainzws2.SearchArtistsResult, error)
	SearchReleases(ctx context.Context, query string, limit, offset int) (musicbrainzws2.SearchReleasesResult, error)
	LookupArtist(ctx context.Context, mbid mbtypes.MBID, includes []string) (musicbrainzws2.Artist, error)
	LookupRelease(ctx context.Context, mbid mbtypes.MBID, includes []string) (musicbrainzws2.Release, error)
	Close() error
}

// Config holds client identification and endpoint overrides from CLI flags.
type Config struct {
	AppName    string
	Version    string
	ContactURL string
	APIURL     string
	UserAgent  string
}

// Client delegates to an underlying musicbrainzws2 client instance.
type Client struct {
	inner *musicbrainzws2.Client
}

// New builds a MusicBrainz client from CLI configuration.
func New(cfg Config) *Client {
	appInfo := musicbrainzws2.AppInfo{
		Name:    cfg.AppName,
		Version: cfg.Version,
		URL:     cfg.ContactURL,
	}
	if appInfo.Name == "" {
		appInfo.Name = DefaultAppName
	}
	if appInfo.URL == "" {
		appInfo.URL = DefaultContactURL
	}

	inner := musicbrainzws2.NewClient(appInfo)
	if cfg.APIURL != "" {
		inner.SetBaseURL(cfg.APIURL)
	}
	// Explicit --user-agent overrides the value derived from AppInfo.
	if cfg.UserAgent != "" {
		inner.SetUserAgent(cfg.UserAgent)
	}
	return &Client{inner: inner}
}

// Close releases resources held by the underlying HTTP client.
func (c *Client) Close() error {
	return c.inner.Close()
}

// SearchArtists runs a Lucene artist search with explicit pagination.
func (c *Client) SearchArtists(ctx context.Context, query string, limit, offset int) (musicbrainzws2.SearchArtistsResult, error) {
	return c.inner.SearchArtists(ctx, musicbrainzws2.SearchFilter{Query: query}, paginator(limit, offset))
}

// SearchReleases runs a Lucene release search with explicit pagination.
func (c *Client) SearchReleases(ctx context.Context, query string, limit, offset int) (musicbrainzws2.SearchReleasesResult, error) {
	return c.inner.SearchReleases(ctx, musicbrainzws2.SearchFilter{Query: query}, paginator(limit, offset))
}

// LookupArtist fetches a single artist by MBID with optional inc= parameters.
func (c *Client) LookupArtist(ctx context.Context, mbid mbtypes.MBID, includes []string) (musicbrainzws2.Artist, error) {
	return c.inner.LookupArtist(ctx, mbid, musicbrainzws2.IncludesFilter{Includes: includes})
}

// LookupRelease fetches a single release by MBID with optional inc= parameters.
func (c *Client) LookupRelease(ctx context.Context, mbid mbtypes.MBID, includes []string) (musicbrainzws2.Release, error) {
	return c.inner.LookupRelease(ctx, mbid, musicbrainzws2.IncludesFilter{Includes: includes})
}

func paginator(limit, offset int) musicbrainzws2.Paginator {
	return musicbrainzws2.Paginator{
		Limit:  limit,
		Offset: offset,
	}
}
