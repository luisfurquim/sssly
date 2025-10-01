package sssly

func (rc *ReadCloser) Close() error {
	if rc.remReader == nil {
		return nil
	}

	rc.buffer = nil

	return rc.remReader.Close()
}
