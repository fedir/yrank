# yrank = Youtube Rank analyzer

[![CI](https://github.com/fedir/yrank/actions/workflows/ci.yml/badge.svg)](https://github.com/fedir/yrank/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/fedir/yrank)](https://goreportcard.com/report/github.com/fedir/yrank)
[![GoDoc](https://godoc.org/github.com/fedir/yrank?status.svg)](https://godoc.org/github.com/fedir/yrank)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

Ranks videos in a YouTube playlist or channel by engagement metrics, so you can prioritise what to watch — especially useful for large conference playlists.

## Installation

Pre-built binaries for Linux, macOS and Windows (amd64/arm64) are attached to every
[GitHub release](https://github.com/fedir/yrank/releases).

### macOS (Homebrew)

```bash
brew install fedir/tap/yrank
```

Recent Homebrew versions require trusting a third-party tap before its formula
will load. If you see `Refusing to load formula ... from untrusted tap`, run:

```bash
brew trust fedir/tap
brew install fedir/tap/yrank
```

Upgrade to the latest release later with:

```bash
brew upgrade yrank
```

Verify the install:

```bash
yrank -version
```

### Install with Go (1.26+)

```bash
go install github.com/fedir/yrank@latest
```

### Build from source

```bash
git clone https://github.com/fedir/yrank
cd yrank
make build
```

## Configuration

You need a YouTube Data API v3 key from [Google Developers Console](https://console.developers.google.com/):

1. Create a project and enable **YouTube Data API v3**
2. Generate an API key
3. Copy `.env.example` to `.env` and set your key:

```
YOUTUBE_API_KEY=your_key_here
```

The key is also read directly from the environment if set there.

## CLI options

| Flag | Default | Description |
|---|---|---|
| `-p` | — | YouTube playlist ID |
| `-c` | — | YouTube channel ID or handle (e.g. `@Squeezie`) |
| `-top-search` | — | Search YouTube for a word/phrase and rank the matching videos |
| `-o` | `table` | Output format: `table`, `markdown`, or `csv` |
| `-out` | — | Write output to file atomically (safer than shell redirection for large exports) |
| `-in` | — | Filter an existing CSV export locally (no API): read this file, apply the view/duration filters, write the same format to `-out` |
| `-check` | — | Validate a yrank CSV export locally (no API) and exit non-zero on failure |
| `-s` | `total-interest` | Sort by metric (see below). Mutually exclusive with `-strategy` |
| `-strategy` | — | Score and rank by evaluation strategy (see below). Mutually exclusive with `-s` |
| `-weights` | — | Override strategy weights: `key=val,key=val` |
| `-from` | — | Only include videos published on or after this date (`YYYY-MM-DD`) |
| `-min-length` | `0` (no min) | Only include videos at least N **seconds** long |
| `-max-length` | `0` (no max) | Only include videos at most N **seconds** long |
| `-min-views` | `0` (no min) | Only include videos with at least N **views** |
| `-m` | `0` (all) | Maximum number of results to return |
| `-local-test` | `false` | Use local `testdata/` fixtures instead of live API calls (no quota consumed) |
| `-d` | `false` | Debug mode — prints API URLs and IDs |
| `-version` / `-V` | — | Print version and exit |

Exactly one input source — `-p`, `-c`, or `-top-search` — must be given; they are mutually exclusive.

`-top-search` uses the YouTube `search.list` endpoint, which costs **100 quota units per page** of 50 results (vs 1 unit for playlist/channel listing). It paginates up to `-m` candidates before ranking them; with the default `-m 0` it fetches a single page (≤50 videos).

See [Quota & limits](#quota--limits) for the full cost model and the hard caps the YouTube API imposes.

### Sorting (`-s`)

`total-interest` (default) · `positive-interest` · `likes` · `total-reaction` · `global-buzz-index` · `positive-negative-coefficient` · `pnc` · `duration`

### Duration filtering (`-min-length` / `-max-length`)

Every result includes a `Duration` column in **seconds** (from the API's `contentDetails.duration`). `-min-length` / `-max-length` keep only videos within the given second bounds (`0` = no limit on that side); both apply before sorting. Videos with an unknown duration (`0`, e.g. live placeholders) are dropped only when `-min-length` is set.

```bash
./yrank -c @Vsauce -min-length 300                 # only videos longer than 5 min
./yrank -c @Vsauce -max-length 60 -s duration      # Shorts-length clips, longest first
```

### View filtering (`-min-views`)

`-min-views N` keeps only videos with at least N views (`0` = no limit), applied before sorting alongside the other filters.

```bash
./yrank -c @Vsauce -min-views 1000000              # only videos past 1M views
./yrank -c @Vsauce -min-length 300 -min-views 50000
```

### Re-filter an existing export locally (`-in`, no API quota)

Already exported a channel and want a tighter cut without spending quota again? Feed an existing
CSV back in with `-in` (or `make local-filter`). It locates the `Views`/`Duration` columns by
header and writes the same format — works on both base and `-strategy all` exports.

```bash
./yrank -in sample_output/vsauce_channel_all.csv -out vsauce_long.csv -min-length 900 -min-views 100000

# or via make (IN and OUT required):
make local-filter IN=sample_output/vsauce_channel_all.csv OUT=vsauce_long.csv MIN_VIEWS=100000 MIN_LENGTH=900
```

### Publishing a channel export (`make publish-channel`)

One command to export a full channel, validate it with Go checks (`-check`), and
`git add`/`commit`/`push` just that CSV with a predefined message:

```bash
make publish-channel CHANNEL=@NASA
# → sample_output/nasa_channel_all.csv, validated, committed as
#   "chore: add @NASA full channel export", then pushed.

make publish-channel CHANNEL=@NASA EXPORT_MSG="chore: refresh NASA export"   # custom message
```

The commit is gated on the checks: if `-check` finds an empty/short export, a missing
column, a `views <= 0` row, or uniform per-video stats, the run stops before committing.

### Evaluation strategies (`-strategy`)

Each strategy scores videos by a weighted formula over raw signals, then sorts by `Score`. A `Score` column is prepended to the output.

Use `-strategy all` to compute **all 6 strategies at once** — adds one score column per strategy (`Score:viral`, `Score:educational`, …) and sorts by `total-interest`. Ideal for CSV export and cross-strategy comparison.

| Slug | Lens | Weight keys |
|---|---|---|
**`viral`** — algo/trending lens

Rewards videos with high engagement rate relative to their audience size.

```
score = 0.5 × (likes+dislikes)/views
      + 0.3 × views/max_views        ← reach normalised across dataset
      + 0.2 × comments/views
```

---

**`educational`** — tutorial/reference lens

Rewards likes and discussion; penalises recency so old-but-solid content ranks high.

```
score = 0.6 × likes/views
      + 0.3 × comments/views
      + 0.1 × 1/age_days
```

---

**`controversial`** — debate/polarising lens

Rewards a high dislike ratio multiplied by total reaction volume (log-scaled to reduce outliers).

```
score = ratio × volume
  ratio  = (dislikes+1) / (likes+1)
  volume = log(likes+dislikes+1)
```

---

**`community`** — fan/community-building lens

Comments are the primary signal; positive sentiment is secondary.

```
score = 0.5 × comments/views
      + 0.5 × norm_sentiment
  norm_sentiment = pnc / (1 + pnc),  pnc = likes/(1+dislikes)
```

---

**`evergreen`** — long-tail/SEO lens

Rewards videos that accumulate steady engagement per day since publication.

```
score = 0.5 × (likes+comments)/age_days
      + 0.5 × 1/age_days
```

---

**`hype`** — launch/premiere lens

Pure view velocity: views per day since publication.

```
score = views / age_days
```

---

`age_days` = days since `PublishedAt` (min 1). `max_views` = highest view count in the dataset (for `viral` normalisation).

| Slug | Weight keys |
|---|---|
| `viral` | `engagement`, `reach`, `comments` |
| `educational` | `likes`, `comments`, `recency` |
| `controversial` | `ratio`, `volume` |
| `community` | `comments`, `sentiment` |
| `evergreen` | `engagement`, `age` |
| `hype` | `velocity` |

**Weight override priority** (highest wins):
1. Strategy defaults (hardcoded in `youtube/strategy.go`)
2. `.env` variables: `WEIGHT_<STRATEGY>_<KEY>=0.7` — e.g. `WEIGHT_VIRAL_ENGAGEMENT=0.7`
3. `-weights` CLI flag: `key=val,key=val` — e.g. `-weights engagement=0.9,reach=0.05,comments=0.05`

## Usage examples

```bash
# Rank a playlist by total interest (default)
./yrank -p PLAYLIST_ID

# Rank in markdown, sorted by positive interest, top 10
./yrank -p PLAYLIST_ID -o markdown -s positive-interest -m 10

# Rank a whole channel using a handle
./yrank -c @Squeezie -s positive-interest -o markdown

# Search YouTube for a phrase and rank the top 20 matches by viral score
./yrank -top-search "kubernetes operator" -strategy viral -m 20

# Only videos from 2025 onwards
./yrank -c @Squeezie -from 2025-01-01 -s positive-interest

# Rank by viral strategy
./yrank -p PLAYLIST_ID -strategy viral

# Rank by viral strategy with custom weights
./yrank -p PLAYLIST_ID -strategy viral -weights engagement=0.9,reach=0.05,comments=0.05

# Export to CSV (pipe-safe, emojis and special chars handled correctly)
./yrank -c @TiboInShape -o csv -out tiboinshape.csv

# Score with all strategies at once (6 score columns)
./yrank -p PLAYLIST_ID -strategy all -o csv -out all_strategies.csv

# Use local fixtures — no API quota consumed
./yrank -p PLiVdPopzGBsV7TgjAw9GH43Ck9QCxrw5w -local-test
./yrank -p PLiVdPopzGBsV7TgjAw9GH43Ck9QCxrw5w -local-test -strategy all -o csv

# Save results to a file atomically (preferred over shell redirection)
./yrank -p PLAYLIST_ID -strategy evergreen -o markdown -out results.md
```

## Quota & limits

yrank talks to the **YouTube Data API v3**, which is metered in **quota units** (the API's
"tokens"), not in number of requests. The default project budget is **10,000 units/day**,
resetting at **midnight US-Pacific time**. This chapter is the cost model and the strategy yrank
follows to stay inside it.

### Per-endpoint cost

| Endpoint | Used for | Cost |
|---|---|---|
| `playlistItems.list` | listing a playlist / channel uploads (≤50 items/page) | **1 unit/page** |
| `videos.list` | per-video statistics + duration (**up to 50 IDs/call**) | **1 unit/call** |
| `channels.list` | resolving `@handle`, listing manual playlists | **1 unit** |
| `search.list` | `-top-search` (≤50 results/page) | **100 units/page** |

### Strategy we follow to minimise consumption

1. **Batch `videos.list` to 50 IDs per call.** Statistics are the dominant cost, and a call
   costs 1 unit whether it carries 1 ID or 50. We collect IDs from the listing, then fetch
   stats in chunks of 50 — roughly a **50× reduction** versus one call per video.
2. **Listing carries only what's free.** Pagination collects `id` + `title` + `publishedAt`
   only; everything else (views, likes, comments, duration) comes from the batched
   `videos.list`, joined back by ID.
3. **Deduplicate IDs before fetching stats.** For a channel, IDs from the uploads playlist and
   manual playlists are deduplicated *before* any `videos.list` call, so a video in several
   playlists is paid for once.
4. **Filters are client-side and free.** `-min-views` / `-min-length` / `-from` cannot be
   pushed to the API (it has no such parameters), but they run on already-fetched data, so
   they add **zero** quota.

**Rule of thumb:** a channel/playlist of *N* videos costs about **`2 × ceil(N/50)`** units
(listing + stats). Examples: Naval Group (239 videos) ≈ **10 units**; a 20k-video channel
≈ **~1,200 units**. A `-top-search` run adds **100 units per 50 candidates**.

### Hard limits imposed by the API (not by yrank)

- **~20,000-video cap per channel.** `playlistItems.list` stops paginating a channel's uploads
  playlist at ~20k items. Channels larger than that **cannot be fully exported** with `-c`
  (e.g. CCTV Video News Agency has 43,634 videos but only ~20,048 are retrievable this way).
- **`search.list` caps** at ~500 results per query (pagination limit), at 100 units/page.
- **No server-side filtering** by view count or precise duration — the data must be fetched
  before any `-min-views`/`-min-length` filter can apply.
- **`dislikeCount` is gone.** YouTube removed public dislikes from the API in December 2021, so
  that field is always `0` and the dislike-based metrics effectively run with `dislikes = 0`.
- **No resume on interruption (yet).** If the daily quota is exhausted or a rate limit hits
  mid-run, the run aborts and partial work is lost. A resumable/queued mode is planned.

### Checking your remaining quota

The exact remaining number is only visible in the Google Cloud Console
(**APIs & Services → YouTube Data API v3 → Quotas**). With just an API key you can *probe*
whether quota is still available with a 1-unit call:

```bash
curl -s -o /dev/null -w '%{http_code}\n' \
  "https://www.googleapis.com/youtube/v3/videos?part=id&id=NCU_Sebq6Tw&key=$YOUTUBE_API_KEY"
# 200 = quota available · 403 = quota exhausted / rate-limited
```

## Sample outputs

Conferences & playlists:

* [KubeCon + CloudNativeCon Europe 2026](sample_output/kubecon_cloudnativecon_europe_2026_positive_interest.md) · [all strategies (CSV)](sample_output/kubecon_2026_all_strategies.csv)
* [CNCF Observability Day Europe 2026](sample_output/cncf_observability_day_europe_2026_positive_interest.md)
* [Cloud Native Days France 2026 (CSV)](sample_output/cloud_native_days_france_2026.csv)
* [GrafanaCon — all strategies (CSV)](sample_output/grafanacon_all_strategies.csv)
* [FOSDEM 2020](sample_output/fosdem2020_positive_interest.md)
* [Squeezie — Concepts originaux](sample_output/squeezie_concepts_originaux_positive_interest.md)

Full-channel exports (all strategies, CSV):

* [Squeezie](sample_output/squeezie_channel_positive_interest.md) · [Vsauce](sample_output/vsauce_channel_all.csv) · [HugoDécrypte Actus](sample_output/hugodecrypteactus_channel_all.csv) · [Mister V](sample_output/mister_v_channel_all.csv) · [Tibo InShape](sample_output/tiboinshape_full.csv)
* [Airbus Defence and Space](sample_output/airbusds_channel_all.csv) · [Naval Group](sample_output/navalgroup_channel_all.csv) · [INA Histoire](sample_output/inahistoire_channel_all.csv)

Filtered exports (duration filters):

* [Vsauce — videos > 300s](sample_output/vsauce_longer_than_300s.csv)
* [Squeezie — videos > 900s](sample_output/squeezie_longer_than_900s.csv)
* [Anyme0233 — videos > 300s](sample_output/anyme0233_longer_than_300s.csv)

## Releases

Releases are fully automated with [GoReleaser](https://goreleaser.com/). Pushing a
`v*` tag triggers the [release workflow](.github/workflows/release.yml), which
cross-compiles binaries (linux/darwin/windows × amd64/arm64), publishes a GitHub
release with checksums and a changelog, and updates the Homebrew tap.

```bash
git tag v1.2.3
git push origin v1.2.3
```

Build a local snapshot to test the release config without publishing:

```bash
make snapshot     # requires goreleaser installed locally
```

The `Homebrew` formula update requires a `HOMEBREW_TAP_TOKEN` repository secret with
write access to `fedir/homebrew-tap`.
