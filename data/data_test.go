package data

import (
	"fmt"
	"io/ioutil"
	"testing"
)

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
