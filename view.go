package main

import (
	"os"

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
		v := []string{vsi.Title, vsi.URL, vsi.ViewCount, vsi.LikeCount, vsi.DislikeCount, vsi.CommentCount}
		table.Append(v)
	}
	table.Render() // Send output
}
