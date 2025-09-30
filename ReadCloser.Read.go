package sssly

import (
	"io"
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (rc *ReadCloser) Read(buf []byte) (int, error) {
	var (
		err error
		resp *s3.GetObjectOutput
		n int
	)

	if rc.rd == nil || rc.consumed>=rc.chunkSize {
		rc.chunk++
		if rc.chunk > rc.chunks {
			return 0, io.EOF
		}

		if rc.rd != nil {
			err = rc.rd.Close()
			if err != nil {
				Goose.Storage.Logf(1, "Error closing %s for last chunk[%d]: %s", rc.key, rc.chunk-1, err)
				return 0, err
			}
		}

		resp, err = rc.cli.Client.GetObject(
			context.TODO(),
			&s3.GetObjectInput{
				Bucket: aws.String(rc.cli.Bucket),
				Key:    aws.String(rc.cli.BasePath + rc.key),
				PartNumber: &rc.chunk,
			},
		)

		if err != nil {
			Goose.Storage.Logf(1, "Error fetching %s for next chunk[%d]: %s", rc.key, rc.chunk, err)
			return 0, err
		}
		
		defer resp.Body.Close()

		n, err = resp.Body.Read(rc.buffer)
		if err != nil {
			Goose.Storage.Logf(1, "Error reading %s on chunk[%d]: %s", rc.key, rc.chunk, err)
			return 0, err
		}

		rc.consumed = 0
		rc.rd = io.NopCloser(bytes.NewReader(rc.buffer[8:rc.chunkSize]))
	}

	n, err = rc.rd.Read(buf)
	if err == io.EOF {
		if rc.chunk < rc.chunks {
			err = nil
		}
	} else if err != nil {
		Goose.Storage.Logf(1, "Error reading %s on chunk[%d]: %s", rc.key, rc.chunk, err)
		return 0, err
	}

	rc.consumed += int32(n)

	return n, nil
}

