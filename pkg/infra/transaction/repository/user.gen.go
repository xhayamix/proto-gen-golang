// Code generated by protoc-gen-all. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/xhayamix/proto-gen-golang/pkg/cerrors"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/database"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/entity/transaction"
	repo "github.com/xhayamix/proto-gen-golang/pkg/domain/repository/transaction"
	"github.com/xhayamix/proto-gen-golang/pkg/infra/mysql"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db mysql.MysqlDB) repo.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) SelectAll(ctx context.Context) (transaction.UserSlice, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM `user`")
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	slice := make(transaction.UserSlice, 0)
	for rows.Next() {
		entity := &transaction.User{}
		ptrs := entity.PtrFromMapping(cols)
		if err := rows.Scan(ptrs...); err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
		slice = append(slice, entity)
	}
	return slice, nil
}

func (r *userRepository) SelectAllOffset(ctx context.Context, offset, limit int) (transaction.UserSlice, error) {
	query := "SELECT * FROM `user` ORDER BY `created_at` ASC LIMIT ? OFFSET ?"
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	slice := make(transaction.UserSlice, 0)
	for rows.Next() {
		entity := &transaction.User{}
		ptrs := entity.PtrFromMapping(cols)
		if err := rows.Scan(ptrs...); err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
		slice = append(slice, entity)
	}
	return slice, nil
}

func (r *userRepository) SelectAllByTx(ctx context.Context, _tx database.ROTx) (transaction.UserSlice, error) {
	tx, err := mysql.ExtractTx(_tx)
	if err != nil {
		return nil, cerrors.Stack(err)
	}
	rows, err := tx.QueryContext(ctx, "SELECT * FROM `user`")
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	slice := make(transaction.UserSlice, 0)
	for rows.Next() {
		entity := &transaction.User{}
		ptrs := entity.PtrFromMapping(cols)
		if err := rows.Scan(ptrs...); err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
		slice = append(slice, entity)
	}
	return slice, nil
}

func (r *userRepository) SelectByPKs(ctx context.Context, pks transaction.UserPKs) (transaction.UserSlice, error) {
	var entities transaction.UserSlice
	for _, pk := range pks {
		entity, err := r.SelectByPK(ctx, pk.ID)
		if err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
		if entity != nil {
			entities = append(entities, entity)
		}
	}
	return entities, nil
}

func (r *userRepository) SelectByPK(ctx context.Context, ID_ string) (*transaction.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM `user` WHERE `id`=?", ID_)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	entity := &transaction.User{}
	ptrs := entity.PtrFromMapping(cols)
	foundOne := false
	for rows.Next() {
		foundOne = true
		if err := rows.Scan(ptrs...); err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
	}
	if !foundOne {
		return nil, nil
	}
	return entity, nil
}

func (r *userRepository) SelectByTx(ctx context.Context, _tx database.ROTx, ID_ string) (*transaction.User, error) {
	tx, err := mysql.ExtractTx(_tx)
	if err != nil {
		return nil, cerrors.Stack(err)
	}
	rows, err := tx.QueryContext(ctx, "SELECT * FROM `user` WHERE `id`=?", ID_)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	entity := &transaction.User{}
	ptrs := entity.PtrFromMapping(cols)
	foundOne := false
	for rows.Next() {
		foundOne = true
		if err := rows.Scan(ptrs...); err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
	}
	if !foundOne {
		return nil, nil
	}
	return entity, nil
}

func (r *userRepository) SearchByID(ctx context.Context, searchText string, limit int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT `id` FROM `user` WHERE `id` LIKE ? LIMIT ?", fmt.Sprintf("%%%s%%", searchText), limit)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	defer rows.Close()

	slice := make([]string, 0)
	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
		slice = append(slice, col)
	}
	return slice, nil
}

