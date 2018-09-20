package main

import (
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

	table.SetHeader([]string{"Title", "URL", "Views", "Likes", "Dislikes", "Comments"})

	for _, vsi := range vs {
		if vsi.Title != "" {
			table.Append(
				[]string{
					vsi.Title,
					vsi.URL,
					strconv.Itoa(vsi.ViewCount),
					strconv.Itoa(vsi.LikeCount),
					strconv.Itoa(vsi.DislikeCount),
					strconv.Itoa(vsi.CommentCount),
				},
			)
		}
	}
	table.Render() // Send output
}
