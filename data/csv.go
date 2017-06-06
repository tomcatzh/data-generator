package data

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const defaultDelimiter = ","
const defaultQuoteChar = "\""
const defaultEscapeChar = ""
const defaultLineTeriminator = "\r\n"

type csv struct {
	file
	haveTitleLine  bool
	delimiter      string
	quoteChar      string
	escapeChar     string
	lineTerminator string
}

func newCsv(rowCount int, namePart []namePart, delimiter string, quoteChar string, escapeChar string, lineTerminator string, haveTitleLine bool) *csv {
	return &csv{
		file: file{
			row:      row{},
			rowCount: rowCount,
			namePart: namePart,
		},
		haveTitleLine:  haveTitleLine,
		delimiter:      delimiter,
		quoteChar:      quoteChar,
		escapeChar:     escapeChar + quoteChar,
		lineTerminator: lineTerminator,
	}
}

func (c *csv) Clone() FileData {
	return newCsv(c.rowCount, c.namePart, c.delimiter, c.quoteChar, c.escapeChar, c.lineTerminator, c.haveTitleLine)
}

func (c *csv) line(columns []string) []byte {
	var buffer bytes.Buffer
	lineStart := true
	for _, column := range columns {
		if lineStart {
			lineStart = false
		} else {
			buffer.WriteString(c.delimiter)
		}

		if c.escapeChar != "" {
			column = strings.Replace(column, c.quoteChar, c.escapeChar, -1)
		}

		buffer.WriteString(fmt.Sprintf("%v%v%v", c.quoteChar, column, c.quoteChar))
	}

	return buffer.Bytes()
}

func (c *csv) Data() (io.ReadSeeker, error) {
	var buffer bytes.Buffer

	if c.haveTitleLine {
		titleData := c.row.Title()
		buffer.Write(c.line(titleData))
		buffer.WriteString(c.lineTerminator)
	}

	for i := 0; i < c.rowCount; i++ {
		rowData, err := c.row.Data()
		if err != nil {
			return nil, err
		}
		buffer.Write(c.line(rowData))

		buffer.WriteString(c.lineTerminator)
	}

	return bytes.NewReader(buffer.Bytes()), nil
}

func (c *csv) Save() (int64, error) {
	name, err := c.Name()
	if err != nil {
		return 0, err
	}
	data, err := c.Data()
	if err != nil {
		return 0, err
	}
	return c.storage.Save(name, data)
}
