package sssly

import (
	"io"
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (lrc *LargeReadCloser) WriteAt(p []byte, off int64) (int, error) {
	var (
		n int
		err error
		buf []byte
		ok book
	)

	buf = make([]byte, len(p))
	copy(buf, p)
	ok = true

	lrc.mtx.Lock()
	if lrc.off == off {
		lrc.ready = append(lrc.ready, buf)
		for ok {
			lrc.off += len(buf)
			if buf, ok = lrc.ahead[lrc.off]; ok {
				lrc.ready = append(lrc.ready, buf)
				delete(lrc.ahead, lrc.off)
			}
		}
	} else {
		lrc.ahead[off] = buf
	}
	lrc.mtx.Unlock()

	return len(p), nil
}
