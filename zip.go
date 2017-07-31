package archive

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type ZipIterator struct {
	reader *zip.Reader
	pos    int
	len    int
}

type ZipFile struct {
	name string
	size int64
	fi   os.FileInfo
	body []byte
}

func (f *ZipFile) Size() int64 {
	return f.size
}

func (f *ZipFile) Mode() os.FileMode {
	return f.fi.Mode()
}

func (f *ZipFile) Name() string {
	return f.name
}

func (f *ZipFile) ModTime() time.Time {
	return f.fi.ModTime()
}

func (f *ZipFile) IsDir() bool {
	return f.fi.IsDir()
}

func (f *ZipFile) Sys() interface{} {
	return f
}

func (f *ZipFile) Bytes() []byte {
	return f.body
}

func NewZipFile(r io.Reader) (*ZipIterator, error) {
	var err error
	z := &ZipIterator{}

	zBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("NewZipFile ioutil.ReadAll() == %s", err))
	}

	zReader := bytes.NewReader(zBytes)

	z.reader, err = zip.NewReader(zReader, int64(len(zBytes)))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("NewZipFile zip.NewReader() == %s", err))
	}
	z.len = len(z.reader.File)
	z.pos = 0
	return z, nil
}

func (t *ZipIterator) Next() (File, error) {
	if t.pos >= t.len {
		return nil, io.EOF
	}
	var err error
	var rc io.ReadCloser
	nZipFile := t.reader.File[t.pos]
	t.pos += 1

	zipFile := &ZipFile{}
	zipFile.name = nZipFile.Name
	zipFile.fi = nZipFile.FileInfo()
	rc, err = nZipFile.Open()
	if err != nil {
		return nil, err
	}

	zipFile.body, err = ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	return zipFile, nil
}
