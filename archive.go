package archive

import (
	archive_tar "archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

func reader(rulesPath string) (io.Reader, error) {
	var err error
	var body io.Reader
	if strings.HasPrefix(rulesPath, "https://") ||
		strings.HasPrefix(rulesPath, "http://") {
		resp, err := http.Get(rulesPath)
		return resp.Body, err
	} else {
		body, err = os.Open(rulesPath)
		return body, err
	}
}

type File interface {
	Name() string
	Bytes() []byte
}

type FileIterator interface {
	Next() (*File, error)
}

func decoder(name string, f io.Reader) (FileIterator, error) {
	if strings.HasSuffix(name, ".tar.gz") {
		gzf, err := gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
		tarReader := archive_tar.NewReader(gzf)
		return TarIterator{tarReader: tarReader}, nil
	}
	return nil, errors.New("unsupported file")
}
