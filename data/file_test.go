package data

import (
	"testing"

	"github.com/tomcatzh/data-generator/data/row"
)

func newRowTmpl() (*row.Set, error) {
	rowTmpl := map[string]interface{}{}
	rowTmpl["RowCount"] = float64(2)

	sequence := []interface{}{}
	sequence = append(sequence, "A")
	sequence = append(sequence, "B")
	rowTmpl["Sequence"] = sequence

	data := map[string]interface{}{}
	dataA := map[string]interface{}{}
	dataA["Type"] = "String"
	dataA["Struct"] = "Fix"
	dataA["Value"] = "ABCDEFG"
	data["A"] = dataA
	dataB := map[string]interface{}{}
	dataB["Type"] = "String"
	dataB["Struct"] = "Fix"
	dataB["Value"] = "1234567"
	data["B"] = dataB
	rowTmpl["Data"] = data
	f, err := row.NewFactory(rowTmpl)
	if err != nil {
		return nil, err
	}
	s := f.Create()
	return s, nil
}

func TestFileName(t *testing.T) {
	s, err := newRowTmpl()
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}

	p, err := parseFileName("Here-${A}-${B}[2-3]")
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}

	n, err := name(p, s)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	if n != "Here-ABCDEFG-34" {
		t.Errorf("Unexcepted name: %v", n)
	}
}
