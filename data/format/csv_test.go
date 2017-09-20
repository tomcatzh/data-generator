package format

import (
	"io/ioutil"
	"testing"

	"github.com/tomcatzh/data-generator/data/row"
)

func newCsvTmpl() (*csv, error) {
	csvTmpl := map[string]interface{}{}
	csvTmpl["HaveTitleLine"] = true

	return newCsv(csvTmpl)
}

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
	dataA["Value"] = "This's A"
	data["A"] = dataA
	dataB := map[string]interface{}{}
	dataB["Type"] = "String"
	dataB["Struct"] = "Fix"
	dataB["Value"] = "This's B"
	data["B"] = dataB
	rowTmpl["Data"] = data
	f, err := row.NewFactory(rowTmpl)
	if err != nil {
		return nil, err
	}
	s := f.Create()
	return s, nil
}

func TestCsvReader(t *testing.T) {
	c, err := newCsvTmpl()
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}

	s, err := newRowTmpl()
	if err != nil {
		t.Errorf("Unexcpted error: %v", err)
	}

	reader, err := c.Data(s)
	if err != nil {
		t.Errorf("Unexcpted error: %v", err)
	}
	line, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Errorf("Unexcpted error: %v", err)
	}
	if string(line) != "\"A\",\"B\"\r\n\"This's A\",\"This's B\"\r\n\"This's A\",\"This's B\"\r\n" {
		t.Errorf("Unexcepted line data: \r\n%v", string(line))
	}
}

func TestCsvLine(t *testing.T) {
	c, err := newCsvTmpl()
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}

	s, err := newRowTmpl()
	if err != nil {
		t.Errorf("Unexcpted error: %v", err)
	}

	line, err := c.line(s, 0, 2)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	if string(line) != "\"A\",\"B\"\r\n\"This's A\",\"This's B\"\r\n" {
		t.Errorf("Unexcepted line data: \r\n%v", string(line))
	}
}
