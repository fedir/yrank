package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// filterCSVFile reads an existing yrank CSV export from inPath, keeps only rows
// matching the view/duration thresholds, and writes the result — same header and
// column order — to outPath (atomically) or to stdout when outPath is empty.
//
// It is format-agnostic: it locates the "Views" and "Duration" columns by header
// name, so both the base and -strategy all exports work. No API calls are made.
func filterCSVFile(inPath, outPath string, minViews, minLength, maxLength int) error {
	f, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return fmt.Errorf("parse %s: %w", inPath, err)
	}
	if len(rows) == 0 {
		return fmt.Errorf("%s is empty", inPath)
	}

	header := rows[0]
	viewsIdx := columnIndex(header, "Views")
	durIdx := columnIndex(header, "Duration")
	if viewsIdx < 0 || durIdx < 0 {
		return fmt.Errorf("%s: missing required columns (need %q and %q)", inPath, "Views", "Duration")
	}

	out := make([][]string, 0, len(rows))
	out = append(out, header)
	for _, row := range rows[1:] {
		views := atoiCol(row, viewsIdx)
		dur := atoiCol(row, durIdx)
		if minViews > 0 && views < minViews {
			continue
		}
		if minLength > 0 && dur < minLength {
			continue
		}
		if maxLength > 0 && dur > maxLength {
			continue
		}
		out = append(out, row)
	}

	return writeCSV(outPath, out)
}

// columnIndex returns the index of the named column in the header, or -1.
func columnIndex(header []string, name string) int {
	for i, h := range header {
		if h == name {
			return i
		}
	}
	return -1
}

// atoiCol parses the integer at row[idx], returning 0 for missing/unparseable cells.
func atoiCol(row []string, idx int) int {
	if idx >= len(row) {
		return 0
	}
	n, _ := strconv.Atoi(row[idx])
	return n
}

// writeCSV writes rows to path (atomically via temp-rename) or stdout when path is empty.
func writeCSV(path string, rows [][]string) error {
	if path == "" {
		w := csv.NewWriter(os.Stdout)
		if err := w.WriteAll(rows); err != nil {
			return err
		}
		w.Flush()
		return w.Error()
	}

	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	if err := w.WriteAll(rows); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	w.Flush()
	if err := w.Error(); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, path)
}
