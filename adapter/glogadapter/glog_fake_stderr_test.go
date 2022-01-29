package glogadapter_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func useFakeStderr(t *testing.T) fakeStderr {
	t.Helper()

	prevStderr := os.Stderr
	tmpFile, err := ioutil.TempFile("", "")
	require.NoError(t, err)

	os.Stderr = tmpFile

	return fakeStderr{
		prevStderr:    prevStderr,
		currentStderr: tmpFile,
	}
}

type fakeStderr struct {
	prevStderr    *os.File
	currentStderr *os.File
}

func (f fakeStderr) Release() {
	os.Stderr = f.prevStderr
}

func (f fakeStderr) FirstLine(t *testing.T) string {
	t.Helper()

	line, err := ioutil.ReadFile(f.currentStderr.Name())
	require.NoError(t, err)

	return string(line)
}
