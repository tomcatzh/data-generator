package row

import "testing"

func TestRow(t *testing.T) {
	tmpl := map[string]interface{}{}

	sequence := []interface{}{}
	sequence = append(sequence, "A")
	sequence = append(sequence, "B")
	tmpl["Sequence"] = sequence

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
	tmpl["Data"] = data

	f, err := NewFactory(tmpl)
	if err != nil {
		t.Errorf("Unexcpted error: %v", err)
	}
	s := f.Create()
	titles := s.Title()
	if titles[0] != "A" || titles[1] != "B" {
		t.Errorf("Unexcpted titles: %v", titles)
	}

	d, err := s.Data()
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	if d[0] != "This's A" || d[1] != "This's B" {
		t.Errorf("Unexcepted data: %v", data)
	}
}
