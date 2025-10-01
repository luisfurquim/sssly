package sssly

func (rc *ReadCloser) Close() error {
	if rc.remReader == nil {
		return nil
	}

	return rc.remReader.Close()
}
