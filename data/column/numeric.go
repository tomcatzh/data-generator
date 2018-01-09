package column

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func newNumericFactory(columnMethod int, c map[string]interface{}) (result Factory, err error) {
	nformat, ok := c["Format"].(string)
	if !ok || nformat == "" {
		return nil, fmt.Errorf("column does not have numeric format")
	}

	nstep, ok := c["Step"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("column does not have numeric step section")
	}

	nStepType, ok := nstep["Type"].(string)
	if !ok {
		return nil, fmt.Errorf("column does not have numeric step type")
	}

	switch nStepType {
	case "Random":
		niMax, ok := nstep["Max"].(float64)
		if !ok {
			return nil, fmt.Errorf("column does not have numeric random max")
		}
		niMin, ok := nstep["Min"].(float64)
		if !ok {
			return nil, fmt.Errorf("column does not have numeric random min")
		}

		switch nformat {
		case "Integer":
			result = newRandomInteger(columnMethod, int(niMax), int(niMin))
		case "Float":
			niDecimal, ok := nstep["Decimal"].(float64)
			if !ok {
				return nil, fmt.Errorf("column does not have numeric random float decimal")
			}
			result = newRandomFloat(columnMethod, niMax, niMin, int(niDecimal))
		default:
			return nil, fmt.Errorf("Unexecpted numeric format %v", nformat)
		}
	default:
		return nil, fmt.Errorf("Unexecpted numeric step type %v", nStepType)
	}

	return
}

type randomIntegerFactory struct {
	n int
	a int
}

func newRandomInteger(columnMethod, max, min int) *randomIntegerFactory {
	return &randomIntegerFactory{
		n: max - min,
		a: min,
	}
}

func (i *randomIntegerFactory) Create() Column {
	return &randomInteger{
		randomIntegerFactory: *i,
		rand:                 rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type randomInteger struct {
	randomIntegerFactory
	rand *rand.Rand
}

func (i *randomInteger) Data() (string, error) {
	return strconv.Itoa(i.rand.Intn(i.n) + i.a), nil
}

type randomFloatFactory struct {
	n       int64
	a       int64
	mod     int64
	decimal int
}

func (f *randomFloatFactory) Create() Column {
	return &randomFloat{
		randomFloatFactory: *f,
		rand:               rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type randomFloat struct {
	randomFloatFactory
	rand *rand.Rand
}

func (f *randomFloat) Data() (string, error) {
	return strconv.FormatFloat(float64(f.rand.Int63n(f.n)+f.a)/float64(f.mod), 'f', f.decimal, 64), nil
}

func newRandomFloat(columnMethod int, max float64, min float64, decimal int) *randomFloatFactory {
	mod := int64(math.Pow10(decimal))
	maxI := int64(max * float64(mod))
	minI := int64(min * float64(mod))

	return &randomFloatFactory{
		n:       maxI - minI,
		a:       minI,
		mod:     mod,
		decimal: decimal,
	}
}
