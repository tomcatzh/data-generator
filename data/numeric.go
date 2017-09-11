package data

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

type randomInteger struct {
	column
	n    int
	a    int
	rand *rand.Rand
}

func (i *randomInteger) Clone() columnData {
	return &randomInteger{
		column: i.column,
		n:      i.n,
		a:      i.a,
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (i *randomInteger) Data() (string, error) {
	return strconv.Itoa(i.rand.Intn(i.n) + i.a), nil
}

func newRandomInteger(title string, max, min int) *randomInteger {
	return &randomInteger{
		column: column{
			title: title,
		},
		n: max - min,
		a: min,
	}
}

type randomFloat struct {
	column
	n    int64
	a    int64
	mod  int64
	rand *rand.Rand
}

func (f *randomFloat) Clone() columnData {
	return &randomFloat{
		column: f.column,
		n:      f.n,
		a:      f.a,
		mod:    f.mod,
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (f *randomFloat) Data() (string, error) {
	return strconv.FormatFloat(float64(f.rand.Int63n(f.n)+f.a)/float64(f.mod), 'f', 6, 64), nil
}

func newRandomFloat(title string, max float64, min float64, decimal int) *randomFloat {
	mod := int64(math.Pow10(decimal))
	maxI := int64(max * float64(mod))
	minI := int64(min * float64(mod))

	return &randomFloat{
		column: column{
			title: title,
		},
		n:   maxI - minI,
		a:   minI,
		mod: mod,
	}
}
