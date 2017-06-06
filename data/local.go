package data

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

type storageLocal struct {
	path string
}

func (l *storageLocal) Save(key string, reader io.ReadSeeker) (int64, error) {
	var path bytes.Buffer
	path.WriteString(l.path)

	lastSeparator := strings.LastIndexByte(key, os.PathSeparator)
	if lastSeparator != -1 {
		path.WriteByte(os.PathSeparator)
		path.WriteString(key[0:lastSeparator])
	}

	err := os.MkdirAll(path.String(), 0770)
	if err != nil {
		return 0, err
	}

	path.Reset()
	path.WriteString(l.path)
	path.WriteByte(os.PathSeparator)
	path.WriteString(key)

	f, err := os.Create(path.String())
	defer f.Close()
	if err != nil {
		return 0, err
	}
	w := bufio.NewWriter(f)
	defer w.Flush()
	len, err := io.Copy(w, reader)
	return len, err
}

func newStorageLocal(path string) *storageLocal {
	if path[len(path)-1] == os.PathSeparator {
		path = path[0 : len(path)-1]
	}

	return &storageLocal{
		path: path,
	}
}
