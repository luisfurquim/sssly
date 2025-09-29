package sssly

func (rc *ReadCloser) Close() error {
	return rc.rd.Close()
}
