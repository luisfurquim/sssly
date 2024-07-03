package sssly

import (
	"io"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *Sssly) UploadFromReader(key string, rd io.Reader, sz int64) error {
	var err error

	_, err = s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s.BasePath + key),
		ContentLength: &sz,
		Body:   rd,
	})

	if err != nil {
		Goose.Storage.Logf(1, "Error opening %s for write: %s", key, err)
	}

	return err
}

