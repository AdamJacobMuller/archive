package archive

import (
	archive_tar "archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func Open(filename string) (FileIterator, error) {
	reader, err := reader(filename)
	if err != nil {
		return nil, err
	}
	decoder, err := decoder(filename, reader)
	if err != nil {
		return nil, err
	}
	return decoder, nil
}

func reader(filename string) (io.Reader, error) {
	var err error
	var body io.Reader
	if strings.HasPrefix(filename, "https://") ||
		strings.HasPrefix(filename, "http://") {
		resp, err := http.Get(filename)
		return resp.Body, err
	} else {
		body, err = os.Open(filename)
		return body, err
	}
}

type File interface {
	Name() string
	Size() int64
	Mode() os.FileMode
	ModTime() time.Time
	IsDir() bool
	Sys() interface{}

	Bytes() []byte
}

type FileIterator interface {
	Next() (File, error)
}

func decoder(name string, f io.Reader) (FileIterator, error) {
	if strings.HasSuffix(name, ".tar.gz") ||
		strings.HasSuffix(name, ".tgz") {
		gzf, err := gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
		tarReader := archive_tar.NewReader(gzf)
		return TarIterator{reader: tarReader}, nil
	} else if strings.HasSuffix(name, ".tar") {
		tarReader := archive_tar.NewReader(f)
		return TarIterator{reader: tarReader}, nil
	} else if strings.HasSuffix(name, ".zip") {
		return NewZipFile(f)
	}
	return nil, errors.New("unsupported file")
}
