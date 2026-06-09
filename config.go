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

func cliParameters() (cid, pid, topSearch, output, sorting, strategy, from, weightsRaw, outFile string, maxResults, minLength, maxLength int, debug, localTest bool) {
	var (
		playlistID    = flag.String("p", "", "Youtube playlist ID")
		channelID     = flag.String("c", "", "Youtube channel ID")
		searchFlag    = flag.String("top-search", "", "Search YouTube for a word/phrase and rank the matching videos")
		out           = flag.String("o", "table", "Output format {table|markdown|csv}")
		sort          = flag.String("s", "", "Sorting {total-interest|positive-interest|global-buzz-index|total-reaction|positive-negative-coefficient|pnc|likes|duration}")
		strategyFlag  = flag.String("strategy", "", fmt.Sprintf("Evaluation strategy {%s}", knownStrategies()))
		maxRes        = flag.Int("m", 0, "Max items to return (0 = all)")
		minLen        = flag.Int("min-length", 0, "Only include videos at least N seconds long (0 = no min)")
		maxLen        = flag.Int("max-length", 0, "Only include videos at most N seconds long (0 = no max)")
		fromDate      = flag.String("from", "", "Only include videos published on or after this date (YYYY-MM-DD)")
		weightsFlag   = flag.String("weights", "", "Override strategy weights: key=val,key=val")
		outFlag       = flag.String("out", "", "Write output to file atomically (safer than shell redirection)")
		localTestFlag = flag.Bool("local-test", false, "Use local testdata/ fixtures instead of live API calls")
		dbg           = flag.Bool("d", false, "Debug mode")
		showVersion   = flag.Bool("version", false, "Print version and exit")
		showVersionV  = flag.Bool("V", false, "Print version and exit (alias of -version)")
	)
	flag.Parse()

	if *showVersion || *showVersionV {
		fmt.Printf("yrank %s\n", version)
		os.Exit(0)
	}

	validateSources(*playlistID, *channelID, *searchFlag)
	validateOutputFormat(*out)
	validateSortStrategy(*sort, *strategyFlag)
	validateFilters(*fromDate, *minLen, *maxLen)

	// Default sort when neither -s nor -strategy is given
	sortVal := *sort
	if sortVal == "" && *strategyFlag == "" {
		sortVal = "total-interest"
	}

	return *channelID, *playlistID, *searchFlag, *out, sortVal, *strategyFlag, *fromDate, *weightsFlag, *outFlag, *maxRes, *minLen, *maxLen, *dbg, *localTestFlag
}

// validateSources enforces that exactly one input source is given.
func validateSources(playlistID, channelID, searchFlag string) {
	sources := 0
	for _, s := range []string{playlistID, channelID, searchFlag} {
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
}

// validateOutputFormat rejects unknown -o values.
func validateOutputFormat(out string) {
	if out != "table" && out != "markdown" && out != "csv" {
		log.Fatalln("Unknown output format")
	}
}

// validateSortStrategy checks -s/-strategy are not combined and are known.
func validateSortStrategy(sort, strategy string) {
	if sort != "" && strategy != "" {
		log.Fatalln("-s and -strategy cannot be used together")
	}
	if sort != "" {
		validSorts := map[string]bool{
			"likes": true, "total-interest": true, "positive-interest": true,
			"global-buzz-index": true, "total-reaction": true,
			"positive-negative-coefficient": true, "pnc": true, "duration": true,
		}
		if !validSorts[sort] {
			log.Fatalln("Unknown sorting column")
		}
	}
	if strategy != "" && strategy != "all" {
		if _, ok := youtube.Strategies[strategy]; !ok {
			log.Fatalf("Unknown strategy %q, available: %s or all", strategy, knownStrategies())
		}
	}
}

// validateFilters checks the -from date and the -min-length/-max-length window.
func validateFilters(fromDate string, minLen, maxLen int) {
	if fromDate != "" {
		if _, err := time.Parse("2006-01-02", fromDate); err != nil {
			log.Fatalln("Invalid -from date, expected format YYYY-MM-DD")
		}
	}
	if minLen < 0 || maxLen < 0 {
		log.Fatalln("-min-length and -max-length must be non-negative (seconds)")
	}
	if minLen > 0 && maxLen > 0 && minLen > maxLen {
		log.Fatalln("-min-length cannot be greater than -max-length")
	}
}

func knownStrategies() string {
	names := make([]string, 0, len(youtube.Strategies))
	for k := range youtube.Strategies {
		names = append(names, k)
	}
	return strings.Join(names, "|")
}
