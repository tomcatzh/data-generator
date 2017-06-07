package data

import (
	"math/rand"
	"time"
)

type fixString struct {
	column
	value string
}

func (s *fixString) Data() (string, error) {
	return s.value, nil
}

func (s *fixString) Clone() columnData {
	return s
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
	rand   *rand.Rand
}

func (s *enumString) Data() (string, error) {
	return s.values[rand.Intn(len(s.values))], nil
}

func (s *enumString) Clone() columnData {
	return &enumString{
		column: s.column,
		values: s.values,
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func newEnumString(title string, values []string) *enumString {
	return &enumString{
		column: column{
			title: title,
		},
		values: values,
	}
}
