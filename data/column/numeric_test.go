package column

import (
	"strconv"
	"testing"
)

func TestRandomInterger(t *testing.T) {
	max := 10
	min := -2

	tmpl := map[string]interface{}{}
	tmpl["Format"] = "Integer"
	step := map[string]interface{}{}
	step["Type"] = "Random"
	step["Max"] = float64(max)
	step["Min"] = float64(min)
	tmpl["Step"] = step

	i, err := newNumericFactory(columnChangePerRow, tmpl)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}

	c := i.Create()
	for j := 0; j < 30; j++ {
		s, _ := c.Data()
		k, err := strconv.Atoi(s)
		if err != nil {
			t.Errorf("Unexcepted error: %v [%v]", err, s)
		} else if k < min || k > max {
			t.Errorf("Unexcetped data: %v", k)
		}
	}
}

func TestRandomFloat(t *testing.T) {
	max := 10.2054
	min := -2.341

	tmpl := map[string]interface{}{}
	tmpl["Format"] = "Float"
	step := map[string]interface{}{}
	step["Type"] = "Random"
	step["Max"] = float64(max)
	step["Min"] = float64(min)
	step["Decimal"] = 3
	tmpl["Step"] = step

	f, err := newNumericFactory(columnChangePerRow, tmpl)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}

	c := f.Create()
	for i := 0; i < 30; i++ {
		s, _ := c.Data()
		n, err := strconv.ParseFloat(s, 32)
		if err != nil {
			t.Errorf("Unexcepted error: %v [%v]", err, s)
		} else if n < min || n > max {
			t.Errorf("Unexcetped data: %v", n)
		}
	}
}
