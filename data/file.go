package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/tomcatzh/data-generator/data/format"
	"github.com/tomcatzh/data-generator/data/row"
	"github.com/tomcatzh/data-generator/data/storage"
)

// Factory is object to generate a file
type Factory struct {
	row       *row.Factory
	fileCount int
	namePart  []namePart
	format    format.Format
	storage   storage.Storage
}

// NewFactoryFile returns a file factory from template file path
func NewFactoryFile(filePath string) (*Factory, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	t := map[string]interface{}{}
	err = json.Unmarshal(content, &t)
	if err != nil {
		return nil, err
	}

	return NewFactory(t)
}

// NewFactory returns a file factory from template
func NewFactory(t map[string]interface{}) (result *Factory, err error) {
	result = &Factory{}

	tstorage, ok := t["Storage"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Template do not have Storage section")
	}
	result.storage, err = storage.NewStorage(tstorage)
	if err != nil {
		return nil, err
	}

	fileCount, ok := t["FileCount"].(float64)
	if !ok || fileCount == 0 {
		return nil, errors.New("Template do not have file count")
	}
	result.fileCount = int(fileCount)

	sformat, ok := t["Format"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Template do not have Format section")
	}
	result.format, err = format.NewFormat(sformat)
	if err != nil {
		return nil, err
	}

	file, ok := t["File"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Template do not have File section")
	}

	fileName, ok := file["Name"].(string)
	if !ok || fileName == "" {
		return nil, errors.New("Template do not have file name")
	}
	result.namePart, err = parseFileName(fileName)
	if err != nil {
		return nil, err
	}

	sRow, ok := file["Row"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Tempate have empty row section")
	}
	result.row, err = row.NewFactory(sRow)
	if err != nil {
		return nil, err
	}

	return
}

const (
	namePartTypeFix = iota
	namePartTypeData
	namePartTypeSubData
)

type namePart struct {
	partType       int
	value          string
	key            string
	substringStart int
	substringEnd   int
}

func parseFileName(fileName string) (part []namePart, err error) {
	isNormal := true
	isData := false
	isPreData := false
	isAfterData := false
	isSubData := false
	var s bytes.Buffer
	substringStart := 0
	substringEnd := 0
	var key string

	for _, c := range []byte(fileName) {
		if isNormal {
			if c == '$' {
				isNormal = false
				isPreData = true

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
		} else if isPreData {
			if c == '{' {
				isData = true
				isPreData = false
			} else {
				return nil, fmt.Errorf("Unexcepted char %v in name", string(c))
			}
		} else if isData {
			if c != '}' {
				s.WriteByte(c)
			} else {
				key = s.String()
				s.Reset()
				isData = false
				isAfterData = true
			}
		} else if isAfterData {
			if c != '[' {
				part = append(part, namePart{
					partType: namePartTypeData,
					key:      key,
				})
				s.WriteByte(c)
				isNormal = true
			} else {
				isSubData = true
			}

			isAfterData = false
		} else if isSubData {
			if c >= '0' && c <= '9' {
				s.WriteByte(c)
			} else if c == '-' {
				substringStart, err = strconv.Atoi(s.String())
				if err != nil {
					return nil, err
				}
				s.Reset()
			} else if c == ']' {
				substringEnd, err = strconv.Atoi(s.String())
				if err != nil {
					return nil, err
				}
				s.Reset()

				part = append(part, namePart{
					partType:       namePartTypeSubData,
					key:            key,
					substringStart: substringStart,
					substringEnd:   substringEnd,
				})

				substringStart = 0
				isSubData = false
				isNormal = true
			} else {
				return nil, fmt.Errorf("Unexcepted char %v in name", string(c))
			}
		}
	}

	if s.Len() > 0 {
		if isNormal {
			part = append(part, namePart{
				partType: namePartTypeFix,
				value:    s.String(),
			})
		} else {
			return nil, errors.New("Unexcepted name format")
		}
	}

	return
}

func name(namePart []namePart, r *row.Set) (string, error) {
	var result bytes.Buffer

	for _, p := range namePart {
		switch p.partType {
		case namePartTypeFix:
			result.WriteString(p.value)
		case namePartTypeData:
			s, err := r.SingleData(p.key)
			if err != nil {
				return "", err
			}
			result.WriteString(s)
		case namePartTypeSubData:
			s, err := r.SingleData(p.key)
			if err != nil {
				return "", err
			}
			result.WriteString(s[p.substringStart : p.substringEnd+1])
		}
	}

	return result.String(), nil
}

func (f *Factory) getFile() (*File, error) {
	s := f.row.Create()
	n, err := name(f.namePart, s)
	if err != nil {
		return nil, err
	}

	return &File{
		row:     s,
		format:  f.format,
		storage: f.storage,
		name:    n,
	}, nil
}

// Iterate returns a file iterate from factory for range
func (f *Factory) Iterate() <-chan *File {
	c := make(chan *File)
	go func() {
		defer close(c)
		for i := 0; i < f.fileCount; i++ {
			file, err := f.getFile()
			if err != nil {
				panic(err)
			}
			c <- file
		}
	}()
	return c
}

// File is the dedicate file with dedicate name, buffer and rand seed etc
type File struct {
	row     *row.Set
	format  format.Format
	storage storage.Storage
	name    string
}

// Save file with with specified name and data
func (f *File) Save() (int64, error) {
	data, err := f.format.Data(f.row)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Data return an error: %v\n", err)
		return 0, err
	}
	return f.storage.Save(f.name, data)
}
