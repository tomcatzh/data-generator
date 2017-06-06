package data

import "math/rand"

type fixString struct {
	column
	value string
}

func (s *fixString) Data() (string, error) {
	return s.value, nil
}

func (s *fixString) Clone() columnData {
	return newFixString(s.title, s.value)
}

func newFixString(title string, value string) *fixString {
	return &fixString{
		column: column{
			title: title,
		},
		value: value,
	}
}

type enumString struct {
	column
	values []string
}

func (s *enumString) Data() (string, error) {
	return s.values[rand.Intn(len(s.values))], nil
}

func (s *enumString) Clone() columnData {
	return newEnumString(s.title, s.values)
}

func newEnumString(title string, values []string) *enumString {
	return &enumString{
		column: column{
			title: title,
		},
		values: values,
	}
}
