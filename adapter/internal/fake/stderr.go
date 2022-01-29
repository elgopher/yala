package fake

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func UseFakeStderr(t *testing.T) Stderr {
	t.Helper()

	prevStderr := os.Stderr
	tmpFile, err := ioutil.TempFile("", "")
	require.NoError(t, err)

	os.Stderr = tmpFile

	return Stderr{
		prevStderr:    prevStderr,
		currentStderr: tmpFile,
	}
}

type Stderr struct {
	prevStderr    *os.File
	currentStderr *os.File
}

func (f Stderr) Release() {
	os.Stderr = f.prevStderr
}

func (f Stderr) FirstLine(t *testing.T) string {
	t.Helper()

	line, err := ioutil.ReadFile(f.currentStderr.Name())
	require.NoError(t, err)

	return string(line)
}
