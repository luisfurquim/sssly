package sssly

func (lrc *LargeReadCloser) Close() error {
	lrc.cli = nil
	lrc.ready = nil
	lrc.ahead = nil

	return nil
}
