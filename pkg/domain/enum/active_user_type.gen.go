// Code generated by protoc-gen-campus. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package enum

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/scylladb/go-set/i32set"

	cstrings "github.com/QualiArts/campus-server/pkg/util/strings"
)

const ActiveUserTypeName = "ActiveUserType"

type ActiveUserType int32

const (
	// 仮登録
	ActiveUserType_provisional ActiveUserType = 1
	// アクティブユーザー
	ActiveUserType_Active ActiveUserType = 2
	// アカウント削除済み
	ActiveUserType_Deleted ActiveUserType = 99
)

var ActiveUserTypeMap = map[string]int32{
	"provisional": 1,
	"Active":      2,
	"Deleted":     99,
}

func (e ActiveUserType) Int() int {
	return int(e)
}

func (e ActiveUserType) Int32() int32 {
	return int32(e)
}

func (e ActiveUserType) Int64() int64 {
	return int64(e)
}

func (e ActiveUserType) String() string {
	switch e {
	case ActiveUserType_provisional:
		return "provisional"
	case ActiveUserType_Active:
		return "Active"
	case ActiveUserType_Deleted:
		return "Deleted"
	case 0:
		return ""
	}
	return strconv.FormatInt(int64(e), 10)
}

func (e ActiveUserType) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *ActiveUserType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	case "provisional":
		*e = ActiveUserType_provisional
	case "active":
		*e = ActiveUserType_Active
	case "deleted":
		*e = ActiveUserType_Deleted
	default:
		i, _ := strconv.ParseInt(s, 10, 64)
		*e = ActiveUserType(int32(i))
	}
	return nil
}

func (e ActiveUserType) EncodeSpanner() (interface{}, error) {
	return e.Int64(), nil
}

func (e *ActiveUserType) DecodeSpanner(val interface{}) error {
	strVal, ok := val.(string)
	if !ok {
		return errors.New(fmt.Sprintf("ActiveUserType.DecodeSpanner failed. %#v", val))
	}
	i, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("ActiveUserType.DecodeSpanner failed. %#v, %#v", val, err))
	}
	*e = ActiveUserType(i)
	return nil
}

func (e ActiveUserType) Validate() bool {
	switch e {
	case ActiveUserType_provisional:
		return true
	case ActiveUserType_Active:
		return true
	case ActiveUserType_Deleted:
		return true
	}
	return false
}

type ActiveUserTypeSlice []ActiveUserType

func (e ActiveUserTypeSlice) First() ActiveUserType {
	if len(e) == 0 {
		return 0
	}
	return e[0]
}

func (e ActiveUserTypeSlice) Last() ActiveUserType {
	if len(e) == 0 {
		return 0
	}
	return e[len(e)-1]
}

func (e ActiveUserTypeSlice) Set() *i32set.Set {
	set := i32set.New()
	for _, i := range e {
		set.Add(i.Int32())
	}
	return set
}

func (e ActiveUserTypeSlice) Each(f func(Enum) bool) {
	for _, i := range e {
		if !f(i) {
			break
		}
	}
}

func (e ActiveUserTypeSlice) Size() int {
	return len(e)
}

func (e ActiveUserTypeSlice) Includes(typ ActiveUserType) bool {
	for _, i := range e {
		if i == typ {
			return true
		}
	}

	return false
}

func (e ActiveUserTypeSlice) Strings() []string {
	s := make([]string, 0, len(e))
	for _, i := range e {
		s = append(s, i.String())
	}
	return s
}

func (e ActiveUserTypeSlice) ToSet() ActiveUserTypeSet {
	s := make(ActiveUserTypeSet, len(e))
	for _, i := range e {
		s.Add(i)
	}
	return s
}

func (e ActiveUserTypeSlice) EncodeSpanner() (interface{}, error) {
	ret := make([]int64, 0, e.Size())
	for _, i := range e {
		ret = append(ret, i.Int64())
	}
	return ret, nil
}

func (e ActiveUserTypeSlice) Validate() bool {
	for _, i := range e {
		if !i.Validate() {
			return false
		}
	}
	return true
}

var ActiveUserTypeValues = ActiveUserTypeSlice{
	ActiveUserType_provisional,
	ActiveUserType_Active,
	ActiveUserType_Deleted,
}

type ActiveUserTypeSet map[ActiveUserType]struct{}

func (s ActiveUserTypeSet) Has(e ActiveUserType) bool {
	_, ok := s[e]
	return ok
}

func (s ActiveUserTypeSet) Size() int {
	return len(s)
}

func (s ActiveUserTypeSet) Add(e ActiveUserType) {
	s[e] = struct{}{}
}

func (s ActiveUserTypeSet) ToSlice() ActiveUserTypeSlice {
	slice := make(ActiveUserTypeSlice, 0, len(s))
	for _, v := range ActiveUserTypeValues {
		if s.Has(v) {
			slice = append(slice, v)
		}
	}
	return slice
}

type ActiveUserTypeCommaSeparated string

func (e ActiveUserTypeCommaSeparated) Split() (Enums, []string) {
	var errs []string
	list := cstrings.SplitComma(string(e))
	res := make(ActiveUserTypeSlice, 0, len(list))

	for _, str := range list {
		i, err := strconv.Atoi(str)
		if err != nil {
			errs = append(errs, fmt.Sprintf("failed to convert enum.ActiveUserType. %v\n", str))
			continue
		}
		res = append(res, ActiveUserType(i))
	}

	return res, errs
}

func (e ActiveUserTypeCommaSeparated) String() string {
	list, _ := e.Split()
	res := make([]string, 0, list.Size())
	list.Each(func(i Enum) bool {
		res = append(res, i.String())
		return true
	})
	return strings.Join(res, ",")
}

func (e ActiveUserTypeCommaSeparated) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *ActiveUserTypeCommaSeparated) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	list := cstrings.SplitComma(s)
	res := make([]string, 0, len(s))
	for _, str := range list {
		var i ActiveUserType
		err := i.UnmarshalJSON([]byte(`"` + str + `"`))
		if err != nil {
			return err
		}
		res = append(res, strconv.FormatInt(i.Int64(), 10))
	}
	*e = ActiveUserTypeCommaSeparated(strings.Join(res, ","))
	return nil
}

var ActiveUserTypeValueDetails = ValueDetails{
	{Type: ActiveUserType_provisional, Comment: "仮登録"},
	{Type: ActiveUserType_Active, Comment: "アクティブユーザー"},
	{Type: ActiveUserType_Deleted, Comment: "アカウント削除済み"},
}
