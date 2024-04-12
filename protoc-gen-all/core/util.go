package core

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"text/template"

	"github.com/huandu/xstrings"
)

var CommentReplacer = strings.NewReplacer("//", "", " ", "", "\n", "")

var IDRegExp = regexp.MustCompile(`^(.*)(Id|id)(s?|\d+|s\d+)$`)
var UUIDRegExp = regexp.MustCompile(`^(.*)(Uuid|uuid)(s?|\d+|s\d+)$`)

const captureLength = 4

type mapString struct {
	value map[string]string
	mutex *sync.RWMutex
	// 累計呼び出し回数
	totalCount int64
	// キャッシュヒットした回数
	hitCount int64
}

func newMapString() *mapString {
	return &mapString{
		// アロケーションしないようにいい感じの数を確保
		value: make(map[string]string, 10000),
		mutex: &sync.RWMutex{},
	}
}

func (m *mapString) Load(key string) (string, bool) {
	m.mutex.Lock()
	v, ok := m.value[key]
	m.totalCount++
	if ok {
		m.hitCount++
	}
	m.mutex.Unlock()
	return v, ok
}

func (m *mapString) Store(key, value string) {
	m.mutex.Lock()
	m.value[key] = value
	m.mutex.Unlock()
}

func (m *mapString) GetTotalCount() int64 {
	m.mutex.RLock()
	v := m.totalCount
	m.mutex.RUnlock()
	return v
}

func (m *mapString) GetHitCount() int64 {
	m.mutex.RLock()
	v := m.hitCount
	m.mutex.RUnlock()
	return v
}

var toSnakeCaseCache = newMapString()

func ToSnakeCase(str string) string {
	if v, ok := toSnakeCaseCache.Load(str); ok {
		return v
	}

	var result string
	if strings.EqualFold(strings.ToLower(str), "i18n") {
		// "I18n"が"i_18n"になってしまうので対応
		result = "i18n"
	} else {
		// "IDs"が"i_ds"になってしまうため事前にreplace
		result = xstrings.ToSnakeCase(strings.ReplaceAll(str, "IDs", "Ids"))
	}

	toSnakeCaseCache.Store(str, result)
	return result
}

var toKebabCaseCache = newMapString()

// ToKebabCase
//
//	user_id -> user-id
func ToKebabCase(str string) string {
	if v, ok := toKebabCaseCache.Load(str); ok {
		return v
	}

	snakeStr := xstrings.ToSnakeCase(str)
	result := xstrings.ToKebabCase(snakeStr)

	toKebabCaseCache.Store(str, result)
	return result
}

var toCamelCaseCache = newMapString()

// ToCamelCase
//
//	user_id -> userId
func ToCamelCase(str string) string {
	if v, ok := toCamelCaseCache.Load(str); ok {
		return v
	}

	snakeStr := xstrings.ToSnakeCase(str)
	camelStr := xstrings.ToCamelCase(snakeStr)
	result := xstrings.FirstRuneToLower(camelStr)

	toCamelCaseCache.Store(str, result)
	return result
}

var toGolangCamelCaseCache = newMapString()

// ToGolangCamelCase
//
//	user_id -> userID
func ToGolangCamelCase(str string) string {
	if v, ok := toGolangCamelCaseCache.Load(str); ok {
		return v
	}

	snakeStr := xstrings.ToSnakeCase(str)
	camelStr := xstrings.ToCamelCase(snakeStr)

	var result string
	if snakeStr == "uuid" {
		result = "uuid"
	} else if snakeStr == "id" {
		result = "id"
	} else if captures := UUIDRegExp.FindStringSubmatch(camelStr); len(captures) == captureLength {
		result = xstrings.FirstRuneToLower(captures[1] + "UUID" + captures[3])
	} else if captures := IDRegExp.FindStringSubmatch(camelStr); len(captures) == captureLength {
		result = xstrings.FirstRuneToLower(captures[1] + "ID" + captures[3])
	} else {
		result = xstrings.FirstRuneToLower(camelStr)
	}

	toGolangCamelCaseCache.Store(str, result)
	return result
}

var toPascalCaseCache = newMapString()

// ToPascalCase
//
//	user_id -> UserId
func ToPascalCase(str string) string {
	if v, ok := toPascalCaseCache.Load(str); ok {
		return v
	}

	snakeStr := xstrings.ToSnakeCase(str)
	result := xstrings.ToCamelCase(snakeStr)

	toPascalCaseCache.Store(str, result)
	return result
}

var toGolangPascalCaseCache = newMapString()

