package fake

import (
	"os"
	"testing"
)

// UseFakeStderr swaps stderr with fake one.
func UseFakeStderr(t *testing.T) SwappedFile {
	t.Helper()

	return swap(t,
		func() *os.File {
			return os.Stderr
		},
		func(f *os.File) {
			os.Stderr = f
		},
	)
}
