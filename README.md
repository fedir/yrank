# yrank = Youtube Rank analyzer

[![CI](https://github.com/fedir/yrank/actions/workflows/ci.yml/badge.svg)](https://github.com/fedir/yrank/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/fedir/yrank)](https://goreportcard.com/report/github.com/fedir/yrank)
[![GoDoc](https://godoc.org/github.com/fedir/yrank?status.svg)](https://godoc.org/github.com/fedir/yrank)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

Ranks videos in a YouTube playlist or channel by engagement metrics, so you can prioritise what to watch — especially useful for large conference playlists.

## Installation

Pre-built binaries for Linux, macOS and Windows (amd64/arm64) are attached to every
[GitHub release](https://github.com/fedir/yrank/releases).

Homebrew:

```bash
brew install fedir/tap/yrank
```

With Go 1.26+:

```bash
go install github.com/fedir/yrank@latest
```

Or build from source:

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
| `-s` | `total-interest` | Sort by metric (see below). Mutually exclusive with `-strategy` |
| `-strategy` | — | Score and rank by evaluation strategy (see below). Mutually exclusive with `-s` |
| `-weights` | — | Override strategy weights: `key=val,key=val` |
| `-from` | — | Only include videos published on or after this date (`YYYY-MM-DD`) |
| `-min-length` | `0` (no min) | Only include videos at least N **seconds** long |
| `-max-length` | `0` (no max) | Only include videos at most N **seconds** long |
| `-m` | `0` (all) | Maximum number of results to return |
| `-local-test` | `false` | Use local `testdata/` fixtures instead of live API calls (no quota consumed) |
| `-d` | `false` | Debug mode — prints API URLs and IDs |
| `-version` / `-V` | — | Print version and exit |

Exactly one input source — `-p`, `-c`, or `-top-search` — must be given; they are mutually exclusive.

`-top-search` uses the YouTube `search.list` endpoint, which costs **100 quota units per page** of 50 results (vs 1 unit for playlist/channel listing). It paginates up to `-m` candidates before ranking them; with the default `-m 0` it fetches a single page (≤50 videos).

### Sorting (`-s`)

`total-interest` (default) · `positive-interest` · `likes` · `total-reaction` · `global-buzz-index` · `positive-negative-coefficient` · `pnc` · `duration`

### Duration filtering (`-min-length` / `-max-length`)

Every result includes a `Duration` column in **seconds** (from the API's `contentDetails.duration`). `-min-length` / `-max-length` keep only videos within the given second bounds (`0` = no limit on that side); both apply before sorting. Videos with an unknown duration (`0`, e.g. live placeholders) are dropped only when `-min-length` is set.

```bash
./yrank -c @Vsauce -min-length 300                 # only videos longer than 5 min
./yrank -c @Vsauce -max-length 60 -s duration      # Shorts-length clips, longest first
```

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

## Sample outputs

* [KubeCon + CloudNativeCon Europe 2026](sample_output/kubecon_cloudnativecon_europe_2026_positive_interest.md)
* [CNCF Observability Day Europe 2026](sample_output/cncf_observability_day_europe_2026_positive_interest.md)
* [Squeezie — Concepts originaux](sample_output/squeezie_concepts_originaux_positive_interest.md)
* [Squeezie — Full channel](sample_output/squeezie_channel_positive_interest.md)
* [FOSDEM 2020](sample_output/fosdem2020_positive_interest.md)

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
