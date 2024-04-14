package cerrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	t.Skip("出力の見た目確認用")
	err := errors.New("new error")
	err = Wrapf(err, InvalidArgument, "wrapf")
	err = Stack(err)
	err = Wrap(err, Internal)
	fmt.Printf("%+v\n%v\n", err, err.Error())
	// -- 出力イメージ --
	// 	wrapf:
	// 		github.com/xhayamix/proto-gen-golang/pkg/cerrors.TestWrap
	// 			/campus-server/pkg/cerrors/cerrors_test.go:14
	// 	  - github.com/xhayamix/proto-gen-golang/pkg/cerrors.TestWrap
	// 			/campus-server/pkg/cerrors/cerrors_test.go:13
	// 	  - wrapf:
	// 		github.com/xhayamix/proto-gen-golang/pkg/cerrors.TestWrap
	// 			/campus-server/pkg/cerrors/cerrors_test.go:12
	// 	  - new error
	// 	error: code = Internal, message = wrapf
}

func TestAs(t *testing.T) {
	campusError, ok := As(errors.New("error"))
	assert.False(t, ok)
	assert.Nil(t, campusError)

	campusError, ok = As(Newf(Internal, "cerror"))
	assert.True(t, ok)
	assert.IsType(t, &CampusError{}, campusError)
}
