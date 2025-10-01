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
		chunkSize int
		hdrBuffer []byte
	)

	Goose.Storage.Logf(4, "Going to read: %d, consumed: %d", rc.chunk, rc.consumed)

	if rc.remReader == nil {
		Goose.Storage.Logf(3, "Obtaining remote reader")

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

		rc.remReader = resp.Body
		rc.consumed  = 1
	}

	if int(rc.consumed+38) >= len(rc.buffer) {
		if rc.done {
			return 0, io.EOF
		}

		rc.chunk++

		Goose.Storage.Logf(3, "Fetching new chunk: %d", rc.chunk)

		hdrBuffer = make([]byte, 1)

		for {
			n, err = rc.remReader.Read(hdrBuffer)
			if err != nil && err != io.EOF {
				Goose.Storage.Logf(1, "Error reading %s chunkSize on chunk[%d]: %s", rc.key, rc.chunk, err)
				return 0, err
			}
			if n==0 {
				Goose.Storage.Logf(1, "Error reading %s chunkSize on chunk[%d]: no bytes available", rc.key, rc.chunk)
				return 0, NoBytesAvailable
			}
			if hdrBuffer[0] == '\r' {
				n, err = rc.remReader.Read(hdrBuffer)
				if err != nil && err != io.EOF {
					Goose.Storage.Logf(1, "Error reading %s linefeed on chunk[%d]: %s", rc.key, rc.chunk, err)
					return 0, err
				}
				if n==0 {
					Goose.Storage.Logf(1, "Error reading %s linefeed on chunk[%d]: %s", rc.key, rc.chunk, NoBytesAvailable)
					return 0, NoBytesAvailable
				}
				if hdrBuffer[0] != '\n' {
					Goose.Storage.Logf(1, "Error reading %s linefeed on chunk[%d]: %s", rc.key, rc.chunk, UnexpectedCharacter)
					return 0, UnexpectedCharacter
				}
				break
			} else if hdrBuffer[0] >= 'a' && hdrBuffer[0] <= 'f' {
				chunkSize <<= 4
				chunkSize += int(hdrBuffer[0] - 'a' + 10)
			} else if hdrBuffer[0] >= 'A' && hdrBuffer[0] <= 'F' {
				chunkSize <<= 4
				chunkSize += int(hdrBuffer[0] - 'A' + 10)
			} else if hdrBuffer[0] >= '0' && hdrBuffer[0] <= '9' {
				chunkSize <<= 4
				chunkSize += int(hdrBuffer[0] - '0')
			} else {
				Goose.Storage.Logf(1, "Error reading %s chunkSize on chunk[%d]: %s", rc.key, rc.chunk, UnexpectedCharacter)
				return 0, UnexpectedCharacter
			}
		}

		chunkSize += 38
		rc.buffer = make([]byte, chunkSize)

		for sz<len(rc.buffer) && sz<chunkSize && err==nil {
			Goose.Storage.Logf(4, "sz: %d", sz)
			n, err = rc.remReader.Read(rc.buffer[sz:])
			if err != nil && err != io.EOF {
				Goose.Storage.Logf(1, "Error reading %s on chunk[%d]: %s", rc.key, rc.chunk, err)
				return 0, err
			}
			sz += n
			if err == io.EOF {
				rc.done = true
			}
		}

		Goose.Storage.Logf(4, "sz=%d", sz)
		rc.consumed = 0
		sz -= 38
		rc.rd = bytes.NewReader(rc.buffer[:sz])
//		rc.rd = bytes.NewReader(rc.buffer[8:rc.chunkSize+8])
		Goose.Storage.Logf(3, "Removing trailer: %d % 2x .. % 2x", rc.chunk,  rc.buffer[:8], rc.buffer[sz:])
		Goose.Storage.Logf(3, "Removing trailer: %d %s .. %s", rc.chunk, rc.buffer[:8], rc.buffer[sz:])
	}

	n, err = rc.rd.Read(buf)
	rc.consumed += int32(n)
	Goose.Storage.Logf(4, "Read %d bytes: %s", n, err)
	if err == io.EOF {
		Goose.Storage.Logf(3,"EOF")
		if !rc.done {
			Goose.Storage.Logf(3,"!rc.done")
			if n == 0 {
				Goose.Storage.Logf(1,"Incomplete reading!")
				return 0, err
			}
			err = nil
		}
	} else if err != nil {
		Goose.Storage.Logf(1, "Error reading %s on chunk[%d]: %s", rc.key, rc.chunk, err)
		return 0, err
	}

	Goose.Storage.Logf(4,"done reading %d bytes", n)

	return n, err
}
