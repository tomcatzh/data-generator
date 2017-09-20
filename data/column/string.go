package column

import (
	"fmt"
	"math/rand"
	"time"
)

func newStringFactory(columnMethod int, c map[string]interface{}) (result Factory, err error) {
	sstruct, ok := c["Struct"].(string)
	if !ok || sstruct == "" {
		return nil, fmt.Errorf("column does not have string struct")
	}

	switch sstruct {
	case "Fix":
		sfixValue, ok := c["Value"].(string)
		if !ok {
			return nil, fmt.Errorf("column does not have string fix value")
		}

		result = &fixString{value: sfixValue}
	case "Enum":
		sEnumValue, ok := c["Values"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("column does not have string enum values")
		}

		var sEnumString []string
		for _, v := range sEnumValue {
			sEnumString = append(sEnumString, v.(string))
		}

		result = &enumStringFactory{values: sEnumString}
	default:
		return nil, fmt.Errorf("Unexecpted string struct %v", sstruct)
	}

	return
}

type fixString struct {
	value string
}

func (s *fixString) Create() Column {
	return s
}

func (s *fixString) Data() (string, error) {
	return s.value, nil
}

type enumStringFactory struct {
	values []string
}

func (s *enumStringFactory) Create() Column {
	return &enumString{
		enumStringFactory: *s,
		rand:              rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type enumString struct {
	enumStringFactory
	rand *rand.Rand
}

func (s *enumString) Data() (string, error) {
	return s.values[s.rand.Intn(len(s.values))], nil
}
