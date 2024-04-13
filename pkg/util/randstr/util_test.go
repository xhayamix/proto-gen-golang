package randstr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	length := 8
	str, err := RandomString(length)
	assert.NoError(t, err)
	assert.Len(t, str, length)
}
