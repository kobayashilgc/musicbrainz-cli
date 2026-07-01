# musicbrainz-cli

<p align="center">
  <img src="doc/mbz.png" alt="mbz" height="200" />
  &nbsp;&nbsp;
  <img src="doc/cobra.png" alt="cobra" height="200" />
</p>

[English](./README.md) | [简体中文](./README.zh-CN.md)

A command-line tool for querying [MusicBrainz Web Service v2](https://musicbrainz.org/doc/MusicBrainz_API). Search and look up artist, release, and release group metadata with JSON output designed for scripting and pipelines.

## Features

- **Search** artists, releases, and release groups using Lucene query syntax
- **Release / release group search by artist MBID** — filter with `--artist-mbid` on `search release` and `search releasegroup`
- **Album-only release & release group search** — both append `primarytype:album` to the API query automatically
- **Lookup** artists, releases, and release groups by MusicBrainz ID (MBID)
- **Pagination** via `--limit` and `--pageno` (search only)
- **Score filtering** — search results with score &lt; 50 are dropped automatically
- **JSON output** — success on stdout, errors on stderr
- **Output modes** — `simple` (default, key fields only) or `full` (raw API structures)

## Requirements

- Go 1.24 or later
- Network access to `https://musicbrainz.org`

## Installation

Install from source:

```bash
git clone https://github.com/liuguancheng/musicbrainz-cli.git
cd musicbrainz-cli
go install .
```

Or build a local binary:

```bash
go build -o mbz .
```

The binary name is `mbz`.

## Quick start

```bash
# Search artists (simplified JSON by default)
mbz search artist "Beatles"

# Search releases with pagination
mbz search release 'release:"Abbey Road" AND artist:"Beatles"' --limit 5

# Look up by MBID
mbz lookup artist b10bbbfc-cf9e-42e6-888b-88b6b374d5d4

# Pipe into jq
mbz search artist "Beatles" | jq '.results[].artist'
```

## Command reference

```
mbz
├── search
│   ├── artist <query>        Search artists
│   ├── release [query]       Search releases (optional query; use --artist-mbid)
│   └── releasegroup [query]  Search release groups (optional query; use --artist-mbid)
└── lookup
    ├── artist <mbid>         Look up an artist
    ├── release <mbid>        Look up a release
    └── releasegroup <mbid>   Look up a release group
```

### Global flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--limit` | `-l` | `25` | Page size (1–100); search only |
| `--pageno` | `-p` | `1` | Page number (≥ 1); search only |
| `--output` | | `simple` | Output mode: `simple` or `full` |
| `--user-agent` | | auto | HTTP User-Agent override |
| `--contact` | | repo URL | Contact URL embedded in User-Agent |
| `--api-url` | | `https://musicbrainz.org/ws/2/` | MusicBrainz WS2 base URL |

### Lookup flags

| Flag | Description |
|------|-------------|
| `--inc` | Include related data (repeatable), e.g. `releases`, `artist-credits`, `media` |

### Search release / releasegroup flags

| Flag | Description |
|------|-------------|
| `--artist-mbid` | Filter by artist MBID (`arid`); optional `query` arg can be combined with `AND` |

Provide at least one of `[query]` or `--artist-mbid` for `search release` and `search releasegroup`.

## Examples

### Search artists

```bash
mbz search artist "Beatles"
mbz search artist 'artist:"The Beatles"' --limit 10 --pageno 1
mbz search artist 'artist:"The Beatles"' --limit 10 --pageno 2   # next page
mbz search artist "Beatles" --output full                         # full API JSON
```

### Search releases

```bash
mbz search release 'release:"Abbey Road" AND artist:"Beatles"'
mbz search release "Abbey Road" --limit 5

# All releases for an artist (by artist MBID)
mbz search release --artist-mbid b10bbbfc-cf9e-42e6-888b-88b6b374d5d4

# Text query + artist MBID filter
mbz search release "Abbey Road" --artist-mbid b10bbbfc-cf9e-42e6-888b-88b6b374d5d4
```

### Search release groups

```bash
mbz search releasegroup "Abbey Road"
mbz search releasegroup --artist-mbid b10bbbfc-cf9e-42e6-888b-88b6b374d5d4
mbz search releasegroup "Abbey Road" --artist-mbid b10bbbfc-cf9e-42e6-888b-88b6b374d5d4
```

### Lookup

```bash
mbz lookup artist b10bbbfc-cf9e-42e6-888b-88b6b374d5d4
mbz lookup release 464a321e-97a0-4654-8a7a-d1d88e8496e0 --inc artist-credits --inc media
mbz lookup releasegroup abbc4905-c25f-4c67-8e2d-19329ec48b1f
```

## Output format

### Modes

| Mode | Flag | Description |
|------|------|-------------|
| `simple` | `--output simple` (default) | Extracts key fields; omits missing data |
| `full` | `--output full` | Full MusicBrainz API entity structures |

**Simple mode fields** (included only when present in the API response):

`mbid`, `score`, `artist`, `release`, `releasegroup`, `type`, `country`, `date`, `format`, `barcode`, `alias`, `primary_alias`, `tag`

### Search response (simple)

Artist search example:

```json
{
  "type": "artist_search",
  "output": "simple",
  "query": "Beatles",
  "pageno": 1,
  "limit": 25,
  "min_score": 50,
  "count": 42,
  "current_count": 1,
  "has_data": true,
  "created": "2026-07-01T12:00:00Z",
  "results": [
    {
      "mbid": "b10bbbfc-cf9e-42e6-888b-88b6b374d5d4",
      "score": 100,
      "artist": "The Beatles",
      "type": "Group",
      "country": "GB"
    }
  ]
}
```

Release search adds `"primary_type": "album"` (API-side filter via Lucene `primarytype:album`):

```json
{
  "type": "release_search",
  "output": "simple",
  "query": "(Abbey Road) AND primarytype:album",
  "pageno": 1,
  "limit": 25,
  "min_score": 50,
  "primary_type": "album",
  "count": 42,
  "current_count": 1,
  "has_data": true,
  "created": "2026-07-01T12:00:00Z",
  "results": [
    {
      "mbid": "464a321e-97a0-4654-8a7a-d1d88e8496e0",
      "score": 100,
      "release": "Abbey Road"
    }
  ]
}
```

Release group search uses the same album filter; simple results use `releasegroup` for the title:

```json
{
  "type": "releasegroup_search",
  "output": "simple",
  "query": "(Abbey Road) AND primarytype:album",
  "pageno": 1,
  "limit": 25,
  "min_score": 50,
  "primary_type": "album",
  "count": 42,
  "current_count": 1,
  "has_data": true,
  "created": "2026-07-01T12:00:00Z",
  "results": [
    {
      "mbid": "abbc4905-c25f-4c67-8e2d-19329ec48b1f",
      "score": 100,
      "releasegroup": "Abbey Road",
      "artist": "The Beatles",
      "type": "Album"
    }
  ]
}
```

### Search response (full)

Includes raw `results` and a top-level `scores` map (MBID → score).

### Lookup response (simple)

```json
{
  "type": "artist_lookup",
  "output": "simple",
  "id": "b10bbbfc-cf9e-42e6-888b-88b6b374d5d4",
  "result": {
    "mbid": "b10bbbfc-cf9e-42e6-888b-88b6b374d5d4",
    "artist": "The Beatles"
  }
}
```

### Error response (stderr)

```json
{
  "error": "limit must be between 1 and 100",
  "code": "INVALID_ARGUMENT"
}
```

| Exit code | Meaning |
|-----------|---------|
| `0` | Success |
| `1` | API or runtime error |
| `2` | Invalid arguments |

## Query syntax

Search commands use [Apache Lucene syntax](https://lucene.apache.org/core/queryparser/org/apache/lucene/queryparser/classic/package-summary.html#package_description). See the [MusicBrainz Search documentation](https://musicbrainz.org/doc/MusicBrainz_API/Search) for field reference.

| Use case | Example |
|----------|---------|
| Artist name | `artist:"The Beatles"` |
| Keyword | `Beatles` |
| Release + artist | `release:"Abbey Road" AND artist:"Beatles"` |
| Barcode | `barcode:602527306377` |
| Releases by artist MBID | `--artist-mbid b10bbbfc-...` or Lucene `arid:b10bbbfc-...` |

## Important notes

1. **Rate limiting** — MusicBrainz allows at most **one request per second** per client. Do not run concurrent requests.
2. **User-Agent** — A meaningful User-Agent (app name, version, contact URL) is required by MusicBrainz.
3. **Pagination** — Default `--limit` is `25`, `--pageno` is `1`. Valid `limit` range: 1–100. Use `has_data` (`(pageno - 1) * limit < count`) to detect when pagination is beyond the last page.
4. **Count fields** — `count` is the API total hit count (across all pages; includes release `primarytype:album` filter, excludes CLI score filter). `current_count` is the number of items in `results` on this page after score filtering. `has_data` is `true` when `(pageno - 1) * limit < count`.
5. **Score filter** — Search drops results with score &lt; 50. This affects `current_count` and `results` only, not `count` or `has_data`.
6. **Release search album filter** — `search release` always adds `primarytype:album` to the Lucene query sent to the API. Only album releases are returned; `primary_type` in the JSON envelope documents this filter.
7. **Release group search** — `search releasegroup` uses the same album and score filters; Lucene text queries default to the `releasegroup` field. Simple output uses `releasegroup` (not `release`) for the title.
8. **Lookup** — Pagination flags are ignored for lookup commands.

## Development

Run tests:

```bash
go test ./...
```

Build:

```bash
go build -o mbz .
```

## Dependencies

- [go.uploadedlobster.com/musicbrainzws2](https://pkg.go.dev/go.uploadedlobster.com/musicbrainzws2) — MusicBrainz WS2 client
- [github.com/spf13/cobra](https://github.com/spf13/cobra) — CLI framework

## License

MIT
