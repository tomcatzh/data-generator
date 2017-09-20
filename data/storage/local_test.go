package storage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestLocal(t *testing.T) {
	path := "/tmp"

	const content = "abcdefghijklmnopqrstuvwxyz\n"
	r := bytes.NewReader([]byte(content))

	template := map[string]interface{}{}
	template["Path"] = path
	template["BufferSizeM"] = (float64)(1)

	l := newStorageLocal(template)

	name := "test/test/test1"
	length, err := l.Save(name, r)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	if length != 27 {
		t.Errorf("Unexcepted write length: %v", length)
	}

	fullname := fmt.Sprintf("%v%c%v", path, os.PathSeparator, name)
	fileContent, err := ioutil.ReadFile(fullname)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	if string(fileContent) != content {
		t.Errorf("Unexcepted content: %v", fileContent)
	}
	err = os.RemoveAll("/tmp/test")
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
}
