package sssly

import (
	"os"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *Sssly) Upload(key, fname string) error {
	var fh *os.File
	var err error

	fh, err = os.Open(fname)
	if err != nil {
		Goose.Storage.Logf(1, "Error opening %s for read: %s", fname, err)
		return err
	}

	_, err = s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s.BasePath + key),
		Body:   fh,
	})

	if err != nil {
		Goose.Storage.Logf(1, "Error opening %s for write: %s", key, err)
	}

	return err
}