// ToGolangPascalCase
//
//	user_id -> UserID
func ToGolangPascalCase(str string) string {
	if v, ok := toGolangPascalCaseCache.Load(str); ok {
		return v
	}

	snakeStr := xstrings.ToSnakeCase(str)
	camelStr := xstrings.ToCamelCase(snakeStr)

	var result string
	if snakeStr == "uuid" {
		result = "UUID"
	} else if snakeStr == "id" {
		result = "ID"
	} else if captures := UUIDRegExp.FindStringSubmatch(camelStr); len(captures) == captureLength {
		result = xstrings.FirstRuneToUpper(captures[1] + "UUID" + captures[3])
	} else if captures := IDRegExp.FindStringSubmatch(camelStr); len(captures) == captureLength {
		result = xstrings.FirstRuneToUpper(captures[1] + "ID" + captures[3])
	} else {
		result = xstrings.FirstRuneToUpper(camelStr)
	}

	toGolangPascalCaseCache.Store(str, result)
	return result
}

var toPkgNameCache = newMapString()

// ToPkgName
//
//	user_id -> userid
func ToPkgName(str string) string {
	if v, ok := toPkgNameCache.Load(str); ok {
		return v
	}

	snakeStr := xstrings.ToSnakeCase(xstrings.ToCamelCase(str))
	result := strings.ReplaceAll(snakeStr, "_", "")

	toPkgNameCache.Store(str, result)
	return result
}

func IsTimeField(snakeName string) bool {
	return strings.HasSuffix(snakeName, "_time")
}

func IsAdminTimeField(snakeName string) bool {
	return strings.HasSuffix(snakeName, "_at")
}

func IsMasterTagKind(snakeName string) bool {
	return strings.HasPrefix(snakeName, "master_tag")
}

func IsMasterVersion(snakeName string) bool {
	return snakeName == "master_version"
}

const baseTemplateString = `
{{ define "autogen_comment" -}}
// Code generated by protoc-gen-all. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.
{{ end }}
`

var baseTemplate = template.Must(template.New("").Parse(baseTemplateString))

func GetBaseTemplate() *template.Template {
	return template.Must(baseTemplate.Clone())
}

func GetCacheInfo() string {
	var f = func(a, b int64) float64 {
		if b == 0 {
			return 0
		}
		return (float64(a) / float64(b)) * 100
	}
	return fmt.Sprintf(`
	toSnakeCaseCache: %v/%v CHR=%.3f%%
	toKebabCaseCache: %v/%v CHR=%.3f%%
	toCamelCaseCache: %v/%v CHR=%.3f%%
	toGolangCamelCaseCache: %v/%v CHR=%.3f%%
	toPascalCaseCache: %v/%v CHR=%.3f%%
	toGolangPascalCaseCache: %v/%v CHR=%.3f%%
	toPkgNameCache: %v/%v CHR=%.3f%%
	`,
		toSnakeCaseCache.GetHitCount(), toSnakeCaseCache.GetTotalCount(), f(toSnakeCaseCache.GetHitCount(), toSnakeCaseCache.GetTotalCount()),
		toKebabCaseCache.GetHitCount(), toKebabCaseCache.GetTotalCount(), f(toKebabCaseCache.GetHitCount(), toKebabCaseCache.GetTotalCount()),
		toCamelCaseCache.GetHitCount(), toCamelCaseCache.GetTotalCount(), f(toCamelCaseCache.GetHitCount(), toCamelCaseCache.GetTotalCount()),
		toGolangCamelCaseCache.GetHitCount(), toGolangCamelCaseCache.GetTotalCount(), f(toGolangCamelCaseCache.GetHitCount(), toGolangCamelCaseCache.GetTotalCount()),
		toPascalCaseCache.GetHitCount(), toPascalCaseCache.GetTotalCount(), f(toPascalCaseCache.GetHitCount(), toPascalCaseCache.GetTotalCount()),
		toGolangPascalCaseCache.GetHitCount(), toGolangPascalCaseCache.GetTotalCount(), f(toGolangPascalCaseCache.GetHitCount(), toGolangPascalCaseCache.GetTotalCount()),
		toPkgNameCache.GetHitCount(), toPkgNameCache.GetTotalCount(), f(toPkgNameCache.GetHitCount(), toPkgNameCache.GetTotalCount()),
	)
}

// JoinPath 空の文字列は無視し、パスの形に結合する
func JoinPath(paths ...string) string {
	joinPaths := make([]string, 0, len(paths))
	for _, path := range paths {
		if path == "" {
			continue
		}
		joinPaths = append(joinPaths, path)
	}

	return strings.Join(joinPaths, "/")
}
