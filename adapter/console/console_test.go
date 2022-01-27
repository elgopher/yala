// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package console_test

import (
	"testing"

	"github.com/elgopher/yala/adapter/console"
	"github.com/stretchr/testify/assert"
)

func TestWriterPrinter_Println(t *testing.T) {
	t.Run("should not panic when Writer is nil", func(t *testing.T) {
		p := console.WriterPrinter{Writer: nil}
		assert.NotPanics(t, func() {
			p.Println(0, "")
		})
	})
}
