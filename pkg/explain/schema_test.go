package explain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadSchema(t *testing.T) {
	_, err := loadSchema(loadCRD(t))
	assert.NoError(t, err)
}
