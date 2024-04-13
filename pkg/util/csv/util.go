package csv

import (
	slicesutil "github.com/xhayamix/proto-gen-golang/pkg/util/slices"
	stringsutil "github.com/xhayamix/proto-gen-golang/pkg/util/strings"
)

func SplitNewLine(text string) []string {
	ret := slicesutil.Filter(stringsutil.SplitNewLine(text), func(e string) bool {
		return e != ""
	})
	return ret
}
