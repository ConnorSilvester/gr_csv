// Package gr_csv provides a simplified API for reading csv files.
//
// It offers a higher-level interface for common read operations,
// such as parsing and finding titles
package gr_csv

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
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
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseReader(file)
}

// ParseReader, parses a io.Reader into a CSVFile.
// Returns nil and error if the file doesn't exist.
func ParseReader(r io.Reader) (*CSVFile, error) {
	r = stripBOM(r)

	reader := csv.NewReader(r)
	reader.FieldsPerRecord = 0
	reader.TrimLeadingSpace = false

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	csvFile := &CSVFile{}

	if len(records) == 0 {
		return csvFile, nil
	}

	csvFile.Titles = records[0]
	csvFile.Rows = make([]CSVRow, 0, len(records)-1)
	for _, rec := range records[1:] {
		csvFile.Rows = append(csvFile.Rows, CSVRow{Fields: rec})
	}

	return csvFile, nil
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

// FindTitleIndexes, finds all the indexs of the given title inside the CSVFile headers.
// Returns a list of indexs or empty list if not found.
func (csv *CSVFile) FindTitleIndexes(title string) []int {
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

// stripBOM returns a reader that skips a leading UTF-8 BOM if one is present.
// It peeks the first three bytes and only consumes them if they match.
// non-BOM input is untouched.
func stripBOM(r io.Reader) io.Reader {
	br := bufio.NewReader(r)
	peek, err := br.Peek(3)
	if err == nil && string(peek) == bom {
		_, _ = br.Discard(3)
	}
	return br
}

