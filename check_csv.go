package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

// checkCSVFile validates a yrank CSV export and returns an error describing the
// first failed check. It is used by `make publish-channel` to gate a commit on a
// sane export, and is intentionally dependency-free (Go, no shell/python).
//
// Checks: the file parses, carries the expected columns, has at least one data
// row, contains no impossible-engagement rows (views <= 0, which the fetch path
// should already drop), and shows per-video variation in views (guards against a
// regression to uniform stats). A summary is printed on success.
func checkCSVFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	if len(rows) < 2 {
		return fmt.Errorf("%s: no data rows", path)
	}

	header := rows[0]
	for _, col := range []string{"Title", "URL", "Views", "Duration"} {
		if columnIndex(header, col) < 0 {
			return fmt.Errorf("%s: missing required column %q", path, col)
		}
	}
	titleIdx := columnIndex(header, "Title")
	viewsIdx := columnIndex(header, "Views")
	durIdx := columnIndex(header, "Duration")

	data := rows[1:]
	distinctViews := make(map[int]struct{})
	minV, maxV := 1<<62, 0
	minD, maxD := 1<<62, 0
	for i, row := range data {
		if titleIdx >= len(row) || row[titleIdx] == "" {
			return fmt.Errorf("%s: row %d has empty Title", path, i+2)
		}
		v := atoiCol(row, viewsIdx)
		if v <= 0 {
			return fmt.Errorf("%s: row %d has non-positive Views (%d) — anomaly filter regression", path, i+2, v)
		}
		dvd := atoiCol(row, durIdx)
		distinctViews[v] = struct{}{}
		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
		if dvd < minD {
			minD = dvd
		}
		if dvd > maxD {
			maxD = dvd
		}
	}

	if len(data) > 1 && len(distinctViews) <= 1 {
		return fmt.Errorf("%s: all %d rows share the same view count — per-video stats look broken", path, len(data))
	}

	fmt.Printf("OK %s: %d videos, views %d-%d, duration %d-%d s, %d distinct view values\n",
		path, len(data), minV, maxV, minD, maxD, len(distinctViews))
	return nil
}
