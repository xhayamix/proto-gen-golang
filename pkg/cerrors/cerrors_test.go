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
}

func TestAs(t *testing.T) {
	hayamiError, ok := As(errors.New("error"))
	assert.False(t, ok)
	assert.Nil(t, hayamiError)

	hayamiError, ok = As(Newf(Internal, "cerror"))
	assert.True(t, ok)
	assert.IsType(t, &HayamiError{}, hayamiError)
}
