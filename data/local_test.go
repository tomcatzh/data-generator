package data

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestLocal(t *testing.T) {
	path := "/tmp/"

	const content = "abcdefghijklmnopqrstuvwxyz\n"
	r := bytes.NewReader([]byte(content))

	l := newStorageLocal(path)

	name := "test1"
	length, err := l.Save(name, r)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	if length != 27 {
		t.Errorf("Unexcepted write length: %v", length)
	}

	fullname := path + name
	fileContent, err := ioutil.ReadFile(fullname)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	if string(fileContent) != content {
		t.Errorf("Unexcepted content: %v", fileContent)
	}
	err = os.Remove(fullname)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
}
