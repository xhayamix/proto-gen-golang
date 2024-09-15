// Code generated by protoc-gen-all. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package transaction

import (
	"fmt"
	"strings"
	"time"

	"github.com/scylladb/go-set/strset"

	"github.com/xhayamix/proto-gen-golang/pkg/domain/constant"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/dto/column"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/entity"
)

const (
	UserTableName = "user"
	UserComment   = "ユーザー"
)

// ユーザー
type User struct {
	// ID
	ID string `json:"id,omitempty"`
	// ユーザー作成ID
	UserID string `json:"user_id,omitempty"`
	// メールアドレス
	Email string `json:"email,omitempty"`
	// パスワード
	Password string `json:"password,omitempty"`
	// ユーザー名
	Name string `json:"name,omitempty"`
	// プロフィール
	Profile string `json:"profile,omitempty"`
	// アイコン画像パス
	IconImg string `json:"icon_img,omitempty"`
	// ヘッダー画像パス
	HeaderImg string `json:"header_img,omitempty"`
	// 作成日時
	CreatedAt time.Time `json:"created_at,omitempty"`
	// 更新日時
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (e *User) GetPK() *UserPK {
	return &UserPK{
		ID: e.ID,
	}
}

func (e *User) PK() string {
	return e.GetPK().Key()
}

func (e *User) ToKeyValue() map[string]interface{} {
	return map[string]interface{}{
		"id":         e.ID,
		"user_id":    e.UserID,
		"email":      e.Email,
		"password":   e.Password,
		"name":       e.Name,
		"profile":    e.Profile,
		"icon_img":   e.IconImg,
		"header_img": e.HeaderImg,
		"created_at": e.CreatedAt,
		"updated_at": e.UpdatedAt,
	}
}

func (e *User) GetTypeMap() map[string]string {
	return map[string]string{
		"id":         "string",
		"user_id":    "string",
		"email":      "string",
		"password":   "string",
		"name":       "string",
		"profile":    "string",
		"icon_img":   "string",
		"header_img": "string",
		"created_at": "time.Time",
		"updated_at": "time.Time",
	}
}

func (e *User) SetKeyValue(columns []string, values []string) []string {
	errs := make([]string, 0, len(columns))
	for index, column := range columns {
		if len(values) <= index {
			break
		}
		value := values[index]
		switch column {
		case "id":
			e.ID = value
		case "user_id":
			e.UserID = value
		case "email":
			e.Email = value
		case "password":
			e.Password = value
		case "name":
			e.Name = value
		case "profile":
			e.Profile = value
		case "icon_img":
			e.IconImg = value
		case "header_img":
			e.HeaderImg = value
		case "created_at":
			if value != "" {
				var v time.Time
				var err error
				switch {
				case constant.NormalDatetimeRegExp.MatchString(value):
					v, err = time.ParseInLocation("2006/01/02 15:04:05", value, time.Local)
				case constant.HyphenDatetimeRegExp.MatchString(value):
					v, err = time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
				default:
					v, err = time.Parse(time.RFC3339, value)
				}
				if err != nil {
					errs = append(errs, fmt.Sprintf("created_at: time.Time parsing %#v: invalid syntax.", value))
				}
				e.CreatedAt = v
			}
		case "updated_at":
			if value != "" {
				var v time.Time
				var err error
				switch {
				case constant.NormalDatetimeRegExp.MatchString(value):
					v, err = time.ParseInLocation("2006/01/02 15:04:05", value, time.Local)
				case constant.HyphenDatetimeRegExp.MatchString(value):
					v, err = time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
				default:
					v, err = time.Parse(time.RFC3339, value)
				}
				if err != nil {
					errs = append(errs, fmt.Sprintf("updated_at: time.Time parsing %#v: invalid syntax.", value))
				}
				e.UpdatedAt = v
			}
		}
	}
	return errs
}

func (e *User) PtrFromMapping(cols []string) []interface{} {
	ptrs := make([]interface{}, 0, len(cols))
	for _, col := range cols {
		switch col {
		case "id":
			ptrs = append(ptrs, &e.ID)
		case "user_id":
			ptrs = append(ptrs, &e.UserID)
		case "email":
			ptrs = append(ptrs, &e.Email)
		case "password":
			ptrs = append(ptrs, &e.Password)
		case "name":
			ptrs = append(ptrs, &e.Name)
		case "profile":
			ptrs = append(ptrs, &e.Profile)
		case "icon_img":
			ptrs = append(ptrs, &e.IconImg)
		case "header_img":
			ptrs = append(ptrs, &e.HeaderImg)
		case "created_at":
			ptrs = append(ptrs, &e.CreatedAt)
		case "updated_at":
			ptrs = append(ptrs, &e.UpdatedAt)
		}
	}

	return ptrs
}

type UserSlice []*User

func (s UserSlice) CreateMapByPK() UserMapByPK {
	m := make(UserMapByPK, len(s))
	for _, e := range s {
		m[e.ID] = e
	}
	return m
}

func (s UserSlice) EachRecord(iterator func(entity.Record) bool) {
	for _, e := range s {
		if !iterator(e) {
			break
		}
	}
}

type UserMapByPK map[string]*User

func (m UserMapByPK) Has(keys ...interface{}) bool {
	m0 := m
	for i, key := range keys {
		switch i {
		case 0:
			k, ok := key.(string)
			if !ok {
				return false
			}
			_, ok = m0[k]
			if !ok {
				return false
			}
		default:
			return false
		}
	}
	return true
}

type UserPK struct {
	ID string
}

type UserPKs []*UserPK

func (e *UserPK) Generate() []interface{} {
	return []interface{}{
		e.ID,
	}
}

func (e *UserPK) Key() string {
	return strings.Join([]string{
		fmt.Sprint(e.ID),
	}, ".")
}

var UserColumnName = struct {
	ID        string
	UserID    string
	Email     string
	Password  string
	Name      string
	Profile   string
	IconImg   string
	HeaderImg string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	UserID:    "user_id",
	Email:     "email",
	Password:  "password",
	Name:      "name",
	Profile:   "profile",
	IconImg:   "icon_img",
	HeaderImg: "header_img",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

var UserColumnMap = map[string]*column.Column{
	UserColumnName.ID: {
		Name:     "id",
		Type:     "string",
		PK:       true,
		Nullable: false,
		Required: false,
		Comment:  "ID",
		CSV:      false,
	},
	UserColumnName.UserID: {
		Name:     "user_id",
		Type:     "string",
		PK:       false,
		Nullable: false,
		Required: false,
		Comment:  "ユーザー作成ID",
		CSV:      false,
	},
	UserColumnName.Email: {
		Name:     "email",
		Type:     "string",
		PK:       false,
		Nullable: false,
		Required: false,
		Comment:  "メールアドレス",
		CSV:      false,
	},
	UserColumnName.Password: {
		Name:     "password",
		Type:     "string",
		PK:       false,
		Nullable: false,
		Required: false,
		Comment:  "パスワード",
		CSV:      false,
	},
	UserColumnName.Name: {
		Name:     "name",
		Type:     "string",
		PK:       false,
		Nullable: false,
		Required: false,
		Comment:  "ユーザー名",
		CSV:      false,
	},
	UserColumnName.Profile: {
		Name:     "profile",
		Type:     "string",
		PK:       false,
		Nullable: false,
		Required: false,
		Comment:  "プロフィール",
		CSV:      false,
	},
	UserColumnName.IconImg: {
		Name:     "icon_img",
		Type:     "string",
		PK:       false,
		Nullable: false,
		Required: false,
		Comment:  "アイコン画像パス",
		CSV:      false,
	},
	UserColumnName.HeaderImg: {
		Name:     "header_img",
		Type:     "string",
		PK:       false,
		Nullable: false,
		Required: false,
		Comment:  "ヘッダー画像パス",
		CSV:      false,
	},
	UserColumnName.CreatedAt: {
		Name:     "created_at",
		Type:     "time.Time",
		PK:       false,
		Nullable: false,
		Required: false,
		Comment:  "作成日時",
		CSV:      false,
	},
	UserColumnName.UpdatedAt: {
		Name:     "updated_at",
		Type:     "time.Time",
		PK:       false,
		Nullable: false,
		Required: false,
		Comment:  "更新日時",
		CSV:      false,
	},
}

var UserCols = []string{
	UserColumnName.ID,
	UserColumnName.UserID,
	UserColumnName.Email,
	UserColumnName.Password,
	UserColumnName.Name,
	UserColumnName.Profile,
	UserColumnName.IconImg,
	UserColumnName.HeaderImg,
	UserColumnName.CreatedAt,
	UserColumnName.UpdatedAt,
}

var UserColumns = column.Columns{
	UserColumnMap[UserColumnName.ID],
	UserColumnMap[UserColumnName.UserID],
	UserColumnMap[UserColumnName.Email],
	UserColumnMap[UserColumnName.Password],
	UserColumnMap[UserColumnName.Name],
	UserColumnMap[UserColumnName.Profile],
	UserColumnMap[UserColumnName.IconImg],
	UserColumnMap[UserColumnName.HeaderImg],
	UserColumnMap[UserColumnName.CreatedAt],
	UserColumnMap[UserColumnName.UpdatedAt],
}

var UserPKCols = []string{
	UserColumnName.ID,
}

var UserPKColumns = column.Columns{
	UserColumnMap[UserColumnName.ID],
}

var UserFKParentTable = strset.New()

var UserFKChildTable = strset.New()
