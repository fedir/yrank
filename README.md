# yrank = Youtube Rank analyzer

[![Build Status](https://travis-ci.org/fedir/yrank.svg?branch=master)](https://travis-ci.org/fedir/yrank)
[![codecov](https://codecov.io/gh/fedir/yrank/branch/master/graph/badge.svg)](https://codecov.io/gh/fedir/yrank)

Ranks videos in a YouTube playlist or channel by engagement metrics, so you can prioritise what to watch — especially useful for large conference playlists.

## Installation

Go 1.26+ is required.

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
| `-o` | `table` | Output format: `table` or `markdown` |
| `-s` | `total-interest` | Sort by metric (see below). Mutually exclusive with `-strategy` |
| `-strategy` | — | Score and rank by evaluation strategy (see below). Mutually exclusive with `-s` |
| `-weights` | — | Override strategy weights: `key=val,key=val` |
| `-from` | — | Only include videos published on or after this date (`YYYY-MM-DD`) |
| `-m` | `0` (all) | Maximum number of results to return |
| `-d` | `false` | Debug mode — prints API URLs and IDs |

### Sorting (`-s`)

`total-interest` (default) · `positive-interest` · `likes` · `total-reaction` · `global-buzz-index` · `positive-negative-coefficient` · `pnc`

### Evaluation strategies (`-strategy`)

Each strategy scores videos by a weighted formula over raw signals, then sorts by `Score`. A `Score` column is prepended to the output.

| Slug | Lens | Weight keys |
|---|---|---|
| `viral` | Algo/trending — engagement rate on a large audience | `engagement`, `reach`, `comments` |
| `educational` | Tutorial/reference — likes + discussion, age-discounted | `likes`, `comments`, `recency` |
| `controversial` | Debate/polarising — dislike ratio × reaction volume | `ratio`, `volume` |
| `community` | Fan engagement — comments first, sentiment second | `comments`, `sentiment` |
| `evergreen` | Long-tail/SEO — steady engagement per day of life | `engagement`, `age` |
| `hype` | Launch velocity — views per day since publication | `velocity` |

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

# Only videos from 2025 onwards
./yrank -c @Squeezie -from 2025-01-01 -s positive-interest

# Rank by viral strategy
./yrank -p PLAYLIST_ID -strategy viral

# Rank by viral strategy with custom weights
./yrank -p PLAYLIST_ID -strategy viral -weights engagement=0.9,reach=0.05,comments=0.05

# Save results to a file
./yrank -p PLAYLIST_ID -strategy evergreen -o markdown > results.md
```

## Sample outputs

* [KubeCon + CloudNativeCon Europe 2026](sample_output/kubecon_cloudnativecon_europe_2026_positive_interest.md)
* [CNCF Observability Day Europe 2026](sample_output/cncf_observability_day_europe_2026_positive_interest.md)
* [Squeezie — Concepts originaux](sample_output/squeezie_concepts_originaux_positive_interest.md)
* [Squeezie — Full channel](sample_output/squeezie_channel_positive_interest.md)
* [FOSDEM 2020](sample_output/fosdem2020_positive_interest.md)
