package format

import (
	"errors"
	"fmt"
	"io"

	"github.com/tomcatzh/data-generator/data/row"
)

// Format is a data file type of data generator, such as csv or parquet
type Format interface {
	// Data return a reader in format
	Data(row *row.Set) (io.Reader, error)
}

// NewFormat returns a Format interface from template
func NewFormat(f map[string]interface{}) (Format, error) {
	fileType, ok := f["Type"].(string)
	if !ok || fileType == "" {
		return nil, errors.New("Template do not have Format type")
	}

	switch fileType {
	case "csv":
		return newCsv(f)
	default:
		return nil, fmt.Errorf("Unknown format: %v", fileType)
	}
}
