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
make build    # compile → ./yrank binary (injects version via -ldflags)
make test     # go test -race ./...
make vet      # go vet ./...
make clean    # remove binary, coverage artifacts and dist/
make snapshot # local GoReleaser snapshot (no publish; needs goreleaser)

# Filter an existing CSV export locally (no API quota); IN and OUT are required
make local-filter IN=sample_output/foo.csv OUT=foo_filtered.csv MIN_VIEWS=100000 MIN_LENGTH=900

# Run a single test
go test -race -run TestName ./youtube/...

# Run the binary
./yrank -p PLAYLIST_ID
./yrank -c CHANNEL_ID -s positive-interest -o markdown -m 10
./yrank -c CHANNEL_ID -o csv -out export.csv

# Run with local fixtures (no API quota consumed)
./yrank -p PLiVdPopzGBsV7TgjAw9GH43Ck9QCxrw5w -local-test
./yrank -p PLiVdPopzGBsV7TgjAw9GH43Ck9QCxrw5w -local-test -strategy all -o csv

# Print version (overridden at build time via -ldflags -X main.version=...)
./yrank -version
```

## CI & releases

- `.github/workflows/ci.yml` — runs `make vet`, `make build`, `make test` and an offline `-local-test` smoke test on every push/PR to `master`.
- `.github/workflows/release.yml` — on a `v*` tag push, runs GoReleaser (`.goreleaser.yaml`): cross-compiles linux/darwin/windows × amd64/arm64, publishes a GitHub release with checksums + changelog, and updates the `fedir/homebrew-tap` formula (needs the `HOMEBREW_TAP_TOKEN` secret).
- `main.go` declares `var version = "dev"`; `-version`/`-V` print it. The Makefile and GoReleaser inject the real version through `-ldflags "-X main.version=..."`.

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
- `testdata/video_stats.json` — `statistics` + `contentDetails.duration` for those 7 video IDs
- `testdata/search_results.json` — `search.list` results reusing those 7 video IDs (one title carries an HTML entity to exercise unescaping)

**Adding new fixtures**: fetch the API response, strip `thumbnails`, `description`, `channelId`, `channelTitle`, `videoOwnerChannelTitle`, `videoOwnerChannelId` from snippet fields, and save to `testdata/`. For `/videos` fixtures keep `contentDetails.duration` so duration parsing/filtering is exercised offline.

**Mock routing** (`youtube/mock_transport.go`):
- URLs containing `playlistItems` → `testdata/playlist_page<N>.json` (`N` from `pageToken`, default `1`)
- URLs containing `/search` → `testdata/search_results.json`
- URLs containing `/videos` → `testdata/video_stats.json`

**In tests**: call `youtube.SetHTTPClient(youtube.NewMockClient("testdata"))` (restored via `defer`). See `youtube/mock_transport_test.go` for examples.

## Architecture

Two-layer design:

**Root package (`main`)** — CLI entrypoint + rendering:
- `main.go`: reads config + CLI flags, calls `youtube` package, sorts, limits, prints
- `config.go`: loads `.env` via `godotenv`, reads `YOUTUBE_API_KEY`; `cliParameters()` parses `-p`, `-c`, `-top-search`, `-in`, `-s`, `-o`, `-out`, `-m`, `-from`, `-min-length`, `-max-length`, `-min-views`, `-strategy`, `-weights`, `-local-test`, `-d` flags. Exactly one of `-p`/`-c`/`-top-search` is required (unless `-in` is used); they are mutually exclusive
- `filter_csv.go`: `-in FILE` mode — `filterCSVFile()` reads an existing yrank CSV export, locates the `Views`/`Duration` columns by header name, applies `-min-views`/`-min-length`/`-max-length`, and writes the same columns/format to `-out` (atomically) or stdout. No API key required; `main()` handles `-in` before calling `configuration()`. Wrapped by `make local-filter`
- `view.go`: `print()` renders results as `table`, `markdown`, or `csv`; `printToFile()` writes atomically via temp-rename; `mdSafe()` escapes `|` in titles for markdown
- `structs.go`: `Configuration` struct

**`youtube` package** — all YouTube API logic:
- `channel.go` / `playlist.go` / `search.go`: entry points `ChannelStatistics()` / `PlaylistStatistics()` / `SearchStatistics()` — paginate the listing endpoint into `[]videoRef` (id + title + publishedAt only, no stats), then hand the refs to `collectStats()`. `SearchStatistics()` hits the `search.list` endpoint (100 quota units/page), paginates up to `maxResults` (single page when `≤0`), and `html.UnescapeString`s titles. `ChannelStatistics()` gathers refs from the uploads playlist + manual playlists and **dedups IDs before** fetching stats, so overlapping videos cost quota only once
- `video_statistics.go`: `collectStats()` fetches stats in **batched `videos.list` calls of up to 50 IDs each** (`maxIDsPerBatch`), one quota unit per call regardless of ID count — the dominant quota saver (~50× fewer stats calls; also avoids the per-100s rate limit on large channels). A bounded worker pool (`statsWorkers`) runs the batches concurrently. `buildVideoStatistics()` maps each returned item back to its ref by ID and computes derived metrics (`parseISO8601Duration()` turns the ISO-8601 `contentDetails.duration` into seconds; `isAnomalousStats()` drops impossible-engagement rows). Note: the API no longer returns `dislikeCount` (removed by YouTube Dec 2021), so it is always `0`
- `sorting.go`: `SortBy()` dispatches sort; `ApplyStrategy()` scores + sorts by one strategy; `ApplyAllStrategies()` scores with all 6 strategies (used by `-strategy all`), stores results in `VideoStatistics.AllScores`
- `mock_transport.go`: `MockTransport` / `NewMockClient()` / `SetHTTPClient()` — injectable HTTP client for tests and `-local-test` mode
- `http_request.go`: shared HTTP helper using injectable `httpClient` var
- `structs.go`: `VideoStatistics`, API response shapes

## Quota / token consumption strategy

The YouTube Data API v3 is metered in **quota units** (default **10,000/day**, reset midnight
US-Pacific), not request count. Keeping consumption low is a core design constraint — respect it
when changing any fetch path.

**Per-endpoint cost:** `playlistItems.list` = 1 unit/page (≤50 items); `videos.list` = **1 unit
per call regardless of how many IDs (≤50) or parts**; `channels.list` = 1 unit; `search.list` =
**100 units/page**.

**Rules the code follows (don't regress these):**
1. **Batch `videos.list` to ≤50 IDs/call** (`collectStats`/`fetchStatsBatch`, `maxIDsPerBatch`
   in `youtube/video_statistics.go`). Stats are the dominant cost; batching is ~50× cheaper than
   one call per video. Never reintroduce a one-ID-per-call fetch.
2. **Listing collects only `id`/`title`/`publishedAt`** (`playlistRefs`, `search.go`); all other
   fields come from the batched `videos.list`, joined by ID.
3. **Dedup IDs before fetching stats** (`channel.go` dedups uploads + manual playlists first).
4. **Filters stay client-side** (`-min-views`/`-min-length`/`-from` in `main.go`) — the API has
   no server-side view/duration filter, so they cost nothing extra but cannot reduce fetches.

**Cost rule of thumb:** a channel/playlist of *N* videos ≈ `2 × ceil(N/50)` units. `-top-search`
adds 100 units per 50 candidates.

**Hard API limits (upstream, not fixable here):**
- `playlistItems.list` caps a channel's uploads playlist at **~20,000 items** — larger channels
  (e.g. CCTV = 43,634 videos) can only be partially exported via `-c` (~20,048 retrievable).
- `search.list` caps at ~500 results/query.
- `dislikeCount` is no longer returned (removed Dec 2021) → always `0`.
- No resume yet if quota/rate-limit aborts a run (resumable mode is planned).

Probe remaining quota with a 1-unit call (`videos.list?part=id&id=…`): HTTP 200 = available,
403 = exhausted/rate-limited. Exact numbers live in the Cloud Console quotas page.

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

`-s` flag values: `likes`, `total-interest` (default), `positive-interest`, `total-reaction`, `global-buzz-index`, `positive-negative-coefficient` (alias: `pnc`), `duration`

`-s` and `-strategy` are mutually exclusive.

## Duration filtering

Every video carries a `Duration` column (seconds), derived from the API's ISO-8601 `contentDetails.duration`.

`-min-length N` / `-max-length N` keep only videos whose duration is within `[N, …]` / `[…, N]` seconds (`0` = no limit on that side). Both are applied (alongside `-from`) before sorting/strategy scoring. Videos with an unknown duration (`0`, e.g. live placeholders) are dropped only when `-min-length` is set.

```bash
./yrank -c CHANNEL_ID -min-length 60                 # exclude Shorts
./yrank -c CHANNEL_ID -max-length 600 -s duration    # videos ≤10 min, longest first
./yrank -c CHANNEL_ID -min-length 120 -max-length 1800
```

## View filtering

`-min-views N` keeps only videos with at least `N` views (`0` = no limit). Applied (alongside `-from` and the duration filters) before sorting/strategy scoring.

```bash
./yrank -c CHANNEL_ID -min-views 1000000             # only videos past 1M views
./yrank -c CHANNEL_ID -min-length 300 -min-views 50000
```

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
