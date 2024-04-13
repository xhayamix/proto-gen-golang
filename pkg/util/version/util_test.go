package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Sort(t *testing.T) {
	res, err := Sort([]string{
		"v0.8.1",
		"v0.8.1-test1",
		"v0.8.1-test-2",
		"v0.8.2",
		"v0.8.1-test2",
		"v0.8.1-test-1",
		"develop",
		"master",
	})
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"master",
		"v0.8.2",
		"v0.8.1",
		"v0.8.1-test2",
		"v0.8.1-test1",
		"v0.8.1-test-2",
		"v0.8.1-test-1",
		"develop",
	}, res)
}
