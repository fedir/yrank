package main

import (
	"os"
	"strconv"

	"github.com/fedir/yrank/youtube"
	"github.com/olekukonko/tablewriter"
)

func print(vs []youtube.VideoStatistics) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true) // Enable row line
	table.SetCenterSeparator("+")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")
	table.SetAlignment(tablewriter.ALIGN_CENTER)

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
