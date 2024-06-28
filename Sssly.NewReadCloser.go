package sssly

import (
	"io"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *Sssly) NewReadCloser(key string) (io.ReadCloser, error) {
	var err error
	var resp *s3.GetObjectOutput

	resp, err = s.Client.GetObject(
		context.TODO(),
		&s3.GetObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(s.BasePath + key),
		},
	)
	if err != nil {
		Goose.Storage.Logf(1, "Error opening %s for read: %s", key, err)
		return nil, err
	}

	return resp.Body, nil
}

