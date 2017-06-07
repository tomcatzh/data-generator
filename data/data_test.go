package data

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func TestRow(t *testing.T) {
	var r row

	title := []string{"Date", "User", "Member"}

	stamp, _ := time.Parse(time.UnixDate, "Thu Jun 1 17:00:00 CST 2017")
	c1, _ := newDatetimeIncrease(title[0], time.UnixDate, "1s", stamp, maxTime)
	r.AddColumn(c1)

	users := []string{"Alice", "Bob", "Charles"}
	c2 := newEnumString(title[1], users)
	r.AddColumn(c2)

	member := "student"
	c3 := newFixString(title[2], member)
	r.AddColumn(c3)

	titleResult := r.Title()
	if titleResult[0] != title[0] || titleResult[1] != title[1] || titleResult[2] != title[2] {
		t.Errorf("Unexcepted title data: %v", titleResult)
	}

	dataResult, err := r.Data()
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	if dataResult[0] != "Thu Jun  1 17:00:01 CST 2017" {
		t.Errorf("Unexcepted Date: %v", dataResult[0])
	}
	if !stringContins(users, dataResult[1]) {
		t.Errorf("Unexcepted User: %v", dataResult[1])
	}
	if dataResult[2] != "student" {
		t.Errorf("Unexcepted Member: %v", dataResult[2])
	}
}

func TestNewTemplate(t *testing.T) {
	template, err := NewTemplate("../templates/sample.json")
	if err != nil {
		t.Errorf("Unexcepted error on template: %v", err)
		return
	}

	f, err := template.getFile()
	if err != nil {
		t.Errorf("Unexcepted error on file: %v", err)
	}

	r, err := f.Data()
	if err != nil {
		t.Errorf("Unexcepted error on get reader: %v", err)
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("Unexcepted error on reading data: %v", err)
	}
	fmt.Println(len(b))

	f, err = template.getFile()
	if err != nil {
		t.Errorf("Unexcepted error on file: %v", err)
	}

	r, err = f.Data()
	if err != nil {
		t.Errorf("Unexcepted error on get reader: %v", err)
	}
	b, err = ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("Unexcepted error on reading data: %v", err)
	}
	fmt.Println(len(b))
}

func TestTemplateIterate(t *testing.T) {
	template, err := NewTemplate("../templates/sample.json")
	if err != nil {
		t.Errorf("Unexcepted error on template: %v", err)
	}

	for f := range template.Iterate() {
		n, err := f.Name()
		if err != nil {
			t.Errorf("Unexcepted error on get reader: %v", err)
		}
		fmt.Println(n)

		r, err := f.Data()
		if err != nil {
			t.Errorf("Unexcepted error on get reader: %v", err)
		}
		b, err := ioutil.ReadAll(r)
		if err != nil {
			t.Errorf("Unexcepted error on reading data: %v", err)
		}
		fmt.Println(len(b))
	}
}

func TestSaveS3(t *testing.T) {
	template, err := NewTemplate("../templates/s3sample.json")
	if err != nil {
		t.Errorf("Unexcepted error on template: %v", err)
		return
	}

	for f := range template.Iterate() {
		f.Save()
	}
}
