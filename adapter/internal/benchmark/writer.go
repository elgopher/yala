// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package benchmark

type DiscardWriter struct{}

func (d DiscardWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
