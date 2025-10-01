package sssly

import (
	"io"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *Sssly) NewReadCloser(key string, chunked ...bool) (io.ReadCloser, error) {
	var err error
	var resp *s3.GetObjectOutput

	if len(chunked)!=0 && len(chunked)!=1 {
		return nil, ErrWrongParmCount
	}

	if len(chunked) == 0 {
		Goose.Storage.Logf(4, "Monolithic reader: %s", key)
		resp, err = s.Client.GetObject(
			context.TODO(),
			&s3.GetObjectInput{
				Bucket: aws.String(s.Bucket),
				Key:    aws.String(s.BasePath + key),
			},
		)
	} else {
		Goose.Storage.Logf(4, "Multipart reader: %s", key)
		return &ReadCloser{
			cli: 			s,
			key: 			key,
		}, nil
	}

	if err != nil {
		Goose.Storage.Logf(1, "Error opening %s for read: %s", key, err)
		return nil, err
	}

	return resp.Body, nil
}

