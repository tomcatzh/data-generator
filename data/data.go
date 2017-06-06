package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/tomcatzh/data-generator/awsstorage"
)

const defaultRowCount = 1000

// A Template is a handler of data generator template
type Template struct {
	file      FileData
	fileCount int
	storage   storage
	columns   []templateColumn
}

func (t *Template) getFile() (FileData, error) {
	result := t.file.Clone()

	for _, c := range t.columns {
		switch c.columnMethod {
		case columnChangePerFile:
			data, err := c.column.Data()
			if err != nil {
				return nil, err
			}
			s := newFixString(c.column.Title(), data)
			result.AddColumn(s)
		case columnChangePerRow:
			result.AddColumn(c.column.Clone())
		}
	}

	result.SetStorage(t.storage)

	return result, nil
}

// Iterate returns FileData rangeable list
func (t *Template) Iterate() <-chan FileData {
	c := make(chan FileData)
	go func() {
		for i := 0; i < t.fileCount; i++ {
			file, err := t.getFile()
			if err != nil {
				panic(err)
			}
			c <- file
		}
		close(c)
	}()
	return c
}

const (
	columnChangePerFile = iota
	columnChangePerRow
)

type templateColumn struct {
	columnMethod int
	column       columnData
}

// FileData is a handler of file for data generator
type FileData interface {
	Data() (io.Reader, error)
	Name() (string, error)
	AddColumn(column columnData)
	Clone() FileData
	SetStorage(s storage)
	Save() (int64, error)
}

const (
	namePartTypeFix = iota
	namePartTypeData
	namePartTypeSubData
)

type namePart struct {
	partType  int
	value     string
	index     int
	substring int
}

type file struct {
	row      row
	rowCount int
	namePart []namePart
	storage  storage
}

func (f *file) SetStorage(s storage) {
	f.storage = s
}

func (f *file) Name() (string, error) {
	var result bytes.Buffer

	for _, p := range f.namePart {
		switch p.partType {
		case namePartTypeFix:
			result.WriteString(p.value)
		case namePartTypeData:
			s, err := f.row.columns[p.index].Data()
			if err != nil {
				return "", err
			}
			result.WriteString(s)
		case namePartTypeSubData:
			s, err := f.row.columns[p.index].Data()
			if err != nil {
				return "", err
			}
			result.WriteString(s[0:p.substring])
		}
	}

	return result.String(), nil
}

