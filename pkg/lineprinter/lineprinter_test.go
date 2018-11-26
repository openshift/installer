package lineprinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type printer struct {
	data [][]interface{}
}

func (p *printer) print(args ...interface{}) {
	p.data = append(p.data, args)
}

func TestLinePrinter(t *testing.T) {
	print := &printer{}
	lp := &LinePrinter{Print: print.print}
	data := []byte("Hello\nWorld\nAnd more")
	n, err := lp.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(data), n)
	assert.Equal(
		t,
		[][]interface{}{
			{"Hello\n"},
			{"World\n"},
		},
		print.data,
	)
	print.data = [][]interface{}{}

	err = lp.Close()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(
		t,
		[][]interface{}{
			{"And more"},
		},
		print.data,
	)
}
