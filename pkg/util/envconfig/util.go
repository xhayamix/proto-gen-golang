package envconfig

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/xhayamix/proto-gen-golang/pkg/cerrors"
	cstrings "github.com/xhayamix/proto-gen-golang/pkg/util/strings"
)

var gatherRegexp = regexp.MustCompile("([^A-Z]+|[A-Z]+[^A-Z]+|[A-Z]+)")
var acronymRegexp = regexp.MustCompile("([A-Z]+)([A-Z][^A-Z]+)")

func getFieldValue(f *reflect.StructField) (string, error) {
	var envName string

	// HogeFugaJSONValueをHOGE_FUGA_JSON_VALUEにする
	if matchedWordsList := gatherRegexp.FindAllStringSubmatch(f.Name, -1); len(matchedWordsList) > 0 {
		var names []string
		for _, words := range matchedWordsList {
			if m := acronymRegexp.FindStringSubmatch(words[0]); len(m) == 3 {
				names = append(names, m[1], m[2])
			} else {
				names = append(names, words[0])
			}
		}

		envName = strings.Join(names, "_") // [Hoge, Fuga] -> Hoge_Fuga
		envName = strings.ToUpper(envName) // Hoge_Fuga -> HOGE_FUGA
	}
	envName = strings.ReplaceAll(envName, "SQLDB", "SQL_DB")   // MySQLDB -> MY_SQLDB -> MY_SQL_DB
	envName = strings.ReplaceAll(envName, "MY_SQL", "MYSQL")   // MySQL -> MY_SQL -> MYSQL
	envName = strings.ReplaceAll(envName, "O_AUTH2", "OAUTH2") // OAuth2 -> O_AUTH2 -> OAUTH2
	envName = strings.ReplaceAll(envName, "QUA_PAY", "QUAPAY") // QuaPay -> QUA_PAY -> QUAPAY

	envValue, ok := syscall.Getenv(envName) // os.Getenv は「環境変数が未設定」であることと「環境変数の値が空文字である」ことを区別できないので syscall.Getenv を使う

	defaultValue := f.Tag.Get("default")
	if defaultValue != "" && !ok { // デフォ値が設定済みかつ環境変数の設定もない
		envValue = defaultValue
	}

	requiredValue := f.Tag.Get("required")
	yamlValue := f.Tag.Get("yaml")
	if defaultValue == "" && !ok { // デフォ値が未設定かつ環境変数の設定もない
		if requiredValue == "true" && yamlValue == "" { // `required:"false"`を指定されない限り環境変数の設定は必須だがyamlが設定されているならSecretManagerから取得するので無視
			return "", cerrors.Newf(cerrors.Internal, "環境変数の設定がされていません。 EnvName = %v", envName)
		}
	}

	return envValue, nil
}

// Process 渡された構造体のフィールド名を元に環境変数から値を取得してそのフィールドに書き込む
func Process(config interface{}) error {
	rv := reflect.ValueOf(config)

	if rv.Kind() != reflect.Ptr {
		return cerrors.Newf(cerrors.Internal, "引数にはポインタ型の構造体を渡してください。 Kind = %v", rv.Kind().String())
	}
	rv = rv.Elem() // ポインタ型の値型を抽出する
	if rv.Kind() != reflect.Struct {
		return cerrors.Newf(cerrors.Internal, "引数にはポインタ型の構造体を渡してください。 Kind = %v", rv.Kind().String())
	}

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		structField := rv.Type().Field(i)

		if field.Type().Name() == "Time" {
			fv, err := getFieldValue(&structField)
			if err != nil {
				return cerrors.Stack(err)
			}
			t, err := time.ParseInLocation("2006-01-02T15:04:05-0700", fv, time.Local)
			if err != nil { // パースに失敗したら現在時刻を入れる
				t = time.Now()
			}
			field.Set(reflect.ValueOf(t))
		} else {
			switch structField.Type.Kind() {
			case reflect.String:
				fv, err := getFieldValue(&structField)
				if err != nil {
					return cerrors.Stack(err)
				}
				field.SetString(fv)
			case reflect.Bool:
				fv, err := getFieldValue(&structField)
				if err != nil {
					return cerrors.Stack(err)
				}
				v, _ := strconv.ParseBool(fv) // フォーマットのエラーは結局値がfalseになるので無視
				field.SetBool(v)
			case reflect.Int32:
				fv, err := getFieldValue(&structField)
				if err != nil {
					return cerrors.Stack(err)
				}
				v, _ := strconv.ParseInt(fv, 10, 32) // フォーマットのエラーは結局値が0になるので無視
				field.SetInt(v)
			case reflect.Int64:
				fv, err := getFieldValue(&structField)
				if err != nil {
					return cerrors.Stack(err)
				}
				v, _ := strconv.ParseInt(fv, 10, 64) // フォーマットのエラーは結局値が0になるので無視
				field.SetInt(v)
			case reflect.Float64:
				fv, err := getFieldValue(&structField)
				if err != nil {
					return cerrors.Stack(err)
				}
				v, _ := strconv.ParseFloat(fv, 64) // フォーマットのエラーは結局値が0になるので無視
				field.SetFloat(v)
			case reflect.Slice:
				fv, err := getFieldValue(&structField)
				if err != nil {
					return cerrors.Stack(err)
				}
				strs := cstrings.SplitComma(fv)
				switch field.Type().Elem().Kind() {
				case reflect.String:
					field.Set(reflect.ValueOf(strs))
				case reflect.Bool:
					s := make([]bool, 0, len(strs))
					for _, e := range strs {
						v, _ := strconv.ParseBool(e)
						s = append(s, v)
					}
					field.Set(reflect.ValueOf(s))
				case reflect.Int32:
					s := make([]int32, 0, len(strs))
					for _, e := range strs {
						v, _ := strconv.ParseInt(e, 10, 32)
						s = append(s, int32(v))
					}
					field.Set(reflect.ValueOf(s))
				case reflect.Int64:
					s := make([]int64, 0, len(strs))
					for _, e := range strs {
						v, _ := strconv.ParseInt(e, 10, 64)
						s = append(s, v)
					}
					field.Set(reflect.ValueOf(s))
				case reflect.Array, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Float32, reflect.Float64,
					reflect.Func, reflect.Int, reflect.Int16, reflect.Int8,
					reflect.Interface, reflect.Invalid, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct,
					reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8,
					reflect.Uintptr, reflect.UnsafePointer:
				default:
					return cerrors.Newf(cerrors.Internal, "サポートされていない型です。 Field = %s, Kind = %s slice", structField.Name, structField.Type.Elem().Kind().String())
				}
			case reflect.Array, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Float32,
				reflect.Func, reflect.Int, reflect.Int16, reflect.Int8,
				reflect.Interface, reflect.Invalid, reflect.Map, reflect.Ptr, reflect.Struct,
				reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8,
				reflect.Uintptr, reflect.UnsafePointer:
			default:
				return cerrors.Newf(cerrors.Internal, "サポートされていない型です。 Field = %s, Kind = %s", structField.Name, structField.Type.Kind().String())
			}
		}
	}

	return nil
}
