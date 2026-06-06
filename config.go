package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fedir/yrank/youtube"
	"github.com/joho/godotenv"
)

func configuration() Configuration {
	godotenv.Load()
	apikey := os.Getenv("YOUTUBE_API_KEY")
	if apikey == "" {
		log.Fatalln("YOUTUBE_API_KEY environment variable is not set")
	}
	return Configuration{apikey: apikey}
}

// envWeights reads WEIGHT_<STRATEGY>_<KEY> variables from the environment
// and returns them as a flat map keyed by "<strategy>_<key>" (lowercase).
func envWeights() youtube.Weights {
	weights := youtube.Weights{}
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := parts[0], parts[1]
		if !strings.HasPrefix(strings.ToUpper(key), "WEIGHT_") {
			continue
		}
		suffix := strings.ToLower(strings.TrimPrefix(strings.ToUpper(key), "WEIGHT_"))
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			log.Fatalf("invalid weight value for %s: %v", key, err)
		}
		weights[suffix] = f
	}
	return weights
}

// parseWeightsFlag parses "key=val,key=val" into a Weights map.
func parseWeightsFlag(raw string) youtube.Weights {
	w := youtube.Weights{}
	if raw == "" {
		return w
	}
	for _, pair := range strings.Split(raw, ",") {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			log.Fatalf("invalid -weights pair %q, expected key=value", pair)
		}
		f, err := strconv.ParseFloat(strings.TrimSpace(kv[1]), 64)
		if err != nil {
			log.Fatalf("invalid weight value %q: %v", kv[1], err)
		}
		w[strings.TrimSpace(kv[0])] = f
	}
	return w
}

func cliParameters() (cid, pid, topSearch, output, sorting, strategy, from, weightsRaw, outFile string, maxResults int, debug, localTest bool) {
	var (
		playlistID    = flag.String("p", "", "Youtube playlist ID")
		channelID     = flag.String("c", "", "Youtube channel ID")
		searchFlag    = flag.String("top-search", "", "Search YouTube for a word/phrase and rank the matching videos")
		out           = flag.String("o", "table", "Output format {table|markdown|csv}")
		sort          = flag.String("s", "", "Sorting {total-interest|positive-interest|global-buzz-index|total-reaction|positive-negative-coefficient|pnc|likes}")
		strat         = flag.String("strategy", "", fmt.Sprintf("Evaluation strategy {%s}", knownStrategies()))
		maxRes        = flag.Int("m", 0, "Max items to return (0 = all)")
		fromDate      = flag.String("from", "", "Only include videos published on or after this date (YYYY-MM-DD)")
		weightsFlag   = flag.String("weights", "", "Override strategy weights: key=val,key=val")
		outFlag       = flag.String("out", "", "Write output to file atomically (safer than shell redirection)")
		localTestFlag = flag.Bool("local-test", false, "Use local testdata/ fixtures instead of live API calls")
		dbg           = flag.Bool("d", false, "Debug mode")
	)
	flag.Parse()

	sources := 0
	for _, s := range []string{*playlistID, *channelID, *searchFlag} {
		if s != "" {
			sources++
		}
	}
	if sources == 0 {
		log.Fatalln("One of -p (playlist), -c (channel) or -top-search must be defined")
	}
	if sources > 1 {
		log.Fatalln("-p, -c and -top-search are mutually exclusive")
	}
	if *out != "table" && *out != "markdown" && *out != "csv" {
		log.Fatalln("Unknown output format")
	}
	if *sort != "" && *strat != "" {
		log.Fatalln("-s and -strategy cannot be used together")
	}
	if *sort != "" {
		validSorts := map[string]bool{
			"likes": true, "total-interest": true, "positive-interest": true,
			"global-buzz-index": true, "total-reaction": true,
			"positive-negative-coefficient": true, "pnc": true,
		}
		if !validSorts[*sort] {
			log.Fatalln("Unknown sorting column")
		}
	}
	if *strat != "" && *strat != "all" {
		if _, ok := youtube.Strategies[*strat]; !ok {
			log.Fatalf("Unknown strategy %q, available: %s or all", *strat, knownStrategies())
		}
	}
	if *fromDate != "" {
		if _, err := time.Parse("2006-01-02", *fromDate); err != nil {
			log.Fatalln("Invalid -from date, expected format YYYY-MM-DD")
		}
	}

	// Default sort when neither -s nor -strategy is given
	sortVal := *sort
	if sortVal == "" && *strat == "" {
		sortVal = "total-interest"
	}

	return *channelID, *playlistID, *searchFlag, *out, sortVal, *strat, *fromDate, *weightsFlag, *outFlag, *maxRes, *dbg, *localTestFlag
}

func knownStrategies() string {
	names := make([]string, 0, len(youtube.Strategies))
	for k := range youtube.Strategies {
		names = append(names, k)
	}
	return strings.Join(names, "|")
}
