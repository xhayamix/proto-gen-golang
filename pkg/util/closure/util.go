package closure

import (
	"log"
)

// QuietClose エラーを伝播させずClose処理を実行する.
func QuietClose(f func() error) {
	err := f()
	if err != nil {
		// エラーのないデフォルトのloggerで出力する。
		log.Printf("終了処理でエラーが発生しました. %v\n", err)
	}
}

type CloseListener struct {
	fs []func()
}

func (cl *CloseListener) Add(f func()) {
	cl.fs = append(cl.fs, f)
}

func (cl *CloseListener) Merge(closeListener *CloseListener) {
	for _, function := range closeListener.fs {
		cl.Add(function)
	}
}

func (cl *CloseListener) Close() {
	// 追加された逆順でCloseする
	for i := len(cl.fs) - 1; i >= 0; i-- {
		f := cl.fs[i]
		f()
	}
	cl.fs = nil
}
