package sssly

func (s *Sssly) NewWriteCloser(key string) *WriteCloser {
	return &WriteCloser{
		cli: s,
		key: key,
	}
}

