package mock_database //nolint:golint,stylecheck // mockに実装を追加するためにアンダーバーを含むパッケージ名にする必要がある

import (
	"context"

	"go.uber.org/mock/gomock"

	"github.com/xhayamix/proto-gen-golang/pkg/domain/database"
)

func (m *MockTxManager) EXPECTReadOnlyTransaction(ctx context.Context, tx database.ROTx, times int) {
	m.EXPECT().ReadOnlyTransaction(ctx, gomock.Any()).Times(times).DoAndReturn(
		func(ctx context.Context, f func(context.Context, database.ROTx) error) error {
			return f(ctx, tx)
		})
}

func (m *MockTxManager) EXPECTBatchReadOnlyTransaction(ctx context.Context, tx database.BatchROTx, times int) {
	m.EXPECT().BatchReadOnlyTransaction(ctx, gomock.Any()).Times(times).DoAndReturn(
		func(ctx context.Context, f func(context.Context, database.BatchROTx) error) error {
			return f(ctx, tx)
		})
}

func (m *MockTxManager) EXPECTTransaction(ctx context.Context, tx database.RWTx, times int) {
	m.EXPECT().Transaction(ctx, gomock.Any()).Times(times).DoAndReturn(
		func(ctx context.Context, f func(context.Context, database.RWTx) error) error {
			return f(ctx, tx)
		})
}
