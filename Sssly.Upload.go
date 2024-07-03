package sssly

import (
	"os"
)

func (s *Sssly) Upload(key, fname string) error {
	var fh *os.File
	var err error
	var fi os.FileInfo
	var sz int64
	

	fh, err = os.Open(fname)
	if err != nil {
		Goose.Storage.Logf(1, "Error opening %s for read: %s", fname, err)
		return err
	}

	fi, err = fh.Stat()
	if err != nil {
		Goose.Collect.Logf(1,"Error checking zip file stat: %s", err)
		return err
	}
	sz = fi.Size()

	err = s.UploadFromReader(key, fh, &sz)
	if err != nil {
		Goose.Storage.Logf(1, "Error opening %s for write: %s", key, err)
	}

	return err
}

