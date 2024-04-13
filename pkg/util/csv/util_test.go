package csv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitNewLine(t *testing.T) {
	t.Run("正常: CRLFとLFを含む場合", func(t *testing.T) {
		assert.Equal(t, []string{"aa", "b"}, SplitNewLine("\naa\r\nb\n"))
	})
}
