package sssly

import (
	"io"
	"os"
)

func (s *Sssly) Download(fname, key string) (int64, error) {
	var err error
	var rd io.ReadCloser
	var fh *os.File

	rd, err = s.NewReadCloser(key)
	if err != nil {
		return 0, err
	}
	defer rd.Close()

	fh, err = os.Create(fname)
	if err != nil {
		Goose.Storage.Logf(1, "Error creating file %s: %s", fname, err)
		return 0, err
	}
	defer fh.Close()

	return io.Copy(fh, rd)
}
