package archive

import (
	archive_tar "archive/tar"
)

type TarIterator struct {
	tarReader *archive_tar.Reader
}

func (t TarIterator) Next() (*File, error) {
	return nil, nil
}
