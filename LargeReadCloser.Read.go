package sssly

import (
	"io"
)

func (lrc *LargeReadCloser) Read(buf []byte) (int, error) {
	var (
		sz, n, i int
	)

	Goose.Storage.Logf(4, "Going to read: %#v %#v", lrc.ready, lrc.ahead)

	lrc.mtx.Lock()
	defer lrc.mtx.Unlock()

	if len(lrc.ready) == 0 {
		if lrc.eof {
			return 0, io.EOF
		}
		return 0, nil
	}

	if len(buf) < int(lrc.off - lrc.consumed) {
		sz = int(lrc.off - lrc.consumed)
	} else {
		sz = len(buf)
	}

	for sz > 0 {
		if sz > len(lrc.ready[0]) {
			n = len(lrc.ready[0])
		} else {
			n = sz
		}
		copy(buf[i:],lrc.ready[0][:n])
		sz -= n
		i  += n
		if n < len(lrc.ready[0]) {
			lrc.ready[0] = lrc.ready[0][n:]
		} else {
			lrc.ready = lrc.ready[1:]
			if len(lrc.ready) == 0 {
				if lrc.eof {
					return i, io.EOF
				}
				return i, nil
			}
		}
	}

	Goose.Storage.Logf(4,"done reading %d bytes", n)

	return i, nil
}
