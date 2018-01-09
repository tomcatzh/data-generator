package column

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

type datetimeFactory struct {
	format       string
	now          time.Time
	end          time.Time
	step         datetimeStep
	fileStep     *datetimeIncreaseStep
	fileNow      time.Time
	columnMethod int
}

func fileStep(duration string) (time.Duration, error) {
	if duration == "" {
		return 0, nil
	}

	return time.ParseDuration(duration)
}

func newDatetimeIncrease(format string, duration string, start time.Time, end time.Time, fileDuration string, columnMehtod int) (*datetimeFactory, error) {
	fDuration, err := fileStep(fileDuration)
	if err != nil {
		return nil, err
	}

	durationNano, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}

	return &datetimeFactory{
		format:       format,
		now:          start.Add(-durationNano),
		end:          end,
		step:         &datetimeIncreaseStep{duration: durationNano},
		fileStep:     &datetimeIncreaseStep{duration: fDuration},
		fileNow:      start.Add(-fDuration),
		columnMethod: columnMehtod,
	}, nil
}

func newDatetimeRandom(format string, unit string, max int, min int, start time.Time, end time.Time, fileDuration string, columnMehtod int) (*datetimeFactory, error) {
	fDuration, err := fileStep(fileDuration)
	if err != nil {
		return nil, err
	}

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

	return &datetimeFactory{
		format: format,
		now:    start,
		end:    end,
		step: &datetimeRandomStep{
			durationN: durationN,
			durationA: durationA,
		},
		fileStep:     &datetimeIncreaseStep{duration: fDuration},
		fileNow:      start.Add(-fDuration),
		columnMethod: columnMehtod,
	}, nil
}

func newDatetimeFactory(columnMethod int, c map[string]interface{}) (*datetimeFactory, error) {
	var result *datetimeFactory

	dformat, ok := c["Format"].(string)
	if !ok || dformat == "" {
		dformat = time.RFC3339
	}

	dstep, ok := c["Step"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("column does not have datetime step section")
	}

	var dFileDuration string
	if columnMethod == columnChangePerRowAndFile || columnMethod == columnChangePerFile {
		var dfileStep map[string]interface{}
		if columnMethod == columnChangePerRowAndFile {
			dfileStep, ok = c["FileStep"].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("column does not have file step section")
			}

			dFileDuration, ok = dfileStep["Duration"].(string)
			if !ok || dFileDuration == "" {
				return nil, fmt.Errorf("column does not have file duration")
			}
		} else {
			dFileDuration, ok = dstep["Duration"].(string)
			if !ok || dFileDuration == "" {
				return nil, fmt.Errorf("column does not have file duration")
			}
		}
	}

	dstartString, ok := dstep["Start"].(string)
	if !ok || dstartString == "" {
		return nil, fmt.Errorf("column does not have datetime start stamp")
	}

	dstart, err := time.Parse(dformat, dstartString)
	if err != nil {
		return nil, err
	}

	var dend time.Time
	dendString, ok := dstep["End"].(string)
	if !ok || dendString == "" {
		dend = maxTime
	} else {
		dend, err = time.Parse(dformat, dendString)
		if err != nil {
			return nil, err
		}
	}

	dstepType, ok := dstep["Type"].(string)
	if !ok || dstepType == "" {
		return nil, fmt.Errorf("column does not have datetime step type")
	}

	switch dstepType {
	case "Increase":
		dstepDuration, ok := dstep["Duration"].(string)
		if !ok || dstepDuration == "" {
			return nil, fmt.Errorf("column does not have datetime fix duration")
		}

		result, err = newDatetimeIncrease(dformat, dstepDuration, dstart, dend, dFileDuration, columnMethod)
		if err != nil {
			return nil, err
		}
	case "Random":
		dstepUnit, ok := dstep["Unit"].(string)
		if !ok || dstepUnit == "" {
			return nil, fmt.Errorf("column does not have datetime random unit")
		}

		dstepMax, ok := dstep["Max"].(float64)
		if !ok {
			return nil, fmt.Errorf("column does not have datetime random max")
		}

		dstepMin, ok := dstep["Min"].(float64)
		if !ok {
			return nil, fmt.Errorf("column does not have datetime random min")
		}

		result, err = newDatetimeRandom(dformat, dstepUnit, int(dstepMax), int(dstepMin), dstart, dend, dFileDuration, columnMethod)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (d *datetimeFactory) Create() Column {
	const bufSize = 64
	var b []byte
	max := len(d.format) + 10
	if max < bufSize {
		var buf [bufSize]byte
		b = buf[:0]
	} else {
		b = make([]byte, 0, max)
	}

	end := d.end

	if d.fileStep != nil && d.fileStep.duration > 0 {
		d.fileNow = d.fileStep.Next(d.fileNow)
		end = d.fileStep.Next(d.fileNow)
	}

	var result Column
	if d.columnMethod == columnChangePerFile {
		if d.fileStep != nil && d.fileStep.duration > 0 {
			d.now = d.fileStep.Next(d.now)
		}

		result = &fixString{
			value: d.now.Format(d.format),
		}
	} else {
		result = &datetime{
			datetimeFactory: datetimeFactory{
				format:  d.format,
				now:     d.now,
				end:     end,
				step:    d.step.Clone(),
				fileNow: d.fileNow,
			},
			buffer: b,
		}

		if d.fileStep != nil && d.fileStep.duration > 0 {
			d.now = d.fileStep.Next(d.now)
		}
	}

	return result
}

type datetime struct {
	datetimeFactory
	buffer []byte
}

func (d *datetime) Data() (string, error) {
	stamp := d.step.Next(d.now)

	if stamp.After(d.end) {
		return "", NewDataOver(fmt.Sprintf("Time is over now[%v] end[%v]", stamp.Format(time.ANSIC), d.end.Format(time.ANSIC)))
	}

	d.now = stamp

	return string(stamp.AppendFormat(d.buffer, d.format)), nil
}
