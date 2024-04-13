package envconfig

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/xhayamix/proto-gen-golang/pkg/cerrors"
	testutil "github.com/xhayamix/proto-gen-golang/pkg/test/util"
)

func TestProcess(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		type Config struct {
			// 型のテスト
			String      string
			Bool        bool
			Int32       int32
			Int64       int64
			Float64     float64
			Time        time.Time
			StringSlice []string
			BoolSlice   []bool
			Int32Slice  []int32
			Int64Slice  []int32

			// 名前変換のテスト
			Hoge              bool
			HogeFuga          bool
			HogeFugaJSON      bool
			HogeFugaJSONValue bool
			// 例外規則
			HogeMySQLDBFuga         bool
			HogeMySQLFuga           bool
			HogeOAuth2Fuga          bool
			HogeQuaPaySecretVersion string

			// オプションのテスト
			Piyo1 bool `default:"true"`
			Piyo2 bool `required:"true"`
			Piyo3 bool `yaml:"Piyo3" required:"true"`

			Ignore bool // requiredじゃないなら無視する
		}
		t.Setenv("STRING", "string value")
		t.Setenv("BOOL", "true")
		t.Setenv("INT32", "1")
		t.Setenv("INT64", "10")
		t.Setenv("FLOAT64", "1.1")
		t.Setenv("TIME", "2000-01-23T04:05:06"+time.Now().In(time.Local).Format("-0700"))
		t.Setenv("STRING_SLICE", "a,b,c")
		t.Setenv("BOOL_SLICE", "true,false,true")
		t.Setenv("INT32_SLICE", "1,10,100")
		t.Setenv("INT64_SLICE", "10,100,1000")

		t.Setenv("HOGE", "true")
		t.Setenv("HOGE_FUGA", "true")
		t.Setenv("HOGE_FUGA_JSON", "true")
		t.Setenv("HOGE_FUGA_JSON_VALUE", "true")

		t.Setenv("HOGE_MYSQL_DB_FUGA", "true")
		t.Setenv("HOGE_MYSQL_FUGA", "true")
		t.Setenv("HOGE_OAUTH2_FUGA", "true")

		t.Setenv("PIYO2", "true")
		t.Setenv("PIYO3", "true")
		t.Setenv("HOGE_QUAPAY_SECRET_VERSION", "latest")

		input := Config{}
		err := Process(&input)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, Config{
			String:      "string value",
			Bool:        true,
			Int32:       1,
			Int64:       10,
			Float64:     1.1,
			Time:        time.Date(2000, 1, 23, 4, 5, 6, 0, time.Local),
			StringSlice: []string{"a", "b", "c"},
			BoolSlice:   []bool{true, false, true},
			Int32Slice:  []int32{1, 10, 100},
			Int64Slice:  []int32{10, 100, 1000},

			Hoge:              true,
			HogeFuga:          true,
			HogeFugaJSON:      true,
			HogeFugaJSONValue: true,

			HogeMySQLDBFuga:         true,
			HogeMySQLFuga:           true,
			HogeOAuth2Fuga:          true,
			HogeQuaPaySecretVersion: "latest",

			Piyo1: true,
			Piyo2: true,
			Piyo3: true,

			Ignore: false,
		}, input)
	})

	t.Run("異常系", func(t *testing.T) {
		testutil.EqualCampusErrorByParam(t,
			Process("hoge"), cerrors.Internal,
			"引数にはポインタ型の構造体を渡してください。 Kind = string",
		)
		testutil.EqualCampusErrorByParam(t,
			Process(struct{}{}), cerrors.Internal,
			"引数にはポインタ型の構造体を渡してください。 Kind = struct",
		)
		pointerPointerStruct := &struct{}{}
		testutil.EqualCampusErrorByParam(t,
			Process(&pointerPointerStruct), cerrors.Internal,
			"引数にはポインタ型の構造体を渡してください。 Kind = ptr",
		)
		testutil.EqualCampusErrorByParam(t, Process(&struct {
			Hoge string `required:"true"`
		}{}),
			cerrors.Internal,
			"環境変数の設定がされていません。 EnvName = HOGE",
		)
	})
}
