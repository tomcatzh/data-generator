package storage

import (
	"bytes"
	"io"
	"os"
	"strings"
)

const defaultBufferSize = 5

type storageLocal struct {
	path        string
	bufferSizeM int
}

func (l *storageLocal) Save(key string, reader io.Reader) (int64, error) {
	var path bytes.Buffer
	path.WriteString(l.path)

	if l.path[len(l.path)-1] != os.PathSeparator && key[0] != os.PathSeparator {
		path.WriteByte(os.PathSeparator)
	}

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

	buffer := make([]byte, 1024*1024*l.bufferSizeM)

	len, err := io.CopyBuffer(f, reader, buffer)
	return len, err
}

func newStorageLocal(s map[string]interface{}) *storageLocal {
	path, ok := s["Path"].(string)
	if !ok || path == "" {
		path = "."
	}
	sbufferSize, ok := s["BufferSizeM"].(float64)
	var bufferSizeM int
	if !ok || sbufferSize <= 0 {
		bufferSizeM = defaultBufferSize
	} else {
		bufferSizeM = int(sbufferSize)
	}

	if path[len(path)-1] == os.PathSeparator {
		path = path[0 : len(path)-1]
	}

	return &storageLocal{
		path:        path,
		bufferSizeM: bufferSizeM,
	}
}
