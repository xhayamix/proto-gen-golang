package version

import (
	"regexp"
	"sort"

	"github.com/hashicorp/go-version"

	"github.com/xhayamix/proto-gen-golang/pkg/cerrors"
)

var (
	// masterのみ一番上に固定
	masterVer     = "master"
	versionRegexp = regexp.MustCompile("^" + version.VersionRegexpRaw + "$")
)

func NewVersionMap(l []string) (map[string]*version.Version, error) {
	res := make(map[string]*version.Version, len(l))
	for _, s := range l {
		// masterなどバージョンフォーマットでない場合はスキップ
		if versionRegexp.FindStringSubmatch(s) == nil {
			continue
		}

		v, err := version.NewVersion(s)
		if err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
		res[s] = v
	}

	return res, nil
}

// Sort master -> version(降順) -> それ以外の文字列 の順でソート
func Sort(l []string) ([]string, error) {
	res := make([]string, len(l))
	copy(res, l)

	versionMap, err := NewVersionMap(res)
	if err != nil {
		return nil, cerrors.Stack(err)
	}

	sort.Slice(res, func(i, j int) bool {
		return GreaterThan(res[i], res[j], versionMap)
	})

	return res, nil
}

// GreaterThan master -> version(降順) -> それ以外の文字列の比較
func GreaterThan(iV, jV string, versionMap map[string]*version.Version) bool {
	// masterであれば優先
	if iV == masterVer {
		return true
	}
	if jV == masterVer {
		return false
	}

	iVer, iVersioned := versionMap[iV]
	jVer, jVersioned := versionMap[jV]
	if iVersioned && jVersioned {
		// patchまで同じ場合はsuffix(prerelease) がついていない方を優先
		if iVer.Core().Compare(jVer.Core()) == 0 {
			if iVer.Prerelease() == "" {
				return true
			}
			if jVer.Prerelease() == "" {
				return false
			}
		}

		return iVer.GreaterThan(jVer)
	}
	if !iVersioned && !jVersioned {
		// 文字列ソート
		return iV < jV
	}

	// バージョン優先
	return iVersioned
}
