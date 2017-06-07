package data

import "testing"
import "strconv"

func TestRandomInterger(t *testing.T) {
	max := 10
	min := -2
	i := newRandomInteger("test", max, min)
	if i.Title() != "test" {
		t.Errorf("Unexcepted title: %v", i.Title())
	}

	c := i.Clone()

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
	f := newRandomFloat("test", max, min, 3)
	if f.Title() != "test" {
		t.Errorf("Unexcepted title: %v", f.Title())
	}

	c := f.Clone()
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
