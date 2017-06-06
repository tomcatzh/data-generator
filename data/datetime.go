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
}

type datetimeFixStep struct {
	duration time.Duration
}

func (s *datetimeFixStep) Next(now time.Time) time.Time {
	return now.Add(s.duration)
}

type datetimeRandomStep struct {
	durationN time.Duration
	durationA time.Duration
}

func (s *datetimeRandomStep) Next(now time.Time) time.Time {
	durationNano := rand.Int63n(s.durationN.Nanoseconds()) + s.durationA.Nanoseconds()

	return now.Add(time.Duration(durationNano))
}

type datetime struct {
	column
	format string
	now    time.Time
	end    time.Time
	step   datetimeStep
}

func (d *datetime) Data() (string, error) {
	stamp := d.step.Next(d.now)

	if stamp.After(d.end) {
		return "", errors.New("Arrived at the end time")
	}

	d.now = stamp

	return stamp.Format(d.format), nil
}

func (d *datetime) Clone() columnData {
	return &datetime{
		column: d.column,
		format: d.format,
		now:    d.now,
		end:    d.end,
		step:   d.step,
	}
}

func newDatetimeFix(title string, format string, duration string, start time.Time, end time.Time) (*datetime, error) {
	durationNano, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}
	step := datetimeFixStep{duration: durationNano}

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