func (r *userRepository) Insert(ctx context.Context, _tx database.RWTx, entity *transaction.User) error {
	tx, err := mysql.ExtractTx(_tx)
	if err != nil {
		return cerrors.Stack(err)
	}
	query := "INSERT INTO `user` (`id`,`user_id`,`email`,`password`,`name`,`profile`,`icon_img`,`header_img`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?,?)"
	vals := entity.PtrFromMapping(transaction.UserCols)
	_, err = tx.ExecContext(ctx, query, vals...)
	if err != nil {
		return cerrors.Wrap(err, cerrors.Internal)
	}
	return nil
}

func (r *userRepository) BulkInsert(ctx context.Context, _tx database.RWTx, entities transaction.UserSlice, replace bool) error {
	tx, err := mysql.ExtractTx(_tx)
	if err != nil {
		return cerrors.Stack(err)
	}

	columnNames := transaction.UserCols
	recordCount := len(entities)
	columnCount := len(columnNames)

	var sqlColumns []string
	var paramPlaces []string
	for _, columnName := range columnNames {
		sqlColumns = append(sqlColumns, fmt.Sprintf("`%s`", columnName))
		paramPlaces = append(paramPlaces, "?")
	}

	recordsList := make([]transaction.UserSlice, 0, int64(math.Ceil(float64(recordCount)/float64(mysql.PlaceholderLimit))))
	tempRecordsCap := recordCount
	if recordCount > mysql.PlaceholderLimit {
		tempRecordsCap = mysql.PlaceholderLimit
	}
	tempRecords := make(transaction.UserSlice, 0, tempRecordsCap)
	placeholderCount := 0
	// プリペアドステートメントの上限に合わせてクエリを分割する
	for i, record := range entities {
		placeholderCount += columnCount

		if placeholderCount > mysql.PlaceholderLimit { // 上限を超えた場合
			placeholderCount = columnCount
			recordsList = append(recordsList, tempRecords)
			tempRecords = make(transaction.UserSlice, 0, mysql.PlaceholderLimit) // 上限を超えないかも知れないが一応上限までCapを確保して初期化
		}

		tempRecords = append(tempRecords, record)
		if i+1 >= recordCount { // 最後のループの場合
			recordsList = append(recordsList, tempRecords)
		}
	}

	var baseQuery string
	if replace {
		baseQuery = fmt.Sprintf("REPLACE INTO `user` (%s) VALUES ", strings.Join(sqlColumns, ","))
	} else {
		baseQuery = fmt.Sprintf("INSERT INTO `user` (%s) VALUES ", strings.Join(sqlColumns, ","))
	}
	valuesPlaceholders := fmt.Sprintf("(%s)", strings.Join(paramPlaces, ","))
	queryStringBuilder := &strings.Builder{}
	for _, records := range recordsList {
		queryStringBuilder.Reset()
		queryStringBuilder.WriteString(baseQuery)

		beforeQueryLength := queryStringBuilder.Len()
		allRecordValues := make([]interface{}, 0, len(records)*columnCount)
		for i, record := range records {
			keyValueMap := record.ToKeyValue()
			for _, name := range columnNames {
				value, ok := keyValueMap[name]
				if !ok {
					return cerrors.Newf(cerrors.Internal, "カラム名が不正です。 name = %q", name)
				}

				allRecordValues = append(allRecordValues, value)
			}

			queryStringBuilder.WriteString(valuesPlaceholders)

			if i+1 >= len(records) {
				queryStringBuilder.WriteString(";")
			} else {
				queryStringBuilder.WriteString(",")
			}
		}

		// クエリが初期状態なら最後のクエリを投げた後だと判断し終了する
		if beforeQueryLength == queryStringBuilder.Len() {
			return nil
		}

		if _, err = tx.ExecContext(ctx, queryStringBuilder.String(), allRecordValues...); err != nil {
			return cerrors.Wrapf(err, cerrors.Internal, "書き込み中にエラーが発生しました。 tableName = user")
		}
	}

	return nil
}

