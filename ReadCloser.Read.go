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

	Goose.Storage.Logf(0, "Going to read: %d", rc.chunk)

	if rc.rd == nil || rc.consumed>=rc.chunkSize {
		rc.chunk++
		if rc.chunk > rc.chunks {
			return 0, io.EOF
		}
		
		Goose.Storage.Logf(0, "Fetching new chunk: %d", rc.chunk)

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

		Goose.Storage.Logf(0, "Fetching new chunk: %d", rc.chunk)
		n, err = resp.Body.Read(rc.buffer)
		if err != nil && err != io.EOF {
			Goose.Storage.Logf(1, "Error reading %s on chunk[%d]: %s", rc.key, rc.chunk, err)
			return 0, err
		}

		Goose.Storage.Logf(0, "Removing header and trailer: % 2x .. % 2x .. % 2x", rc.chunk, rc.buffer[:8], rc.buffer[8:16], rc.buffer[rc.chunkSize:rc.chunkSize+16])
		rc.consumed = 0
		rc.rd = io.NopCloser(bytes.NewReader(rc.buffer[8:rc.chunkSize]))
	}

	n, err = rc.rd.Read(buf)
	Goose.Storage.Logf(0, "Read %d bytes: %s", n, err)
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

