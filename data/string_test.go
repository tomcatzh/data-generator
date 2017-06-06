package data

import "testing"

func stringContins(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}

	return false
}

func TestFixString(t *testing.T) {
	const s = "test"

	d := newFixString(s, s)

	if d.Title() != s {
		t.Errorf("Title error: %v", d.Title())
	}

	result, _ := d.Data()

	if result != s {
		t.Errorf("Unexcapted value of fixString: %v", result)
	}
}

func TestEnumString(t *testing.T) {
	values := []string{"test1", "test2", "test3"}

	d := newEnumString("test", values)

	if d.Title() != "test" {
		t.Errorf("Title error: %v", d.Title())
	}

	result, _ := d.Data()
	found := stringContins(values, result)

	if !found {
		t.Errorf("Unexcapted value of enumString: %v", result)
	}
}
