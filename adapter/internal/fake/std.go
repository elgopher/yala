package fake

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func useStd(t *testing.T, get func() *os.File, set func(*os.File)) Std {
	t.Helper()

	prev := get()
	tmpFile, err := ioutil.TempFile("", "")
	require.NoError(t, err)

	set(tmpFile)

	return Std{
		prev:    prev,
		current: tmpFile,
		set:     set,
	}
}

type Std struct {
	prev    *os.File
	current *os.File
	set     func(file *os.File)
}

func (f Std) Release() {
	f.set(f.prev)
}

func (f Std) FirstLine(t *testing.T) string {
	t.Helper()

	line, err := ioutil.ReadFile(f.current.Name())
	require.NoError(t, err)

	return string(line)
}
