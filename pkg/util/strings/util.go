package strings

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	alphabetDigitRegExp = regexp.MustCompile("^([0-9a-zA-Z])+$")
	newLineRegExp       = regexp.MustCompile("\r\n|\n")
	symbolRegExp        = regexp.MustCompile(`[!$%^&*()_+\|~\-=` + "`" + `{}\[\]:";'<>?,./\p{So}]`)
)

func IsAlphabetDigit(text string) bool {
	return alphabetDigitRegExp.MatchString(text)
}

func SplitComma(text string) []string {
	return Split(text, ",")
}

// Split 空文字のときに空配列になる
func Split(text, separator string) []string {
	if text == "" {
		return []string{}
	}
	return strings.Split(text, separator)
}

// SplitN Splitできなくても指定された要素数で空文字を埋める
func SplitN(text, separator string, n int) []string {
	if text == "" {
		return make([]string, n)
	}
	l := strings.SplitN(text, separator, n)
	return l[:n]
}

func SplitCommaToInt32(str string) ([]int32, error) {
	strs := SplitComma(str)
	ints := make([]int32, 0, len(strs))
	for _, s := range strs {
		if s == "" {
			continue
		}
		i, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, err
		}
		ints = append(ints, int32(i))
	}
	return ints, nil
}

func SplitCommaToInt64(str string) ([]int64, error) {
	strs := SplitComma(str)
	ints := make([]int64, 0, len(strs))
	for _, s := range strs {
		if s == "" {
			continue
		}
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		ints = append(ints, i)
	}
	return ints, nil
}

func SplitCommaToBool(str string) ([]bool, error) {
	strs := SplitComma(str)
	bools := make([]bool, 0, len(strs))
	for _, s := range strs {
		if s == "" {
			continue
		}
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, err
		}
		bools = append(bools, b)
	}
	return bools, nil
}

// SplitNewLine 空文字のときに空配列になる
func SplitNewLine(text string) []string {
	return newLineRegExp.Split(text, -1)
}

// JoinNewLine 改行文字列で結合する
func JoinNewLine(texts []string) string {
	return strings.Join(texts, "\n")
}

// ParseInt 文字列を数値に変換 空文字の場合は0を返す
func ParseInt(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseInt(s, 10, 64)
}

// ParseBool 文字列をBoolに変換 空文字の場合はfalseを返す
func ParseBool(s string) bool {
	// strconv.ParseBoolは "1" でもtrueになる
	return s == "true"
}

// IsContainsNewLine 改行文字が含まれているか
func IsContainsNewLine(text string) bool {
	ret := strings.ToValidUTF8(text, "")

	return newLineRegExp.MatchString(ret)
}

// IsContainsSymbol 記号・絵文字が含まれているか
func IsContainsSymbol(text string) bool {
	ret := strings.ToValidUTF8(text, "")

	return symbolRegExp.MatchString(ret)
}
