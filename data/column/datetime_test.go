package column

import (
	"testing"
	"time"
)

func TestCreateDatetimeFactory(t *testing.T) {
	tmpl := map[string]interface{}{}

	format := "2006-01-02 15:04:05"

	tmpl["Format"] = format
	step := map[string]interface{}{}
	step["Type"] = "Random"
	step["Unit"] = "us"
	step["Max"] = float64(10000)
	step["Min"] = float64(100)
	step["Start"] = "2015-01-01 00:00:00"
	tmpl["Step"] = step
	fileStep := map[string]interface{}{}
	fileStep["Duration"] = "1h"
	tmpl["FileStep"] = fileStep

	factory, err := newDatetimeFactory(columnChangePerRowAndFile, tmpl)
	if err != nil {
		t.Errorf("Someting wrong when new datetime factory: %v", err)
	}
	data := factory.Create()
	_, err = data.Data()
	if err != nil {
		t.Errorf("Something wrong when get datetime data: %v", err)
	}
}

func TestDatetimeRandom(t *testing.T) {
	stamp := time.Now()

	_, err := newDatetimeRandom(time.UnixDate, "m", 1, 1, stamp, maxTime, "", columnChangePerRowAndFile)
	if err == nil {
		t.Errorf("Excepted an error, but nil")
	}

	for i := 0; i < 10; i++ {
		data1, err := newDatetimeRandom(time.UnixDate, "m", 10, -2, stamp, maxTime, "", columnChangePerRowAndFile)
		if err != nil {
			t.Errorf("Someting wrong when new datetime: %v", err)
		}
		data1A := data1.Create()
		result1, err := data1A.Data()
		if err != nil {
			t.Errorf("Something wrong when get datetime data: %v", err)
		} else {
			resultStamp1, _ := time.Parse(time.UnixDate, result1)
			minutes := resultStamp1.Sub(stamp) / time.Minute
			if minutes > 10 || minutes < -2 {
				t.Errorf("Unexcepted duration: %v mins", minutes)
				return
			}
		}
	}

	for i := 0; i < 10; i++ {
		data2, err := newDatetimeRandom(time.UnixDate, "m", 10, -2, stamp, stamp.Add(-5*time.Minute), "", columnChangePerRowAndFile)
		if err != nil {
			t.Errorf("Someting wrong when new datetime: %v", err)
		}
		data2A := data2.Create()
		result, err := data2A.Data()
		if err == nil {
			t.Errorf("Excepted an error, but nil: %v vs %v", stamp.Add(-5*time.Minute).Format(time.UnixDate), result)
		}
	}
}

func TestDatetimeFix(t *testing.T) {
	stamp, _ := time.Parse(time.UnixDate, "Thu Jun 1 17:00:00 CST 2017")
	data1A, err := newDatetimeIncrease(time.UnixDate, "3m", stamp, maxTime, "", columnChangePerRowAndFile)
	if err != nil {
		t.Errorf("Someting wrong when new datetime: %v", err)
	}
	data1 := data1A.Create()
	result1, err := data1.Data()
	if err != nil {
		t.Errorf("Something wrong when get datetime data: %v", err)
	} else if result1 != "Thu Jun  1 17:00:00 CST 2017" {
		t.Errorf("Unexcepted result: %s", result1)
	}

	stampEnd, _ := time.Parse(time.UnixDate, "Thu Jun 1 17:00:00 CST 2017")
	data2A, err := newDatetimeIncrease(time.UnixDate, "3m", stamp, stampEnd, "", columnChangePerRowAndFile)
	if err != nil {
		t.Errorf("Someting wrong when new datetime: %v", err)
	}
	data2 := data2A.Create()
	_, err = data2.Data()
	_, err = data2.Data()
	if err == nil {
		t.Error("Excepted an error, but nil")
	}
}

func TestFileStep(t *testing.T) {
	stamp, _ := time.Parse(time.UnixDate, "Thu Jun 1 17:00:00 CST 2017")
	data1, err := newDatetimeIncrease(time.UnixDate, "30m", stamp, maxTime, "1h", columnChangePerRowAndFile)
	if err != nil {
		t.Errorf("Someting wrong when new datetime: %v", err)
	}
	data1A := data1.Create()
	result1, err := data1A.Data()
	if err != nil {
		t.Errorf("Something wrong when get datetime data: %v", err)
	} else if result1 != "Thu Jun  1 17:00:00 CST 2017" {
		t.Errorf("Unexcepted result: %s", result1)
	}
	result2, err := data1A.Data()
	if err != nil {
		t.Errorf("Something wrong when get datetime data: %v", err)
	} else if result2 != "Thu Jun  1 17:30:00 CST 2017" {
		t.Errorf("Unexcepted result: %s", result2)
	}
	result3, err := data1A.Data()
	if err != nil {
		t.Errorf("Something wrong when get datetime data: %v", err)
	} else if result3 != "Thu Jun  1 18:00:00 CST 2017" {
		t.Errorf("Unexcepted result: %s", result3)
	}
	_, err = data1A.Data()
	if err == nil {
		t.Error("Excepted an error, but nil")
	}
	if !IsDataOver(err) {
		t.Error("Not a time over error")
	}

	data1B := data1.Create()
	result4, err := data1B.Data()
	if err != nil {
		t.Errorf("Something wrong when get datetime data: %v", err)
	} else if result4 != "Thu Jun  1 18:00:00 CST 2017" {
		t.Errorf("Unexcepted result: %s", result4)
	}
	result5, err := data1B.Data()
	if err != nil {
		t.Errorf("Something wrong when get datetime data: %v", err)
	} else if result5 != "Thu Jun  1 18:30:00 CST 2017" {
		t.Errorf("Unexcepted result: %s", result5)
	}
}
