package format

import (
	"compress/gzip"
	"fmt"
	"io"
	"strings"

	"github.com/tomcatzh/data-generator/data/column"

	"github.com/tomcatzh/data-generator/data/compress"
	"github.com/tomcatzh/data-generator/data/row"
	"github.com/tomcatzh/data-generator/misc"
)

const defaultDelimiter = ","
const defaultQuoteChar = "\""
const defaultEscapeChar = ""
const defaultLineTeriminator = "\r\n"

const (
	csvNoCompress = iota
	csvGzipCompress
	csvZlibCompress
)

type csv struct {
	haveTitleLine  bool
	delimiter      string
	quoteChar      string
	escapeChar     string
	lineTerminator string
	compress       int
	compressLevel  int
}

// Data returns a csv reader
func (c *csv) Data(row *row.Set) (io.Reader, error) {
	var rowNow int

	if c.haveTitleLine {
		rowNow = 0
	} else {
		rowNow = 1
	}

	r := &csvReader{
		csvContent: c,
		rowContent: row,
		rowNow:     rowNow,
	}

	switch c.compress {
	case csvGzipCompress:
		return compress.NewGzipReader(r, c.compressLevel), nil
	default:
		return r, nil
	}
}

func newCsv(s map[string]interface{}) (*csv, error) {
	var result csv

	result.haveTitleLine, _ = s["HaveTitleLine"].(bool)

	result.compress = csvNoCompress
	compress, ok := s["Compress"].(string)
	if !ok {
		compress = "none"
	}
	c := strings.Split(compress, ":")
	var compressType, compressLevel string
	if len(c) == 1 {
		compressType = compress
		compressLevel = "normal"
	} else {
		compressType = c[0]
		compressLevel = c[1]
	}

	switch compressType {
	case "gzip":
		result.compress = csvGzipCompress
		switch compressLevel {
		case "normal":
			result.compressLevel = gzip.DefaultCompression
		case "fastest":
			result.compressLevel = gzip.BestSpeed
		case "best":
			result.compressLevel = gzip.BestCompression
		default:
			return nil, fmt.Errorf("unknow compress level: %v", compressLevel)
		}
	case "none":
		result.compress = csvNoCompress
	default:
		return nil, fmt.Errorf("unknow compress type: %v", compressType)
	}

	result.delimiter, ok = s["Delimiter"].(string)
	if !ok || result.delimiter == "" {
		result.delimiter = defaultDelimiter
	}
	result.quoteChar, ok = s["Quotechar"].(string)
	if !ok {
		result.quoteChar = defaultQuoteChar
	}
	result.escapeChar, ok = s["Escapechar"].(string)
	if !ok || result.escapeChar == "" {
		result.escapeChar = defaultEscapeChar
	}
	result.lineTerminator, ok = s["Lineterminator"].(string)
	if !ok || result.lineTerminator == "" {
		result.lineTerminator = defaultLineTeriminator
	}

	return &result, nil
}

func (c *csv) line(row *row.Set, rowStart int, rowMax int) ([]byte, error) {
	row.Buffer.Reset()

	var err error
	var columns []string
	endNow := false

	for i := rowStart; i < rowMax; i++ {
		if i == 0 {
			columns = row.Title()
		} else {
			columns, err = row.Data()
			if err != nil {
				if column.IsDataOver(err) {
					endNow = true
				} else {
					return nil, err
				}
			}
		}

		lineStart := true
		for _, column := range columns {
			if lineStart {
				lineStart = false
			} else {
				row.Buffer.WriteString(c.delimiter)
			}

			if c.escapeChar != "" {
				column = strings.Replace(column, c.quoteChar, c.escapeChar, -1)
			}

			row.Buffer.WriteString(c.quoteChar)
			row.Buffer.WriteString(column)
			row.Buffer.WriteString(c.quoteChar)
		}

		if endNow {
			break
		}

		row.Buffer.WriteString(c.lineTerminator)
	}
	return row.Buffer.Bytes(), nil
}

type csvReader struct {
	csvContent  *csv
	rowContent  *row.Set
	rowNow      int
	buffer      []byte
	bufferPoint int
	step        int
}

func (r *csvReader) Read(b []byte) (copyed int, err error) {
	max := len(b)
	if r.step == 0 {
		r.step = 1
	}

	bufferLength := len(r.buffer)
	bufferLeft := bufferLength - r.bufferPoint
	if r.rowNow > r.rowContent.RowCount && bufferLeft <= 0 {
		return 0, io.EOF
	}

	readTimes := 1
	for {
		if copyed >= max {
			break
		}

		bufferLength = len(r.buffer)
		bufferLeft = bufferLength - r.bufferPoint
		if bufferLeft > 0 {
			willCopy := misc.MinInt(max-copyed, bufferLeft)
			n := copy(b[copyed:], r.buffer[r.bufferPoint:r.bufferPoint+willCopy])
			copyed += n
			r.bufferPoint += n
			continue
		}

		if r.rowNow > r.rowContent.RowCount {
			break
		}

		step := r.step * readTimes
		r.buffer, err = r.csvContent.line(r.rowContent, r.rowNow, misc.MinInt(r.rowNow+step, r.rowContent.RowCount+1))
		if err != nil {
			if column.IsDataOver(err) {
				r.rowNow = r.rowContent.RowCount + 1
				break // end the file
			} else {
				return 0, err
			}
		}
		readTimes += step
		r.bufferPoint = 0
		r.rowNow += step
	}

	if readTimes > r.step+1 {
		r.step = r.step*readTimes + 1
	}

	return
}
