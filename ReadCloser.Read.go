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
		sz, n int
	)

	Goose.Storage.Logf(4, "Going to read: %d, consumed: %d", rc.chunk, rc.consumed)

	if rc.remReader == nil {
		Goose.Storage.Logf(0, "Obtaining remote reader")

		resp, err = rc.cli.Client.GetObject(
			context.TODO(),
			&s3.GetObjectInput{
				Bucket: aws.String(rc.cli.Bucket),
				Key:    aws.String(rc.cli.BasePath + rc.key),
//				PartNumber: &rc.chunk,
			},
		)

		if err != nil {
			Goose.Storage.Logf(1, "Error fetching %s for next chunk[%d]: %s", rc.key, rc.chunk, err)
			return 0, err
		}

		rc.remReader = res.Body
		rc.consumed  = rc.chunkSize
	}

	if rc.consumed >= rc.chunkSize {
		rc.chunk++
		if rc.chunk > rc.chunks {
			return 0, io.EOF
		}

		Goose.Storage.Logf(0, "Fetching new chunk: %d", rc.chunk)

		for sz<len(rc.buffer) && err==nil {
			Goose.Storage.Logf(4, "sz: %d", sz)
			n, err = rc.remReader.Read(rc.buffer[sz:])
			if err != nil && err != io.EOF {
				Goose.Storage.Logf(1, "Error reading %s on chunk[%d]: %s", rc.key, rc.chunk, err)
				return 0, err
			}
			sz += n
		}

		Goose.Storage.Logf(0, "sz=%d", sz)
		rc.consumed = 0
		sz -= 38
		rc.rd = bytes.NewReader(rc.buffer[8:sz])
//		rc.rd = bytes.NewReader(rc.buffer[8:rc.chunkSize+8])
		Goose.Storage.Logf(0, "Removing header and trailer: %d % 2x .. % 2x .. % 2x", rc.chunk, rc.buffer[:8], rc.buffer[8:16], rc.buffer[sz:sz+16])
		Goose.Storage.Logf(0, "Removing header and trailer: %d %s .. %s .. %s", rc.chunk, rc.buffer[:8], rc.buffer[8:16], rc.buffer[sz:sz+16])
	}

	n, err = rc.rd.Read(buf)
	Goose.Storage.Logf(4, "Read %d bytes: %s", n, err)
	if err == io.EOF {
		if rc.chunk < rc.chunks {
			err = nil
			if n == 0 {
				defer panic("Stopped!")
			}
		}
	} else if err != nil {
		Goose.Storage.Logf(1, "Error reading %s on chunk[%d]: %s", rc.key, rc.chunk, err)
		return 0, err
	}

	rc.consumed += int32(n)

	return n, nil
}

