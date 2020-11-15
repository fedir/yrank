package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fedir/yrank/youtube"
	"github.com/olekukonko/tablewriter"
)

func print(vs []youtube.VideoStatistics, of string) {

	table := tablewriter.NewWriter(os.Stdout)

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

	table.SetHeader([]string{
		"Title",
		"URL",
		"Published at",
		"Positive interestingness",
		"Total interestingness",
		"Views",
		"Likes",
		"Dislikes",
		"Comments",
		"Total reaction",
		"Global buzz index",
	})

	for _, vsi := range vs {
		if vsi.Title != "" {
			table.Append(
				[]string{
					vsi.Title,
					vsi.URL,
					fmt.Sprintf(vsi.PublishedAt.Format("2006-01-02 15:04:05")),
					fmt.Sprintf("%.4f", vsi.PositiveInterestingness),
					fmt.Sprintf("%.4f", vsi.TotalInterestingness),
					strconv.Itoa(vsi.ViewCount),
					strconv.Itoa(vsi.LikeCount),
					strconv.Itoa(vsi.DislikeCount),
					strconv.Itoa(vsi.CommentCount),
					strconv.Itoa(vsi.TotalReaction),
					strconv.Itoa(vsi.GlobalBuzzIndex),
				},
			)
		}
	}
	table.Render() // Send output
}
