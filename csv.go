// Package gr_csv provides a simplified API for reading csv files.
//
// It offers a higher-level interface for common read operations,
// such as parsing and finding titles
package gr_csv

import (
	"os"
	"strings"
)

// CSVFile represents an in-memory csv file, consisting of the headers and the row data.
type CSVFile struct {
	Titles []string
	Rows   []CSVRow
}

// CSVRow represents a single row of data in the file.
type CSVRow struct {
	Fields []string
}

const (
	bom = "\uFEFF"
)

// ParseFile, parses a file into a CSVFile type given a filepath.
// Returns nil and error if the file doesn't exist.
func ParseFile(filePath string) (*CSVFile, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	fileContents := string(raw)
	fileContents = strings.TrimPrefix(fileContents, bom)

	csv := &CSVFile{}
	row := &CSVRow{}
	var field strings.Builder

	inQuotes := false
	isFirstRow := true

	i := 0
	for i < len(fileContents) {
		c := fileContents[i]

		switch {
		case c == '"' && !inQuotes:
			inQuotes = true
		case c == '"' && inQuotes:
			if i+1 < len(fileContents) && fileContents[i+1] == '"' {
				field.WriteByte('"')
				i++
			} else {
				inQuotes = false
			}
		case c == ',' && !inQuotes:
			row.Fields = append(row.Fields, field.String())
			field.Reset()
		case c == '\n' && !inQuotes:
			row.Fields = append(row.Fields, field.String())
			field.Reset()

			if isFirstRow {
				csv.Titles = row.Fields
				isFirstRow = false
			} else {
				csv.Rows = append(csv.Rows, *row)
			}
			row = &CSVRow{}
		default:
			if c != '\r' {
				field.WriteByte(c)
			}
		}
		i++
	}

	if field.Len() > 0 || len(row.Fields) > 0 {
		if field.Len() > 0 {
			row.Fields = append(row.Fields, field.String())
		}
		csv.Rows = append(csv.Rows, *row)
	}

	return csv, err
}

// FindTitleIndex, finds the first index of the given title inside the CSVFile headers.
// Returns the index or -1 if not found.
func (csv *CSVFile) FindTitleIndex(title string) int {
	for i, s := range csv.Titles {
		if s == title {
			return i
		}
	}
	return -1
}

// FindTitleIndexs, finds all the indexs of the given title inside the CSVFile headers.
// Returns a list of indexs or empty list if not found.
func (csv *CSVFile) FindTitleIndexs(title string) []int {
	var result []int
	for i, s := range csv.Titles {
		if s == title {
			result = append(result, i)
		}
	}
	return result
}

// RowCount, returns the number of rows in the csv file, excluding the header row.
func (csv *CSVFile) RowCount() int {
	return len(csv.Rows)
}
