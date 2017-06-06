package data

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCsv(t *testing.T) {
	namePart := []namePart{
		namePart{
			partType: namePartTypeData,
			index:    0,
		},
		namePart{
			partType:  namePartTypeSubData,
			index:     2,
			substring: 2,
		},
		namePart{
			partType: namePartTypeFix,
			value:    "abc",
		},
	}
	c := newCsv(2, namePart, defaultDelimiter, defaultQuoteChar, defaultEscapeChar, defaultLineTeriminator, true)
	c.AddColumn(newFixString("test", "a"))
	c.AddColumn(newFixString("test", "b"))
	c.AddColumn(newFixString("test", "cccc"))

	const content = "\"test\",\"test\",\"test\"\r\n\"a\",\"b\",\"cccc\"\r\n\"a\",\"b\",\"cccc\"\r\n"

	r, err := c.Data()
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	s := string(b)
	if s != content {
		t.Errorf("Unexcepted result: %v", s)
	}

	tmpPath := "/tmp/"

	n, err := c.Name()
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	if n != "accabc" {
		t.Errorf("Unexcepted name: %v", s)
	}

	c.SetStorage(newStorageLocal(tmpPath))

	len, err := c.Save()
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	if len != 54 {
		t.Errorf("Unexcepted write length: %v", len)
	}

	tmpFile := tmpPath + n
	tmpContent, err := ioutil.ReadFile(tmpFile)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	if string(tmpContent) != content {
		t.Errorf("Unexcepted result: %v", s)
	}
	err = os.Remove(tmpFile)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
}

func TestCsv2(t *testing.T) {
	namePart := []namePart{}
	c := newCsv(2, namePart, defaultDelimiter, defaultQuoteChar, "\\", defaultLineTeriminator, false)
	c.AddColumn(newFixString("test", "a\""))
	c.AddColumn(newFixString("test", "b"))
	c.AddColumn(newFixString("test", "c"))

	r, err := c.Data()
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	s := string(b)
	if s != "\"a\\\"\",\"b\",\"c\"\r\n\"a\\\"\",\"b\",\"c\"\r\n" {
		t.Errorf("Unexcepted result: %v", s)
	}
}

func TestSmallRead(t *testing.T) {
	namePart := []namePart{}
	c := newCsv(2, namePart, defaultDelimiter, defaultQuoteChar, "\\", defaultLineTeriminator, false)
	c.AddColumn(newFixString("test", "a"))
	c.AddColumn(newFixString("test", "b"))
	c.AddColumn(newFixString("test", "c"))

	r, err := c.Data()
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}

	buf := make([]byte, 10)

	n, err := r.Read(buf)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	if n != 10 {
		t.Errorf("Unexcepted bytes: %v", n)
	}

	n, err = r.Read(buf)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	if n != 10 {
		t.Errorf("Unexcepted bytes: %v", n)
	}

	n, err = r.Read(buf)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	if n != 6 {
		t.Errorf("Unexcepted bytes: %v", n)
	}
}
