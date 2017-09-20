package compress

import (
	"compress/gzip"
	"io"
)

// NewGzipReader wraps a io.reader to gziped io.reader
func NewGzipReader(source io.Reader, gzipLevel int) io.Reader {
	r, w := io.Pipe()
	go func() {
		defer w.Close()

		buffer := make([]byte, 1024)
		zip, err := gzip.NewWriterLevel(w, gzipLevel)
		defer zip.Close()
		if err != nil {
			w.CloseWithError(err)
		}

		io.CopyBuffer(zip, source, buffer)
	}()
	return r
}
