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
		sz, n, i int
		chunkSize int
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
		rc.consumed  = rc.chunkSize
	}

	if rc.consumed >= rc.chunkSize {
		rc.chunk++
		if rc.chunk > rc.chunks {
			return 0, io.EOF
		}

		Goose.Storage.Logf(3, "Fetching new chunk: %d", rc.chunk)

		for {
			n, err = rc.remReader.Read(rc.buffer[i:i+1])
			if err != nil && err != io.EOF {
				Goose.Storage.Logf(1, "Error reading %s chunkSize on chunk[%d]: %s", rc.key, rc.chunk, err)
				return 0, err
			}
			if n==0 {
				Goose.Storage.Logf(1, "Error reading %s chunkSize on chunk[%d]: no bytes available", rc.key, rc.chunk)
				return 0, NoBytesAvailable
			}
			if rc.buffer[i] == '\r' {
				n, err = rc.remReader.Read(rc.buffer[i:i+1])
				if err != nil && err != io.EOF {
					Goose.Storage.Logf(1, "Error reading %s linefeed on chunk[%d]: %s", rc.key, rc.chunk, err)
					return 0, err
				}
				if n==0 {
					Goose.Storage.Logf(1, "Error reading %s linefeed on chunk[%d]: %s", rc.key, rc.chunk, NoBytesAvailable)
					return 0, NoBytesAvailable
				}
				if rc.buffer[i] != '\n' {
					Goose.Storage.Logf(1, "Error reading %s linefeed on chunk[%d]: %s", rc.key, rc.chunk, UnexpectedCharacter)
					return 0, UnexpectedCharacter
				}
				break
			} else if rc.buffer[i] >= 'a' && rc.buffer[i] <= 'f' {
				chunkSize <<= 4
				chunkSize += int(rc.buffer[i] - 'a' + 10)
			} else if rc.buffer[i] >= 'A' && rc.buffer[i] <= 'F' {
				chunkSize <<= 4
				chunkSize += int(rc.buffer[i] - 'A' + 10)
			} else if rc.buffer[i] >= '0' && rc.buffer[i] <= '9' {
				chunkSize <<= 4
				chunkSize += int(rc.buffer[i] - '0')
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
		}

		Goose.Storage.Logf(4, "sz=%d", sz)
		rc.consumed = 0
		sz -= 38
		rc.rd = bytes.NewReader(rc.buffer[:sz])
//		rc.rd = bytes.NewReader(rc.buffer[8:rc.chunkSize+8])
		Goose.Storage.Logf(0, "Removing trailer: %d % 2x .. % 2x", rc.chunk,  rc.buffer[:8], rc.buffer[sz:])
		Goose.Storage.Logf(0, "Removing trailer: %d %s .. %s", rc.chunk, rc.buffer[:8], rc.buffer[sz:])
	}

	n, err = rc.rd.Read(buf)
	rc.consumed += int32(n)
	Goose.Storage.Logf(4, "Read %d bytes: %s", n, err)
	if err == io.EOF {
		Goose.Storage.Logf(3,"EOF")
		if rc.chunk < rc.chunks {
			Goose.Storage.Logf(3,"rc.chunk %d < rc.chunks %d", rc.chunk, rc.chunks)
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
