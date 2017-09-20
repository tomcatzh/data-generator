package row

import (
	"bytes"
	"fmt"

	"github.com/tomcatzh/data-generator/data/column"
)

const defaultRowCount = 1000

// Factory is object to generate a delicate rowset
type Factory struct {
	titles         []string
	columnFactorys map[string]column.Factory
	rowCount       int
}

// NewFactory returns a row set factory from template
func NewFactory(c map[string]interface{}) (result *Factory, err error) {
	result = &Factory{}

	rowCount, ok := c["RowCount"].(float64)
	if !ok || rowCount <= 0 {
		rowCount = defaultRowCount
	}
	result.rowCount = int(rowCount)

	sSequence, ok := c["Sequence"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("file does not have row sequence")
	}
	for _, v := range sSequence {
		result.titles = append(result.titles, v.(string))
	}

	result.columnFactorys = map[string]column.Factory{}
	r, ok := c["Data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("file does not have row data")
	}
	for key, data := range r {
		c, ok := data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Unexecpted column section in key: %v", key)
		}
		column, err := column.NewFactory(c)
		if err != nil {
			return nil, fmt.Errorf("%v: %v", key, err)
		}
		result.columnFactorys[key] = column
	}

	return
}

// Create returns a set of template row
func (f *Factory) Create() *Set {
	columns := map[string]column.Column{}

	for key, factory := range f.columnFactorys {
		columns[key] = factory.Create()
	}

	b := make([]byte, 1024*1024)
	buf := bytes.NewBuffer(b)

	return &Set{
		columns:  columns,
		titles:   f.titles,
		Buffer:   buf,
		RowCount: f.rowCount,
	}
}

// Set is a columnset of data
type Set struct {
	columns  map[string]column.Column
	titles   []string
	Buffer   *bytes.Buffer
	RowCount int
}

// Title return the title set of row set
func (r *Set) Title() []string {
	return r.titles
}

// SingleData returns the next data of special key
func (r *Set) SingleData(key string) (string, error) {
	return r.columns[key].Data()
}

// Data return the next row data of row set
func (r *Set) Data() ([]string, error) {
	result := []string{}

	for _, title := range r.titles {
		data, err := r.columns[title].Data()
		if err != nil {
			return nil, err
		}
		result = append(result, data)
	}

	return result, nil
}
