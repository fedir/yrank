package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/fedir/yrank/youtube"
	"github.com/olekukonko/tablewriter"
)

func print(vs []youtube.VideoStatistics, of string, showScore bool) {
	printTo(os.Stdout, vs, of, showScore)
}

func printTo(out io.Writer, vs []youtube.VideoStatistics, of string, showScore bool) {
	headers := []string{
		"Title",
		"URL",
		"Published at",
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
	if showScore {
		headers = append([]string{"Score"}, headers...)
	}

	if of == "csv" {
		w := csv.NewWriter(out)
		_ = w.Write(headers)
		for _, vsi := range vs {
			if vsi.Title == "" {
				continue
			}
			row := []string{
				vsi.Title,
				vsi.URL,
				vsi.PublishedAt.Format("2006-01-02 15:04:05"),
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
			if showScore {
				row = append([]string{fmt.Sprintf("%.6f", vsi.Score)}, row...)
			}
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
		row := []string{
			vsi.Title,
			vsi.URL,
			vsi.PublishedAt.Format("2006-01-02 15:04:05"),
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
		if showScore {
			row = append([]string{fmt.Sprintf("%.6f", vsi.Score)}, row...)
		}
		table.Append(row)
	}
	table.Render()
}
