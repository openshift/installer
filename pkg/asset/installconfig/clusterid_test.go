package installconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generateInfraID(t *testing.T) {
	tests := []struct {
		input string

		expLen     int
		expNonRand string
	}{{
		input:      "qwertyuiop",
		expLen:     10 + randomLen + 1,
		expNonRand: "qwertyuiop",
	}, {
		input:      "qwertyuiopasdfghjklzxcvbnm",
		expLen:     27,
		expNonRand: "qwertyuiopasdfghjklzx",
	}, {
		input:      "qwe.rty.@iop!",
		expLen:     13 + randomLen + 1,
		expNonRand: "qwe-rty--iop-",
	}}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			got := generateInfraID(test.input, 27)
			t.Log("InfraID", got)
			assert.Equal(t, test.expLen, len(got))
			assert.Equal(t, test.expNonRand, got[:len(got)-randomLen-1])
		})
	}
}
