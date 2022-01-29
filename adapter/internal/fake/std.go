package fake

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func swap(t *testing.T, get func() *os.File, set func(*os.File)) SwappedFile {
	t.Helper()

	prev := get()
	tmpFile, err := ioutil.TempFile("", "")
	require.NoError(t, err)

	set(tmpFile)

	return SwappedFile{
		original: prev,
		current:  tmpFile,
		set:      set,
	}
}

type SwappedFile struct {
	original *os.File
	current  *os.File
	set      func(file *os.File)
}

// Release brings back the original file.
func (f SwappedFile) Release() {
	f.set(f.original)
}

// String returns the entire contents of current file.
func (f SwappedFile) String(t *testing.T) string {
	t.Helper()

	line, err := ioutil.ReadFile(f.current.Name())
	require.NoError(t, err)

	return string(line)
}
