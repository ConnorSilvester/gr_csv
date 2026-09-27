package gr_csv

import (
	"os"
	"strings"
)

type CSVFile struct {
	Titles []string
	Rows   []CSVRow
}

type CSVRow struct {
	Fields []string
}

const (
	bom = "\uFEFF"
)

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

func (csv *CSVFile) FindTitleIndex(title string) int {
	for i, s := range csv.Titles {
		if s == title {
			return i
		}
	}
	return -1
}

func (csv *CSVFile) FindTitleIndexs(title string) []int {
	var result []int
	for i, s := range csv.Titles {
		if s == title {
			result = append(result, i)
		}
	}
	return result
}

func (csv *CSVFile) RowCount() int {
	return len(csv.Rows)
}
