package sssly

import (
	"os"
)

func (s *Sssly) Upload(key, fname string) error {
	var fh *os.File
	var err error

	fh, err = os.Open(fname)
	if err != nil {
		Goose.Storage.Logf(1, "Error opening %s for read: %s", fname, err)
		return err
	}

	err = s.UploadFromReader(key, fh)
	if err != nil {
		Goose.Storage.Logf(1, "Error opening %s for write: %s", key, err)
	}

	return err
}

