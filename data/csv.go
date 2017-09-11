package data

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
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
	step        int
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (r *csvReader) Read(b []byte) (copyed int, err error) {
	max := len(b)
	if r.step == 0 {
		r.step = 1
	}

	bufferLength := len(r.buffer)
	bufferLeft := bufferLength - r.bufferPoint
	if r.rowNow > r.csvContent.rowCount && bufferLeft <= 0 {
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
			willCopy := min(max-copyed, bufferLeft)
			n := copy(b[copyed:], r.buffer[r.bufferPoint:r.bufferPoint+willCopy])
			copyed += n
			r.bufferPoint += n
			continue
		}

		if r.rowNow > r.csvContent.rowCount {
			break
		}

		step := r.step * readTimes
		r.buffer, err = r.csvContent.line(r.rowNow, min(r.rowNow+step, r.csvContent.rowCount+1))
		if err != nil {
			_, ok := err.(*TimeOver)

			if ok {
				r.rowNow = r.csvContent.rowCount + 1
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

const (
	csvNoCompress = iota
	csvGzipCompress
	csvZlibCompress
)

type csv struct {
	file
	haveTitleLine  bool
	delimiter      string
	quoteChar      string
	escapeChar     string
	lineTerminator string
	buffer         *bytes.Buffer
	compress       int
	compressLevel  int
}

func newCsv(rowCount int, namePart []namePart, delimiter string, quoteChar string, escapeChar string, lineTerminator string, haveTitleLine bool, compress int, compressLevel int) *csv {
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
		compress:       compress,
		compressLevel:  compressLevel,
	}
}

func (c *csv) Clone() FileData {
	b := make([]byte, 1024*1024)
	buf := bytes.NewBuffer(b)

	return &csv{
		file:           c.file,
		haveTitleLine:  c.haveTitleLine,
		delimiter:      c.delimiter,
		quoteChar:      c.quoteChar,
		escapeChar:     c.escapeChar,
		lineTerminator: c.lineTerminator,
		buffer:         buf,
		compress:       c.compress,
		compressLevel:  c.compressLevel,
	}
}

func (c *csv) line(rowStart int, rowMax int) ([]byte, error) {
	c.buffer.Reset()

	var err error
	var columns []string
	endNow := false

	for i := rowStart; i < rowMax; i++ {
		if i == 0 {
			columns = c.row.Title()
			i++
		} else {
			columns, err = c.row.Data()
			if err != nil {
				_, ok := err.(*TimeOver)

				if ok {
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
				c.buffer.WriteString(c.delimiter)
			}

			if c.escapeChar != "" {
				column = strings.Replace(column, c.quoteChar, c.escapeChar, -1)
			}

			c.buffer.WriteString(c.quoteChar)
			c.buffer.WriteString(column)
			c.buffer.WriteString(c.quoteChar)
		}

		if endNow {
			break
		}

		c.buffer.WriteString(c.lineTerminator)
	}
	return c.buffer.Bytes(), nil
}

func newGzipReader(source io.Reader, gzipLevel int) io.Reader {
	r, w := io.Pipe()
	go func() {
		defer w.Close()

		buffer := make([]byte, 1024)
		zip, err := gzip.NewWriterLevel(w, gzipLevel)
		defer zip.Close()
		if err != nil {
			w.CloseWithError(err)
		}

		io.CopyBuffer(zip, source, buffer)
	}()
	return r
}

func (c *csv) Data() (io.Reader, error) {
	var rowNow int

	if c.haveTitleLine {
		rowNow = 0
	} else {
		rowNow = 1
	}

	r := &csvReader{
		csvContent: c,
		rowNow:     rowNow,
	}

	switch c.compress {
	case csvGzipCompress:
		return newGzipReader(r, c.compressLevel), nil
	default:
		return r, nil
	}
}

func (c *csv) Save() (int64, error) {
	name, err := c.Name()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Name return an error: %v\n", err)
		return 0, err
	}
	data, err := c.Data()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Data return an error: %v\n", err)
		return 0, err
	}
	return c.storage.Save(name, data)
}
