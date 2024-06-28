package sssly

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (wr *WriteCloser) Close() error {
	var err error

	_, err = wr.cli.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(wr.cli.Bucket),
		Key:    aws.String(wr.cli.BasePath + wr.key),
		Body:   wr,
	})

	if err != nil {
		Goose.Storage.Logf(1, "Error opening %s for write: %s", wr.key, err)
	}

	return err
}

