package fake

import (
	"os"
	"testing"
)

func UseFakeStdout(t *testing.T) Std {
	t.Helper()

	return useStd(t,
		func() *os.File {
			return os.Stdout
		},
		func(f *os.File) {
			os.Stdout = f
		},
	)
}
