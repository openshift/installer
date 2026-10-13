package lineprinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimmer(t *testing.T) {
	print := &printer{}
	trimmer := &Trimmer{WrappedPrint: print.print}
	trimmer.Print("Hello\n", "World\n")
	trimmer.Print(123)
	assert.Equal(
		t,
		[][]interface{}{
			{"Hello\n", "World"},
			{123},
		},
		print.data,
	)
}
