package sssly

import (
	"io"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *Sssly) NewReadCloser(key string, chunkSchema ...int32) (io.ReadCloser, error) {
	var err error
	var resp *s3.GetObjectOutput

	if len(chunkSchema)!=0 && len(chunkSchema)!=2 {
		return nil, ErrWrongParmCount
	}

	if len(chunkSchema) == 0 {
		resp, err = s.Client.GetObject(
			context.TODO(),
			&s3.GetObjectInput{
				Bucket: aws.String(s.Bucket),
				Key:    aws.String(s.BasePath + key),
			},
		)
	} else {
		Goose.Storage.Logf(0, "Multipart reader: %s", key)
		return &ReadCloser{
			chunks:     chunkSchema[0],
			chunkSize:  chunkSchema[1],
			cli: 			s,
			key: 			key,
			buffer:     make([]byte, chunkSchema[1]+46),
		}, nil
	}

	if err != nil {
		Goose.Storage.Logf(1, "Error opening %s for read: %s", key, err)
		return nil, err
	}

	return resp.Body, nil
}

