package column

import "testing"

func stringContins(s []interface{}, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}

	return false
}

func TestFixString(t *testing.T) {
	const s = "test"

	tmpl := map[string]interface{}{}
	tmpl["Struct"] = "Fix"
	tmpl["Value"] = s

	df, err := newStringFactory(columnChangePerRow, tmpl)
	if err != nil {
		t.Errorf("Unexcapted error: %v", err)
	}
	d := df.Create()
	result, _ := d.Data()

	if result != s {
		t.Errorf("Unexcapted value of fixString: %v", result)
	}
}

func TestEnumString(t *testing.T) {
	values := []interface{}{"test1", "test2", "test3"}

	tmpl := map[string]interface{}{}
	tmpl["Struct"] = "Enum"
	tmpl["Values"] = values

	df, err := newStringFactory(columnChangePerRow, tmpl)
	if err != nil {
		t.Errorf("Unexcapted error: %v", err)
	}
	d := df.Create()
	result, _ := d.Data()
	found := stringContins(values, result)

	if !found {
		t.Errorf("Unexcapted value of enumString: %v", result)
	}
}
