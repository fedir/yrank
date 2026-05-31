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
```

## Configuration

Copy `.env.example` to `.env` and set your YouTube Data API v3 key:

```
YOUTUBE_API_KEY=your_key_here
```

The key must have **YouTube Data API v3** enabled in Google Developers Console. The app also reads `YOUTUBE_API_KEY` from the environment directly (`.env` is loaded automatically via `godotenv`).

## Architecture

Two-layer design:

**Root package (`main`)** — CLI entrypoint + rendering:
- `main.go`: reads config + CLI flags, calls `youtube` package, sorts, limits, prints
- `config.go`: loads `.env` via `godotenv`, reads `YOUTUBE_API_KEY`; `cliParameters()` parses `-p`, `-c`, `-s`, `-o`, `-m`, `-d` flags
- `view.go`: `print()` renders results as table or markdown using `tablewriter`
- `structs.go`: `Configuration` struct

**`youtube` package** — all YouTube API logic:
- `channel.go` / `playlist.go`: entry points `ChannelStatistics()` / `PlaylistStatistics()` — paginate the API, collect video IDs, then fan out goroutines via `sync.WaitGroup` + buffered channel to fetch per-video stats concurrently
- `video_statistics.go`: fetches and computes derived metrics for a single video
- `sorting.go`: `SortBy()` dispatches sort on any `[]VideoStatistics` field
- `http_request.go`: shared HTTP helper — returns body bytes directly
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

`-s` flag values: `likes`, `total-interest` (default), `positive-interest`, `total-reaction`, `global-buzz-index`, `positive-negative-coefficient` (alias: `pnc`)

`-s` and `-strategy` are mutually exclusive.

## Evaluation strategies

`-strategy` scores and sorts videos by a weighted formula. A `Score` column is prepended to the output.

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