func (r *userRepository) Update(ctx context.Context, _tx database.RWTx, entity *transaction.User) error {
	tx, err := mysql.ExtractTx(_tx)
	if err != nil {
		return cerrors.Stack(err)
	}
	query := "UPDATE `user` SET `id`=?,`user_id`=?,`email`=?,`password`=?,`name`=?,`profile`=?,`icon_img`=?,`header_img`=?,`created_at`=?,`updated_at`=? WHERE `id`=?"
	vals := entity.PtrFromMapping([]string{"id", "user_id", "email", "password", "name", "profile", "icon_img", "header_img", "created_at", "updated_at", "id"})
	_, err = tx.ExecContext(ctx, query, vals...)
	if err != nil {
		return cerrors.Wrap(err, cerrors.Internal)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, _tx database.RWTx, entity *transaction.User) error {
	tx, err := mysql.ExtractTx(_tx)
	if err != nil {
		return cerrors.Stack(err)
	}
	query := "DELETE FROM `user` WHERE `id`=?"
	vals := entity.PtrFromMapping([]string{"id"})
	_, err = tx.ExecContext(ctx, query, vals...)
	if err != nil {
		return cerrors.Wrap(err, cerrors.Internal)
	}
	return nil
}

func (r *userRepository) BulkDelete(ctx context.Context, _tx database.RWTx, entities transaction.UserSlice) error {
	if len(entities) == 0 {
		return nil
	}

	tx, err := mysql.ExtractTx(_tx)
	if err != nil {
		return cerrors.Stack(err)
	}

	pkNames := []string{"id"}

	recordCount := len(entities)
	pkCount := len(pkNames)

	recordsList := make([]transaction.UserSlice, 0, int64(math.Ceil(float64(recordCount)/float64(mysql.PlaceholderLimit))))
	tempRecordsCap := recordCount
	if recordCount > mysql.PlaceholderLimit {
		tempRecordsCap = mysql.PlaceholderLimit
	}
	tempRecords := make(transaction.UserSlice, 0, tempRecordsCap)
	placeholderCount := 0
	// プリペアドステートメントの上限に合わせてクエリを分割する
	for i, record := range entities {
		placeholderCount += pkCount

		if placeholderCount > mysql.PlaceholderLimit { // 上限を超えた場合
			placeholderCount = pkCount
			recordsList = append(recordsList, tempRecords)
			tempRecords = make(transaction.UserSlice, 0, mysql.PlaceholderLimit) // 上限を超えないかも知れないが一応上限までCapを確保して初期化
		}

		tempRecords = append(tempRecords, record)
		if i+1 >= recordCount { // 最後のループの場合
			recordsList = append(recordsList, tempRecords)
		}
	}

	baseQuery := "DELETE FROM `user` WHERE (`id`) IN ("
	queryStringBuilder := &strings.Builder{}
	for _, records := range recordsList {
		queryStringBuilder.Reset()
		queryStringBuilder.WriteString(baseQuery)

		beforeQueryLength := queryStringBuilder.Len()
		allRecordValues := make([]interface{}, 0, len(records)*pkCount)

		for i, record := range records {
			keyValueMap := record.ToKeyValue()
			for _, pkName := range pkNames {
				value, ok := keyValueMap[pkName]
				if !ok {
					return cerrors.Newf(cerrors.Internal, "カラム名が不正です。 pkName = %q", pkName)
				}

				allRecordValues = append(allRecordValues, value)
			}

			queryStringBuilder.WriteString("(?)")
			if i+1 >= len(records) {
				queryStringBuilder.WriteString(");")
			} else {
				queryStringBuilder.WriteString(",")
			}
		}

		// クエリが初期状態なら最後のクエリを投げた後だと判断し終了する
		if beforeQueryLength == queryStringBuilder.Len() {
			return nil
		}

		if _, err = tx.ExecContext(ctx, queryStringBuilder.String(), allRecordValues...); err != nil {
			return cerrors.Wrapf(err, cerrors.Internal, "削除中にエラーが発生しました。 tableName = user")
		}
	}

	return nil
}

func (r *userRepository) DeleteAll(ctx context.Context, _tx database.RWTx) error {
	tx, err := mysql.ExtractTx(_tx)
	if err != nil {
		return cerrors.Stack(err)
	}
	query := "DELETE FROM `user`"
	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		return cerrors.Wrap(err, cerrors.Internal)
	}
	return nil
}
