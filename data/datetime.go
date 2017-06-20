package data

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var maxTime = time.Unix(1<<63-62135596801, 999999999)

type datetimeStep interface {
	Next(now time.Time) time.Time
	Clone() datetimeStep
}

type datetimeIncreaseStep struct {
	duration time.Duration
}

func (s *datetimeIncreaseStep) Clone() datetimeStep {
	return s
}

func (s *datetimeIncreaseStep) Next(now time.Time) time.Time {
	return now.Add(s.duration)
}

type datetimeRandomStep struct {
	durationN time.Duration
	durationA time.Duration
	rand      *rand.Rand
}

func (s *datetimeRandomStep) Clone() datetimeStep {
	return &datetimeRandomStep{
		durationN: s.durationN,
		durationA: s.durationA,
		rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *datetimeRandomStep) Next(now time.Time) time.Time {
	durationNano := s.rand.Int63n(s.durationN.Nanoseconds()) + s.durationA.Nanoseconds()
	return now.Add(time.Duration(durationNano))
}

type datetime struct {
	column
	format string
	now    time.Time
	end    time.Time
	step   datetimeStep
	buffer []byte
}

func (d *datetime) Data() (string, error) {
	stamp := d.step.Next(d.now)

	if stamp.After(d.end) {
		return "", errors.New("Arrived at the end time")
	}

	d.now = stamp

	return string(stamp.AppendFormat(d.buffer, d.format)), nil
}

func (d *datetime) Clone() columnData {
	const bufSize = 64
	var b []byte
	max := len(d.format) + 10
	if max < bufSize {
		var buf [bufSize]byte
		b = buf[:0]
	} else {
		b = make([]byte, 0, max)
	}

	result := &datetime{
		column: d.column,
		format: d.format,
		now:    d.now,
		end:    d.end,
		step:   d.step.Clone(),
		buffer: b,
	}

	return result
}

func newDatetimeIncrease(title string, format string, duration string, start time.Time, end time.Time) (*datetime, error) {
	durationNano, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}
	step := datetimeIncreaseStep{duration: durationNano}

	return &datetime{
		column: column{
			title: title,
		},
		format: format,
		now:    start,
		end:    end,
		step:   &step,
	}, nil
}

func newDatetimeRandom(title string, format string, unit string, max int, min int, start time.Time, end time.Time) (*datetime, error) {
	n := max - min
	if n <= 0 {
		return nil, errors.New("Max must bigger the min")
	}
	durationN, err := time.ParseDuration(fmt.Sprintf("%v%v", n, unit))
	if err != nil {
		return nil, err
	}
	durationA, err := time.ParseDuration(fmt.Sprintf("%v%v", min, unit))
	if err != nil {
		return nil, err
	}
	step := datetimeRandomStep{
		durationN: durationN,
		durationA: durationA,
	}

	return &datetime{
		column: column{
			title: title,
		},
		format: format,
		now:    start,
		end:    end,
		step:   &step,
	}, nil
}
