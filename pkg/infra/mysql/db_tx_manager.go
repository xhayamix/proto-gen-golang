package mysql

import (
	"context"
	"database/sql"
	"log"

	"github.com/xhayamix/proto-gen-golang/pkg/cerrors"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/database"
)

type dbTxManager struct {
	db *sql.DB
}

func NewDBTxManager(db *sql.DB) database.TxManager {
	return &dbTxManager{db}
}

func (m *dbTxManager) Transaction(ctx context.Context, f func(ctx context.Context, tx database.RWTx) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return cerrors.Wrap(err, cerrors.Internal)
	}
	defer func() {
		// panic
		if p := recover(); p != nil {
			e := tx.Rollback()
			if e != nil {
				log.Fatal("failed to MySQL Rollback")
			}
			panic(p) // re-throw panic after Rollback
		}
		// error
		if err != nil {
			e := tx.Rollback()
			if e != nil {
				log.Fatal("failed to MySQL Rollback")
			}
			return
		}
		// 正常
		e := tx.Commit()
		if e != nil {
			log.Fatal("failed to MySQL Commit")
		}
	}()
	err = f(ctx, &dbTx{tx})
	if err != nil {
		return cerrors.Stack(err)
	}
	return nil
}

func (m *dbTxManager) ReadOnlyTransaction(ctx context.Context, f func(ctx context.Context, tx database.ROTx) error) error {
	// not support
	return cerrors.New(cerrors.Internal)
}

func (m *dbTxManager) BatchReadOnlyTransaction(ctx context.Context, f func(ctx context.Context, batchROTx database.BatchROTx) error) error {
	// not support
	return cerrors.New(cerrors.Internal)
}

type dbTx struct {
	value *sql.Tx
}

func (tx *dbTx) GetTxObject() any {
	return tx.value
}

func (tx *dbTx) ReadOnlyImpl() {}

func (tx *dbTx) ReadWriteImpl() {}

func ExtractTx(tx database.ROTx) (*sql.Tx, error) {
	txObject, ok := tx.GetTxObject().(*sql.Tx)
	if !ok {
		return nil, cerrors.New(cerrors.Internal)
	}
	return txObject, nil
}
