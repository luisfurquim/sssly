package sssly

func (lrc *LargeReadCloser) WriteAt(p []byte, off int64) (int, error) {
	var (
		buf []byte
		ok bool
	)

	buf = make([]byte, len(p))
	copy(buf, p)
	ok = true

	lrc.mtx.Lock()
	if lrc.off == off {
		lrc.ready = append(lrc.ready, buf)
		for ok {
			lrc.off += int64(len(buf))
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
