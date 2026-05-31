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
