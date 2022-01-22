package benchmark

type DiscardWriter struct{}

func (d DiscardWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
