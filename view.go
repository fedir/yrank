package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/fedir/yrank/youtube"
	"github.com/olekukonko/tablewriter"
)

func print(vs []youtube.VideoStatistics, of string, showScore, allScores bool) {
	printTo(os.Stdout, vs, of, showScore, allScores)
}

// printToFile writes output atomically: data goes to a temp file first,
// then the temp file is renamed to path only on success.
func printToFile(path string, vs []youtube.VideoStatistics, of string, showScore, allScores bool) error {
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	printTo(f, vs, of, showScore, allScores)
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, path)
}

// mdSafe escapes pipe characters in a string so they don't break markdown tables.
func mdSafe(s string) string {
	return strings.ReplaceAll(s, "|", `\|`)
}

func printTo(out io.Writer, vs []youtube.VideoStatistics, of string, showScore, allScores bool) {
	baseHeaders := []string{
		"Title",
		"URL",
		"Published at",
		"Duration",
		"Positive interestingness",
		"Positive negative coefficient",
		"Total interestingness",
		"Views",
		"Likes",
		"Dislikes",
		"Comments",
		"Total reaction",
		"Global buzz index",
	}
	var headers []string
	switch {
	case allScores:
		scoreHeaders := make([]string, len(youtube.StrategyOrder))
		for i, slug := range youtube.StrategyOrder {
			scoreHeaders[i] = "Score:" + slug
		}
		// Title, URL first, then strategy scores, then remaining metrics
		headers = append([]string{"Title", "URL"}, scoreHeaders...)
		headers = append(headers, baseHeaders[2:]...)
	case showScore:
		headers = append([]string{"Score"}, baseHeaders...)
	default:
		headers = baseHeaders
	}

	if of == "csv" {
		w := csv.NewWriter(out)
		_ = w.Write(headers)
		for _, vsi := range vs {
			if vsi.Title == "" {
				continue
			}
			row := buildRow(vsi, of, showScore, allScores)
			_ = w.Write(row)
		}
		w.Flush()
		return
	}

	table := tablewriter.NewWriter(out)

	if of == "table" {
		table.SetRowLine(true)
		table.SetCenterSeparator("+")
		table.SetColumnSeparator("|")
		table.SetRowSeparator("-")
		table.SetAlignment(tablewriter.ALIGN_CENTER)
	} else if of == "markdown" {
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
		table.SetAutoWrapText(false)
	}

	table.SetHeader(headers)

	for _, vsi := range vs {
		if vsi.Title == "" {
			continue
		}
		table.Append(buildRow(vsi, of, showScore, allScores))
	}
	table.Render()
}

func buildRow(vsi youtube.VideoStatistics, of string, showScore, allScores bool) []string {
	title := vsi.Title
	if of == "markdown" {
		title = mdSafe(title)
	}
	base := []string{
		title,
		vsi.URL,
		vsi.PublishedAt.Format("2006-01-02 15:04:05"),
		strconv.Itoa(vsi.Duration),
		fmt.Sprintf("%.4f", vsi.PositiveInterestingness),
		fmt.Sprintf("%.4f", vsi.PositiveNegativeCoefficient),
		fmt.Sprintf("%.4f", vsi.TotalInterestingness),
		strconv.Itoa(vsi.ViewCount),
		strconv.Itoa(vsi.LikeCount),
		strconv.Itoa(vsi.DislikeCount),
		strconv.Itoa(vsi.CommentCount),
		strconv.Itoa(vsi.TotalReaction),
		strconv.Itoa(vsi.GlobalBuzzIndex),
	}
	switch {
	case allScores:
		scores := make([]string, len(youtube.StrategyOrder))
		for i, slug := range youtube.StrategyOrder {
			scores[i] = fmt.Sprintf("%.6f", vsi.AllScores[slug])
		}
		// Title, URL first, then scores, then remaining metrics
		row := append([]string{base[0], base[1]}, scores...)
		return append(row, base[2:]...)
	case showScore:
		return append([]string{fmt.Sprintf("%.6f", vsi.Score)}, base...)
	default:
		return base
	}
}
