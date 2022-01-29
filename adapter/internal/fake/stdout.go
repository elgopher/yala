package fake

import (
	"os"
	"testing"
)

// UseFakeStdout swaps stdout with fake one.
func UseFakeStdout(t *testing.T) SwappedFile {
	t.Helper()

	return swap(t,
		func() *os.File {
			return os.Stdout
		},
		func(f *os.File) {
			os.Stdout = f
		},
	)
}
