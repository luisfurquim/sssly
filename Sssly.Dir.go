package sssly

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (s *Sssly) Dir() ([]string, error) {
	var result *s3.ListObjectsV2Output
	var err error
	var dir []string
	var content types.Object

	result, err = s.Client.ListObjectsV2(
		context.TODO(),
		&s3.ListObjectsV2Input{
			Bucket: aws.String(s.Bucket),
		},
	)

	if err != nil {
		Goose.Storage.Logf(1,"Error listing bucket %s: %s", s.Bucket, err)
		return nil, err
	}

	dir = make([]string, 0, len(result.Contents))
	for _, content = range result.Contents {
		dir = append(dir, *content.Key)
	}

	return dir, nil
}
