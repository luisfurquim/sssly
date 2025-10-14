package sssly

import (
	"io"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *Sssly) NewLargeReadCloser(key string) (io.ReadCloser, error) {
	var (
		err error
		resp *s3.GetObjectOutput
		downloader *manager.Downloader
		lrc LargeReadCloser
	)

	lrc.cli = s
	lrc.key = key

	downloader = manager.NewDownloader(s.Client, func(d *manager.Downloader) {
		d.PartSize = s.MaxChunk
	})

	go func() {
		var err error

		lrc.ahead = map[int64][]byte{}

		_, err = downloader.Download(
			s.ctx,
			&lrc,
			&s3.GetObjectInput{
				Bucket: aws.String(s.Bucket),
				Key:    aws.String(s.BasePath + key),
			},
		)
		if err != nil {
			Goose.Storage.Logf(1, "Couldn't download large object from %s/%s:%s", s.Bucket, s.BasePath + key, err)
		}

		lrc.mtx.Lock()
		lrc.eof = true
		lrc.mtx.Unlock()

	}

	if err != nil {
		Goose.Storage.Logf(1, "Error opening %s for read: %s", key, err)
		return nil, err
	}

	return &lrc, nil
}
