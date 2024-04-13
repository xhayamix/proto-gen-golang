package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAlphabetDigit(t *testing.T) {
	assert.True(t, IsAlphabetDigit("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"))

	invalidChars := []string{"あ", "ア", "亜", "-", "+", "=", "/", "*", " ", ""}
	for _, char := range invalidChars {
		assert.False(t, IsAlphabetDigit(char))
	}
}

func TestSplitN(t *testing.T) {
	for _, tt := range []struct {
		text      string
		separator string
		n         int
		expected  []string
	}{
		{
			text:      "1,2,3,4",
			separator: ",",
			n:         4,
			expected:  []string{"1", "2", "3", "4"},
		},
		{
			text:      "1,2",
			separator: ",",
			n:         4,
			expected:  []string{"1", "2", "", ""},
		},
		{
			text:      "",
			separator: ",",
			n:         4,
			expected:  []string{"", "", "", ""},
		},
		{
			text:      "1.2.3.4",
			separator: ",",
			n:         4,
			expected:  []string{"1.2.3.4", "", "", ""},
		},
		{
			text:      "1.2",
			separator: ".",
			n:         2,
			expected:  []string{"1", "2"},
		},
	} {
		assert.Equal(t, tt.expected, SplitN(tt.text, tt.separator, tt.n))
	}
}

func TestSplitNewLine(t *testing.T) {
	t.Run("正常: 空文字", func(t *testing.T) {
		assert.Equal(t, []string{""}, SplitNewLine(""))
	})

	t.Run("正常: CRLFとLF", func(t *testing.T) {
		assert.Equal(t, []string{"aa", "b", ""}, SplitNewLine("aa\r\nb\n"))
	})
}

func TestJoinNewLine(t *testing.T) {
	t.Run("正常: LF", func(t *testing.T) {
		assert.Equal(t, "aa\nb\n", JoinNewLine([]string{"aa", "b", ""}))
	})
}

func TestParseInt(t *testing.T) {
	res, err := ParseInt("1")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res)

	res, err = ParseInt("")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), res)
}

func TestParseBool(t *testing.T) {
	assert.True(t, ParseBool("true"))
	assert.False(t, ParseBool("false"))
	assert.False(t, ParseBool("1"))
	assert.False(t, ParseBool(""))
}

func Test_IsContainsNewLine(t *testing.T) {
	assert.False(t, IsContainsNewLine("abc"))
	assert.True(t, IsContainsNewLine("abc\n"))
	assert.True(t, IsContainsNewLine("\nabc"))
}

func Test_IsContainsSymbol(t *testing.T) {
	for _, str := range "09azAZあア亜" {
		assert.False(t, IsContainsSymbol(string(str)))
	}
	// 許容する記号
	for _, str := range "！？ー〜「」＝＠＃％＆（）￥＊＋ー＜＞。、＿：；…・｜＄”’｀｛｝＾／＼" {
		assert.False(t, IsContainsSymbol(string(str)))
	}

	for _, str := range "!$%^&*()_+|~-=`{}[]:\";'<>?,./" {
		assert.True(t, IsContainsSymbol(string(str)))
	}

	assert.True(t, IsContainsSymbol("🐳"))
}
