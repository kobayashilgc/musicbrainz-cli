# musicbrainz-cli

<p align="center">
  <img src="doc/mbz.png" alt="mbz" height="200" />
  &nbsp;&nbsp;
  <img src="doc/cobra.png" alt="cobra" height="200" />
</p>

[English](./README.md) | [简体中文](./README.zh-CN.md)

A command-line tool for querying [MusicBrainz Web Service v2](https://musicbrainz.org/doc/MusicBrainz_API). Search and look up artist and release metadata with JSON output designed for scripting and pipelines.

## Features

- **Search** artists and releases using Lucene query syntax
- **Lookup** artists and releases by MusicBrainz ID (MBID)
- **Pagination** via `--limit` and `--offset` (search only)
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
│   ├── artist <query>     Search artists
│   └── release <query>    Search releases
└── lookup
    ├── artist <mbid>       Look up an artist
    └── release <mbid>      Look up a release
```

### Global flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--limit` | `-l` | `25` | Page size (1–100); search only |
| `--offset` | `-o` | `0` | Result offset (≥ 0); search only |
| `--output` | | `simple` | Output mode: `simple` or `full` |
| `--user-agent` | | auto | HTTP User-Agent override |
| `--contact` | | repo URL | Contact URL embedded in User-Agent |
| `--api-url` | | `https://musicbrainz.org/ws/2/` | MusicBrainz WS2 base URL |

### Lookup flags

| Flag | Description |
|------|-------------|
| `--inc` | Include related data (repeatable), e.g. `releases`, `artist-credits`, `media` |

## Examples

### Search artists

```bash
mbz search artist "Beatles"
mbz search artist 'artist:"The Beatles"' --limit 10 --offset 0
mbz search artist 'artist:"The Beatles"' --limit 10 --offset 10   # next page
mbz search artist "Beatles" --output full                         # full API JSON
```

### Search releases

```bash
mbz search release 'release:"Abbey Road" AND artist:"Beatles"'
mbz search release "Abbey Road" --limit 5
```

### Lookup

```bash
mbz lookup artist b10bbbfc-cf9e-42e6-888b-88b6b374d5d4
mbz lookup release 464a321e-97a0-4654-8a7a-d1d88e8496e0 --inc artist-credits --inc media
```

## Output format

### Modes

| Mode | Flag | Description |
|------|------|-------------|
| `simple` | `--output simple` (default) | Extracts key fields; omits missing data |
| `full` | `--output full` | Full MusicBrainz API entity structures |

**Simple mode fields** (included only when present in the API response):

`mbid`, `score`, `artist`, `release`, `type`, `country`, `date`, `format`, `barcode`, `alias`, `primary_alias`, `tag`

### Search response (simple)

```json
{
  "type": "artist_search",
  "output": "simple",
  "query": "Beatles",
  "offset": 0,
  "limit": 25,
  "min_score": 50,
  "count": 1,
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

## Important notes

1. **Rate limiting** — MusicBrainz allows at most **one request per second** per client. Do not run concurrent requests.
2. **User-Agent** — A meaningful User-Agent (app name, version, contact URL) is required by MusicBrainz.
3. **Pagination** — Default `--limit` is `25`, `--offset` is `0`. Valid `limit` range: 1–100.
4. **Score filter** — Search drops results with score &lt; 50. The `count` field reflects filtered results.
5. **Lookup** — Pagination flags are ignored for lookup commands.

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
