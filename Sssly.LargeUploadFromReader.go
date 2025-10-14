package sssly

import (
	"io"
	"time"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)


func (s *Sssly) LargeUploadFromReader(key string, rd io.Reader) error {
	var (
		err error
		uploader *manager.Uploader
		fullKey string
	)

	fullKey = s.BasePath + key

	uploader = manager.NewUploader(s.Client, func(u *manager.Uploader) {
		u.PartSize = int64(s.MaxChunk)
	})

	_, err = uploader.Upload(s.ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(fullKey),
		Body:   rd,
	})
	if err != nil {
		Goose.Storage.Logf(1, "Error uploading to %s: %s", s.Bucket, err)
	} else {
		err = s3.NewObjectExistsWaiter(s.Client).Wait(
			s.ctx,
			&s3.HeadObjectInput{Bucket: aws.String(s.Bucket), Key: aws.String(fullKey)},
			time.Minute,
		)
		if err != nil {
			Goose.Storage.Logf(1, "Failed attempt to wait for object %s to exist: %s", fullKey, err)
		}
	}

	return err
}
