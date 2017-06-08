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

type csvReader struct {
	csvContent  *csv
	rowNow      int
	buffer      []byte
	bufferPoint int
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (r *csvReader) Read(b []byte) (copyed int, err error) {
	max := len(b)

	bufferLength := len(r.buffer)
	bufferLeft := bufferLength - r.bufferPoint
	if r.rowNow > r.csvContent.rowCount && bufferLeft <= 0 {
		return 0, io.EOF
	}

	for {
		if copyed >= max {
			break
		}

		bufferLength = len(r.buffer)
		bufferLeft = bufferLength - r.bufferPoint
		if bufferLeft > 0 {
			willCopy := min(max-copyed, bufferLeft)
			n := copy(b[copyed:], r.buffer[r.bufferPoint:r.bufferPoint+willCopy])
			copyed += n
			r.bufferPoint += n
			continue
		}

		if r.rowNow > r.csvContent.rowCount {
			break
		}

		r.buffer, err = r.csvContent.line(r.rowNow)
		if err != nil {
			return 0, err
		}
		r.bufferPoint = 0
		r.rowNow++
	}

	return
}

type csv struct {
	file
	haveTitleLine  bool
	delimiter      string
	quoteChar      string
	escapeChar     string
	lineTerminator string
	buffer         *bytes.Buffer
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
	b := make([]byte, 16*1024)
	buf := bytes.NewBuffer(b)

	return &csv{
		file:           c.file,
		haveTitleLine:  c.haveTitleLine,
		delimiter:      c.delimiter,
		quoteChar:      c.quoteChar,
		escapeChar:     c.escapeChar,
		lineTerminator: c.lineTerminator,
		buffer:         buf,
	}
}

func (c *csv) line(row int) ([]byte, error) {
	c.buffer.Reset()

	var err error
	var columns []string
	if row == 0 {
		columns = c.row.Title()
	} else {
		columns, err = c.row.Data()
		if err != nil {
			return nil, err
		}
	}

	lineStart := true
	for _, column := range columns {
		if lineStart {
			lineStart = false
		} else {
			c.buffer.WriteString(c.delimiter)
		}

		if c.escapeChar != "" {
			column = strings.Replace(column, c.quoteChar, c.escapeChar, -1)
		}

		c.buffer.WriteString(fmt.Sprintf("%v%v%v", c.quoteChar, column, c.quoteChar))
	}

	c.buffer.WriteString(c.lineTerminator)

	return c.buffer.Bytes(), nil
}

func (c *csv) Data() (io.Reader, error) {
	var rowNow int

	if c.haveTitleLine {
		rowNow = 0
	} else {
		rowNow = 1
	}

	return &csvReader{
		csvContent: c,
		rowNow:     rowNow,
	}, nil
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
