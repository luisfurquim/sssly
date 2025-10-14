package sssly

import (
	"io"
	"sync"
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (s *Sssly) UploadFromReader(key string, rd io.Reader, sz int64) error {
	var (
		err error
		MultipartUpload *s3.CreateMultipartUploadOutput
		ctx context.Context
		MaxChunk int
		uploadID string
		i int32
		buf []byte
		n int
		wg sync.WaitGroup
//		mtx sync.Mutex
		parts []types.CompletedPart
	)

	if s.MaxChunk > 0 {
		MaxChunk = s.MaxChunk
	} else {
		MaxChunk =  1024
	}

	if sz < int64(MaxChunk) {
		Goose.Storage.Logf(3, "Singlepart upload started size: %s", sz)
		_, err = s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(s.BasePath + key),
			ContentLength: &sz,
			Body:   rd,
		})

		if err != nil {
			Goose.Storage.Logf(1, "Error opening %s for write: %s", key, err)
		}
	} else {

		ctx = context.Background()

		// 2. Iniciar multipart upload
		MultipartUpload, err = s.Client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(s.BasePath + key),
//			// Opcional: definir metadados, tipo de conteúdo, etc.
//			ContentType: aws.String("application/octet-stream"),
//			Metadata: map[string]string{
//				"original-filename": filepath.Base(filePath),
//				"file-size":        fmt.Sprintf("%d", fileSize),
//			},
		})
		if err != nil {
			Goose.Storage.Logf(1, "%s: %s", ErrStartingMultipartUpload, err)
			return err
		}

		uploadID = *MultipartUpload.UploadId
		Goose.Storage.Logf(3, "Multipart upload started. Upload ID: %s", uploadID)
		parts = make([]types.CompletedPart, (sz/int64(MaxChunk))*2 + 1)

		for sz > 0 {
			// yes, this is the right place to increment, before the loop body, because part numbers range from 1~10000 and not from 0~9999...
			i++
			if sz < int64(MaxChunk) {
				buf = make([]byte, sz)
			} else {
				buf = make([]byte, MaxChunk)
			}

			n, err = rd.Read(buf)
			if err != nil {
				Goose.Storage.Logf(1, "Error reading data: %s", err)
				s.Client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
					Bucket:     aws.String(s.Bucket),
					Key:        aws.String(s.BasePath + key),
					UploadId:	aws.String(uploadID),
				})
				return err
			}
			
			func(part int32, size int64, buffer []byte) {
				var e error
				var upl *s3.UploadPartOutput

				Goose.Storage.Logf(3,"%d: buffersize:%d", part, size)

				wg.Add(1)
				defer wg.Done()

				upl, e = s.Client.UploadPart(
					ctx,
					&s3.UploadPartInput{
						Bucket:     aws.String(s.Bucket),
						Key:        aws.String(s.BasePath + key),
						PartNumber: &part,
						UploadId:   aws.String(uploadID),
						Body:       bytes.NewReader(buffer),
						ContentLength: &size,
					},
//					optFns ...func(*Options),
				)

				if e != nil {
					err = e
					Goose.Storage.Logf(1, "Error uploading chunk: %s", err)
				} else {
					Goose.Storage.Logf(3, "%d: %s", part, *upl.ETag)
					parts[part] = types.CompletedPart{
						ETag:       upl.ETag,
						PartNumber: aws.Int32(part),
					}
				}
			}(i, int64(n), buf[:n])

			sz -= int64(n)
		}

		wg.Wait()

		if err != nil {
			s.Client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
				Bucket:     aws.String(s.Bucket),
				Key:        aws.String(s.BasePath + key),
				UploadId:	aws.String(uploadID),
			})
			return err
		}

		_, err = s.Client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
			Bucket:     aws.String(s.Bucket),
			Key:        aws.String(s.BasePath + key),
			UploadId:	aws.String(uploadID),
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: parts[1:i+1],
			},
		})

		if err != nil {
			Goose.Storage.Logf(1, "Error completing upload: %s", err)
			return err
		}
	}

	return err
}