// NewTemplate returns a template handler from template path.
func NewTemplate(templateFile string) (*Template, error) {
	templateContent, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return nil, err
	}

	var result Template

	t := map[string]interface{}{}
	err = json.Unmarshal(templateContent, &t)
	if err != nil {
		return nil, err
	}

	fileCount, ok := t["FileCount"].(float64)
	if !ok || fileCount == 0 {
		return nil, errors.New("Template do not have file count")
	}

	result.fileCount = int(fileCount)

	format, ok := t["Format"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Template do not have Format section")
	}

	fileType, ok := format["Type"].(string)
	if !ok || fileType == "" {
		return nil, errors.New("Template do not have Format type")
	}

	file, ok := t["File"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Template do not have File section")
	}

	rowCount, ok := file["RowCount"].(float64)
	if !ok || rowCount <= 0 {
		rowCount = defaultRowCount
	}

	fileName, ok := file["Name"].(string)
	if !ok || fileName == "" {
		return nil, errors.New("Template do not have file name")
	}

	isNormal := true
	isData := false
	isSubData := false
	var s bytes.Buffer
	var index int
	var part []namePart
	for _, c := range []byte(fileName) {
		if isNormal {
			if c == '$' {
				isNormal = false
				isData = true

				if s.Len() > 0 {
					part = append(part, namePart{
						partType: namePartTypeFix,
						value:    s.String(),
					})
				}

				s.Reset()
				continue
			} else {
				s.WriteByte(c)
			}
		} else if isData {
			if c >= '0' && c <= '9' {
				s.WriteByte(c)
			} else {
				index, err = strconv.Atoi(s.String())
				if err != nil {
					return nil, err
				}
				s.Reset()
				isData = false

				if c != '[' {
					part = append(part, namePart{
						partType: namePartTypeData,
						index:    index,
					})
					s.WriteByte(c)
					isNormal = true
				} else {
					isSubData = true
				}
			}
		} else if isSubData {
			if c >= '0' && c <= '9' {
				s.WriteByte(c)
			} else if c == ']' {
				substring, err := strconv.Atoi(s.String())
				if err != nil {
					return nil, err
				}
				s.Reset()

				part = append(part, namePart{
					partType:  namePartTypeSubData,
					index:     index,
					substring: substring,
				})

				isSubData = false
				isNormal = true
			} else {
				return nil, fmt.Errorf("Unexcepted char %v in name", c)
			}
		}
	}

	if s.Len() > 0 {
		if isNormal {
			part = append(part, namePart{
				partType: namePartTypeFix,
				value:    s.String(),
			})
		} else if isData {
			index, err = strconv.Atoi(s.String())
			if err != nil {
				return nil, err
			}
			part = append(part, namePart{
				partType: namePartTypeData,
				index:    index,
			})
		} else if isSubData {
			return nil, errors.New("Unexcepted name format")
		}
	}

	switch fileType {
	case "csv":
		haveTitleLine, ok := format["haveTitleLine"].(bool)

		delimiter, ok := format["delimiter"].(string)
		if !ok || delimiter == "" {
			delimiter = defaultDelimiter
		}
		quotechar, ok := format["quotechar"].(string)
		if !ok || quotechar == "" {
			quotechar = defaultQuoteChar
		}
		escapechar, ok := format["escapechar"].(string)
		if !ok || escapechar == "" {
			escapechar = defaultEscapeChar
		}
		lineterminator, ok := format["lineterminator"].(string)
		if !ok || lineterminator == "" {
			lineterminator = defaultLineTeriminator
		}
		csv := newCsv(int(rowCount), part, delimiter, quotechar, escapechar, lineterminator, haveTitleLine)

		result.file = csv
	default:
		return nil, fmt.Errorf("Unknown format: %v", fileType)
	}

	data, ok := file["Data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Tempate have empty data")
	}
	for i := 0; ; i++ {
		var cdata templateColumn

		idx := strconv.Itoa(i)

		c, ok := data[idx].(map[string]interface{})
		if !ok {
			break
		}

		ctitle, ok := c["Title"].(string)
		if !ok || ctitle == "" {
			return nil, fmt.Errorf("%v column does not have column", i)
		}

		cchange, ok := c["Change"].(string)
		if !ok || cchange == "" {
			cdata.columnMethod = columnChangePerRow
		} else {
			switch cchange {
			case "PerFile":
				cdata.columnMethod = columnChangePerFile
			case "PerRow":
				cdata.columnMethod = columnChangePerRow
			default:
				return nil, fmt.Errorf("Unexecpted colomn change method in column %v: %v", i, cchange)
			}
		}

		ctype, ok := c["Type"].(string)
		if !ok || ctype == "" {
			return nil, fmt.Errorf("%v column does not have type", i)
		}

		switch ctype {
		case "datetime":
			dformat, ok := c["Format"].(string)
			if !ok || dformat == "" {
				dformat = time.RFC3339
			}

			dstartString, ok := c["Start"].(string)
			if !ok || dstartString == "" {
				return nil, fmt.Errorf("%v column does not have datetime start stamp", i)
			}

			dstart, err := time.Parse(dformat, dstartString)
			if err != nil {
				return nil, err
			}

			var dend time.Time
			dendString, ok := c["End"].(string)
			if !ok || dendString == "" {
				dend = maxTime
			} else {
				dend, err = time.Parse(dformat, dendString)
				if err != nil {
					return nil, err
				}
			}

			dstep, ok := c["Step"].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("%v column does not have datetime step section", i)
			}

			dstepType, ok := dstep["Type"].(string)
			if !ok || dstepType == "" {
				return nil, fmt.Errorf("%v column does not have datetime step type", i)
			}

			switch dstepType {
			case "Fix":
				dstepDuration, ok := dstep["Duration"].(string)
				if !ok || dstepDuration == "" {
					return nil, fmt.Errorf("%v column does not have datetime fix duration", i)
				}

				cdata.column, err = newDatetimeFix(ctitle, dformat, dstepDuration, dstart, dend)
				if err != nil {
					return nil, err
				}
			case "Random":
				dstepUnit, ok := dstep["Unit"].(string)
				if !ok || dstepUnit == "" {
					return nil, fmt.Errorf("%v column does not have datetime random unit", i)
				}

				dstepMax, ok := dstep["Max"].(float64)
				if !ok {
					return nil, fmt.Errorf("%v column does not have datetime random max", i)
				}

				dstepMin, ok := dstep["Min"].(float64)
				if !ok {
					return nil, fmt.Errorf("%v column does not have datetime random min", i)
				}

				cdata.column, err = newDatetimeRandom(ctitle, dformat, dstepUnit, int(dstepMax), int(dstepMin), dstart, dend)
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("Unexecpted datetime step type in column %v: %v", i, dstepType)
			}
		case "string":
			sstruct, ok := c["Struct"].(string)
			if !ok || sstruct == "" {
				return nil, fmt.Errorf("%v column does not have string struct", i)
			}

			switch sstruct {
			case "fix":
				sfixValue, ok := c["Value"].(string)
				if !ok {
					return nil, fmt.Errorf("%v column does not have string fix value", i)
				}

				cdata.column = newFixString(ctitle, sfixValue)
			case "enum":
				sEnumValue, ok := c["Value"].([]interface{})
				if !ok {
					return nil, fmt.Errorf("%v column does not have string enum value", i)
				}

				var sEnumString []string
				for _, v := range sEnumValue {
					sEnumString = append(sEnumString, v.(string))
				}

				cdata.column = newEnumString(ctitle, sEnumString)
			default:
				return nil, fmt.Errorf("Unexecpted string struct in column %v: %v", i, sstruct)
			}
		default:
			return nil, fmt.Errorf("Unexecpted column type in column %v: %v", i, ctype)
		}

		storage, ok := t["Storage"].(map[string]interface{})
		if !ok {
			return nil, errors.New("Template do not have Storage section")
		}

		stype, ok := storage["Type"].(string)
		if !ok || stype == "" {
			return nil, errors.New("Template do not have Storage type")
		}

		switch stype {
		case "Local":
			spath, ok := storage["Path"].(string)
			if !ok || spath == "" {
				spath = "."
			}
			result.storage = newStorageLocal(spath)
		case "S3":
			sregion, ok := storage["Region"].(string)
			if !ok || sregion == "" {
				return nil, errors.New("Template do not have S3 region")
			}
			sbucket, ok := storage["Bucket"].(string)
			if !ok || sbucket == "" {
				return nil, errors.New("Template do not have S3 bucket")
			}
			result.storage = awsstorage.NewStorageS3(sregion, sbucket)
		default:
			return nil, fmt.Errorf("Unexecepted Storage type: %v", stype)
		}

		result.columns = append(result.columns, cdata)
	}

	return &result, nil
}

type storage interface {
	Save(key string, reader io.Reader) (int64, error)
}

func (f *file) AddColumn(column columnData) {
	f.row.AddColumn(column)
}

type row struct {
	columns []columnData
}

func (r *row) AddColumn(column columnData) {
	r.columns = append(r.columns, column)
}

func (r *row) Title() []string {
	result := []string{}

	for _, column := range r.columns {
		data := column.Title()
		result = append(result, data)
	}

	return result
}

func (r *row) Data() ([]string, error) {
	result := []string{}

	for _, column := range r.columns {
		data, err := column.Data()
		if err != nil {
			return nil, err
		}
		result = append(result, data)
	}

	return result, nil
}

type column struct {
	title string
}

type columnData interface {
	Title() string
	Data() (string, error)
	Clone() columnData
}

func (c *column) Title() string {
	return c.title
}
