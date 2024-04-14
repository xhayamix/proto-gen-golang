//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_interactor.go
//go:generate goimports -w --local "github.com/xhayamix/proto-gen-golang" mock_$GOPACKAGE/mock_interactor.go
package api

import (
	"context"
	"encoding/json"
	"github.com/xhayamix/proto-gen-golang/pkg/cerrors"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/database"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/port/api"
)

type Interactor interface {
	List(ctx context.Context) []api.Method
	Request(ctx context.Context, method api.Method, authToken, param string) (string, error)
}

type interactor struct {
	env       string
	locations []string
	api       api.API
	txManager database.UserTxManager
}

func New(
	api api.API,
	txManager database.UserTxManager,
) Interactor {
	return &interactor{
		api:       api,
		txManager: txManager,
	}
}

func (i *interactor) List(ctx context.Context) []api.Method {
	return i.list(ctx)
}

func (i *interactor) Request(ctx context.Context, method api.Method, userID, param string) (string, error) {
	res, err := i.request(ctx, method, param)
	if err != nil {
		return "", cerrors.Stack(err)
	}
	resultByte, err := json.Marshal(res)
	if err != nil {
		return "", cerrors.Wrap(err, cerrors.Internal)
	}

	return string(resultByte), nil
}
