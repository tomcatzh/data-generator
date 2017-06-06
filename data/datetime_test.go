package data

import (
	"math/rand"
	"testing"
	"time"
)

func TestDatetimeRandom(t *testing.T) {
	stamp := time.Now()

	_, err := newDatetimeRandom("test", time.UnixDate, "m", 1, 1, stamp, maxTime)
	if err == nil {
		t.Errorf("Excepted an error, but nil")
	}

	for i := 0; i < 10; i++ {
		data1, err := newDatetimeRandom("test", time.UnixDate, "m", 10, -2, stamp, maxTime)
		if err != nil {
			t.Errorf("Someting wrong when new datetime: %v", err)
		} else if data1.title != "test" {
			t.Errorf("Title error: %v", data1.Title())
		}
		result1, err := data1.Data()
		if err != nil {
			t.Errorf("Something wrong when get datetime data: %v", err)
		}
		resultStamp1, _ := time.Parse(time.UnixDate, result1)
		minutes := resultStamp1.Sub(stamp) / time.Minute
		if minutes > 10 || minutes < -2 {
			t.Errorf("Unexcepted duration: %v mins", minutes)
			return
		}
	}

	for i := 0; i < 10; i++ {
		data2, err := newDatetimeRandom("test", time.UnixDate, "m", 10, -2, stamp, stamp.Truncate(5*time.Minute))
		if err != nil {
			t.Errorf("Someting wrong when new datetime: %v", err)
		} else if data2.title != "test" {
			t.Errorf("Title error: %v", data2.Title())
		}
		_, err = data2.Data()
		if err == nil {
			t.Errorf("Excepted an error, but nil")
		}
	}
}

func TestDatetimeFix(t *testing.T) {
	stamp, _ := time.Parse(time.UnixDate, "Thu Jun 1 17:00:00 CST 2017")
	data1, err := newDatetimeFix("test", time.UnixDate, "3m", stamp, maxTime)
	if err != nil {
		t.Errorf("Someting wrong when new datetime: %v", err)
	} else if data1.title != "test" {
		t.Errorf("Title error: %v", data1.Title())
	}
	result1, err := data1.Data()
	if err != nil {
		t.Errorf("Something wrong when get datetime data: %v", err)
	} else if result1 != "Thu Jun  1 17:03:00 CST 2017" {
		t.Errorf("Unexcepted result: %s", result1)
	}

	stampEnd, _ := time.Parse(time.UnixDate, "Thu Jun 1 17:02:00 CST 2017")
	data2, _ := newDatetimeFix("test", time.UnixDate, "3m", stamp, stampEnd)
	if err != nil {
		t.Errorf("Someting wrong when new datetime: %v", err)
	}
	_, err = data2.Data()
	if err == nil {
		t.Error("Excepted an error, but nil")
	}
}

func TestStepFix(t *testing.T) {
	duration, _ := time.ParseDuration("3m")
	step := datetimeFixStep{duration: duration}

	stamp, _ := time.Parse(time.UnixDate, "Thu Jun 1 17:00:00 CST 2017")

	data := datetime{
		column: column{
			title: "test",
		},
		format: time.UnixDate,
		now:    stamp,
		end:    time.Unix(1<<63-62135596801, 999999999),
		step:   &step,
	}

	result, err := data.Data()
	if err != nil {
		t.Errorf("Something wrong when get datetime data: %v", err)
	} else if result != "Thu Jun  1 17:03:00 CST 2017" {
		t.Errorf("Unexcepted result: %s", result)
	}
}

func TestStepRandom(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())

	durationN, _ := time.ParseDuration("12m")
	durationA, _ := time.ParseDuration("-2m")

	step := datetimeRandomStep{
		durationN: durationN,
		durationA: durationA,
	}

	stamp := time.Now()

	for i := 0; i < 10; i++ {
		data := datetime{
			column: column{
				title: "test",
			},
			format: time.UnixDate,
			now:    stamp,
			end:    time.Unix(1<<63-62135596801, 999999999),
			step:   &step,
		}

		result, err := data.Data()
		if err != nil {
			t.Errorf("Something wrong when get datetime data: %v", err)
			return
		}

		resultStamp, _ := time.Parse(time.UnixDate, result)
		minutes := resultStamp.Sub(stamp) / time.Minute
		if minutes > 10 || minutes < -2 {
			t.Errorf("Unexcepted duration: %v mins", minutes)
			return
		}
	}

}
