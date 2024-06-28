package sssly

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (s *Sssly) Delete(keys ...string) error {
	var key string
	var err error
	var objectIds []types.ObjectIdentifier

	for _, key = range keys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(s.BasePath + key)})
	}

	_, err = s.Client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(s.Bucket),
		Delete: &types.Delete{Objects: objectIds},
	})

	if err != nil {
		Goose.Storage.Logf(1, "Error deleting object: %s", err)
	}

	return err
}
