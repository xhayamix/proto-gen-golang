// Code generated by protoc-gen-all. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package enum

import "github.com/scylladb/go-set/i32set"

var EnumTypeSlice = []string{
	ActiveUserTypeName,
	ErrorCodeName,
}

var EnumTypeMap = map[string]map[string]int32{
	ActiveUserTypeName: ActiveUserTypeMap,
	ErrorCodeName:      ErrorCodeMap,
}

var EnumValueDetailsMap = map[string]ValueDetails{
	ActiveUserTypeName: ActiveUserTypeValueDetails,
	ErrorCodeName:      ErrorCodeValueDetails,
}

type Enum interface {
	Int() int
	Int32() int32
	Int64() int64
	String() string
	MarshalJSON() ([]byte, error)
	Validate() bool
}

type Enums interface {
	Set() *i32set.Set
	Each(f func(Enum) bool)
	Size() int
	Validate() bool
}

type EnumCommaSeparated interface {
	Split() (Enums, []string)
	String() string
	MarshalJSON() ([]byte, error)
}

type ValueDetail struct {
	Type    Enum
	Comment string
}

type ValueDetails []*ValueDetail
