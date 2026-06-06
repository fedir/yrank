# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Git commits

Use simple conventional commit one-liners. No `Co-Authored-By` or other trailers.

```
feat: add markdown output format
fix: handle missing API key gracefully
refactor: extract pagination logic
```

## Commands

```bash
make build    # compile → ./yrank binary
make test     # go test -race ./...
make vet      # go vet ./...
make clean    # remove binary and coverage artifacts

# Run a single test
go test -race -run TestName ./youtube/...

# Run the binary
./yrank -p PLAYLIST_ID
./yrank -c CHANNEL_ID -s positive-interest -o markdown -m 10
./yrank -c CHANNEL_ID -o csv -out export.csv

# Run with local fixtures (no API quota consumed)
./yrank -p PLiVdPopzGBsV7TgjAw9GH43Ck9QCxrw5w -local-test
./yrank -p PLiVdPopzGBsV7TgjAw9GH43Ck9QCxrw5w -local-test -strategy all -o csv
```

## Configuration

Copy `.env.example` to `.env` and set your YouTube Data API v3 key:

```
YOUTUBE_API_KEY=your_key_here
```

The key must have **YouTube Data API v3** enabled in Google Developers Console. The app also reads `YOUTUBE_API_KEY` from the environment directly (`.env` is loaded automatically via `godotenv`).

## Local test mode

The `-local-test` flag replaces the live HTTP client with a mock that serves pre-recorded JSON fixtures from `testdata/`. This avoids consuming API quota during development and enables deterministic testing.

**Fixture files** (playlist `PLiVdPopzGBsV7TgjAw9GH43Ck9QCxrw5w`, 7 videos):
- `testdata/playlist_page1.json` — playlist items page 1 (no sensitive data)
- `testdata/video_stats.json` — video statistics for those 7 video IDs
- `testdata/search_results.json` — `search.list` results reusing those 7 video IDs (one title carries an HTML entity to exercise unescaping)

**Adding new fixtures**: fetch the API response, strip `thumbnails`, `description`, `channelId`, `channelTitle`, `videoOwnerChannelTitle`, `videoOwnerChannelId` from snippet fields, and save to `testdata/`.

**Mock routing** (`youtube/mock_transport.go`):
- URLs containing `playlistItems` → `testdata/playlist_page<N>.json` (`N` from `pageToken`, default `1`)
- URLs containing `/search` → `testdata/search_results.json`
- URLs containing `/videos` → `testdata/video_stats.json`

**In tests**: call `youtube.SetHTTPClient(youtube.NewMockClient("testdata"))` (restored via `defer`). See `youtube/mock_transport_test.go` for examples.

## Architecture

Two-layer design:

**Root package (`main`)** — CLI entrypoint + rendering:
- `main.go`: reads config + CLI flags, calls `youtube` package, sorts, limits, prints
- `config.go`: loads `.env` via `godotenv`, reads `YOUTUBE_API_KEY`; `cliParameters()` parses `-p`, `-c`, `-top-search`, `-s`, `-o`, `-out`, `-m`, `-from`, `-strategy`, `-weights`, `-local-test`, `-d` flags. Exactly one of `-p`/`-c`/`-top-search` is required; they are mutually exclusive
- `view.go`: `print()` renders results as `table`, `markdown`, or `csv`; `printToFile()` writes atomically via temp-rename; `mdSafe()` escapes `|` in titles for markdown
- `structs.go`: `Configuration` struct

**`youtube` package** — all YouTube API logic:
- `channel.go` / `playlist.go` / `search.go`: entry points `ChannelStatistics()` / `PlaylistStatistics()` / `SearchStatistics()` — paginate the API, collect video IDs, then fan out goroutines via `sync.WaitGroup` + buffered channel to fetch per-video stats concurrently. `SearchStatistics()` hits the `search.list` endpoint (100 quota units/page), paginates up to `maxResults` (single page when `≤0`), and `html.UnescapeString`s titles
- `video_statistics.go`: fetches and computes derived metrics for a single video
- `sorting.go`: `SortBy()` dispatches sort; `ApplyStrategy()` scores + sorts by one strategy; `ApplyAllStrategies()` scores with all 6 strategies (used by `-strategy all`), stores results in `VideoStatistics.AllScores`
- `mock_transport.go`: `MockTransport` / `NewMockClient()` / `SetHTTPClient()` — injectable HTTP client for tests and `-local-test` mode
- `http_request.go`: shared HTTP helper using injectable `httpClient` var
- `structs.go`: `VideoStatistics`, API response shapes

## Metrics computed per video

| Field | Formula |
|---|---|
| `TotalReaction` | likes + dislikes + comments |
| `PositiveInterestingness` | (likes − dislikes) / views |
| `TotalInterestingness` | (likes + dislikes + comments) / views |
| `GlobalBuzzIndex` | views × (likes + dislikes + comments) |
| `PositiveNegativeCoefficient` | likes / (1 + dislikes) |

## Sorting options

`-o` flag values: `table` (default), `markdown`, `csv`

`-out FILE` writes output atomically to a file (temp-rename pattern). Prefer over shell redirection for large exports.

`-s` flag values: `likes`, `total-interest` (default), `positive-interest`, `total-reaction`, `global-buzz-index`, `positive-negative-coefficient` (alias: `pnc`)

`-s` and `-strategy` are mutually exclusive.

## Evaluation strategies

`-strategy` scores and sorts videos by a weighted formula. A `Score` column is prepended to the output.

`-strategy all` runs all 6 strategies simultaneously — prepends one `Score:<slug>` column per strategy and sorts by `total-interest`. Ideal for CSV export and cross-strategy comparison. `-weights` has no effect in `all` mode (default weights are used per strategy).

| Slug | Lens | Weight keys |
|---|---|---|
| `viral` | Algo/trending | `engagement`, `reach`, `comments` |
| `educational` | Tutorial/reference | `likes`, `comments`, `recency` |
| `controversial` | Debate/polarising | `ratio`, `volume` |
| `community` | Fan engagement | `comments`, `sentiment` |
| `evergreen` | Long-tail/SEO | `engagement`, `age` |
| `hype` | Launch velocity | `velocity` |

**Weight resolution order** (highest priority wins):
1. Strategy defaults (in `youtube/strategy.go`)
2. `.env` variables — `WEIGHT_<STRATEGY>_<KEY>=0.7` (e.g. `WEIGHT_VIRAL_ENGAGEMENT=0.7`)
3. `-weights` CLI flag — `key=val,key=val` (e.g. `-weights engagement=0.9,reach=0.05,comments=0.05`)

```bash
./yrank -p PLAYLIST_ID -strategy viral
./yrank -c @Squeezie -strategy evergreen -o markdown -m 10
./yrank -p PLAYLIST_ID -strategy viral -weights engagement=0.9,reach=0.05,comments=0.05
```
