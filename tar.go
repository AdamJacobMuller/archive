package archive

import (
	"archive/tar"
	"fmt"
	"os"
	"time"
)

type TarIterator struct {
	reader *tar.Reader
}

type TarFile struct {
	name string
	body []byte
	size int64
	fi   os.FileInfo
}

func (f *TarFile) Size() int64 {
	return f.size
}

func (f *TarFile) Mode() os.FileMode {
	return f.fi.Mode()
}

func (f *TarFile) Name() string {
	return f.name
}

func (f *TarFile) ModTime() time.Time {
	return f.fi.ModTime()
}

func (f *TarFile) IsDir() bool {
	return f.fi.IsDir()
}

func (f *TarFile) Sys() interface{} {
	return f
}

func (f *TarFile) Bytes() []byte {
	return f.body
}

func (t TarIterator) Next() (File, error) {
	header, err := t.reader.Next()
	if err != nil {
		return nil, err
	}
	tf := &TarFile{}
	tf.name = header.Name
	tf.fi = header.FileInfo().(os.FileInfo)
	tf.size = header.Size
	if header.Size > 0 {
		tf.body = make([]byte, header.Size)
		_, err = t.reader.Read(tf.body)
		if err != nil {
			panic(fmt.Sprintf("TarIterator.Next(%s)", err))
		}
	}
	return tf, nil
}
