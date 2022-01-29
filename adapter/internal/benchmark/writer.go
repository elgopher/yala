// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package benchmark

// DiscardWriter is io.Writer discarding all written bytes.
type DiscardWriter struct{}

func (d DiscardWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
