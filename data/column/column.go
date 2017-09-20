package column

import (
	"fmt"
)

const (
	columnChangePerFile = iota
	columnChangePerRow
	columnChangePerRowAndFile
)

// Factory is a interface of column factory
type Factory interface {
	Create() Column
}

// NewFactory returns a column factory from template
func NewFactory(c map[string]interface{}) (result Factory, err error) {
	ctype, ok := c["Type"].(string)
	if !ok || ctype == "" {
		return nil, fmt.Errorf("column does not have type")
	}

	var columnMethod int
	cchange, ok := c["Change"].(string)
	if !ok || cchange == "" {
		columnMethod = columnChangePerRow
	} else {
		switch cchange {
		case "PerFile":
			columnMethod = columnChangePerFile
		case "PerRow":
			columnMethod = columnChangePerRow
		case "PerRowAndFile":
			columnMethod = columnChangePerRowAndFile
		default:
			return nil, fmt.Errorf("Unexecpted colomn change method in column: %v", cchange)
		}
	}

	switch ctype {
	case "Datetime":
		result, err = newDatetimeFactory(columnMethod, c)
	case "String":
		result, err = newStringFactory(columnMethod, c)
	case "Numeric":
		result, err = newNumericFactory(columnMethod, c)
	case "IPv4":
		result, err = newIPv4Factory(columnMethod, c)
	}
	if err != nil {
		return nil, err
	}

	return
}

// Column is a interface of data column
type Column interface {
	Data() (string, error)
}
