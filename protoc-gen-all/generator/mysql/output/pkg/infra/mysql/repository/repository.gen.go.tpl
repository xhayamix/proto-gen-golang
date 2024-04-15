{{ template "autogen_comment" }}
{{- $name := .GoName -}}
{{- $camelName := .CamelName -}}
{{- $pkColumns := .PKColumns -}}
{{- $pkgName := .PkgName -}}
{{- $indexMethods := list -}}
{{- $tableName := .SnakeName }}

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/xhayamix/proto-gen-golang/pkg/cerrors"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/database"
	mysqlentity "github.com/xhayamix/proto-gen-golang/pkg/domain/entity/mysql"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/enum"
	repo "github.com/xhayamix/proto-gen-golang/pkg/domain/repository/mysql"
	"github.com/xhayamix/proto-gen-golang/pkg/infra/mysql"
)

type {{ .CamelName }}Repository struct {
	db *sql.DB
}

func New{{ .GoName }}Repository(db mysql.MysqlDB) repo.{{ .GoName }}Repository {
	return &{{ .CamelName }}Repository{
		db: db,
	}
}

func (r *{{ .CamelName }}Repository) SelectAll(ctx context.Context) (mysqlentity.{{ .GoName }}Slice, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM `{{ $tableName }}`")
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	slice := make(mysqlentity.{{ .GoName }}Slice, 0)
	for rows.Next() {
		entity := &mysqlentity.{{ .GoName }}{}
		ptrs := entity.PtrFromMapping(cols)
		if err := rows.Scan(ptrs...); err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
		slice = append(slice, entity)
	}
	return slice, nil
}

func (r *{{ .CamelName }}Repository) SelectAllByTx(ctx context.Context, _tx database.ROTx) (mysqlentity.{{ .GoName }}Slice, error) {
	tx, err := mysql.ExtractTx(_tx)
	if err != nil {
		return nil, cerrors.Stack(err)
	}
	rows, err := tx.QueryContext(ctx, "SELECT * FROM `{{ $tableName }}`")
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	slice := make(mysqlentity.{{ .GoName }}Slice, 0)
	for rows.Next() {
		entity := &mysqlentity.{{ .GoName }}{}
		ptrs := entity.PtrFromMapping(cols)
		if err := rows.Scan(ptrs...); err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
		slice = append(slice, entity)
	}
	return slice, nil
}

func (r *{{ .CamelName }}Repository) SelectByPKs(ctx context.Context, pks mysqlentity.{{ .GoName }}PKs) (mysqlentity.{{ .GoName }}Slice, error) {
    var entities mysqlentity.{{ .GoName }}Slice
    for _, pk := range pks {
        entity, err := r.SelectByPK(ctx{{ range $pk := .PKColumns }}, pk.{{ .GoName }}{{ end }})
        if err != nil {
            return nil, cerrors.Wrap(err, cerrors.Internal)
        }
        if entity != nil {
            entities = append(entities, entity)
        }
    }
    return entities, nil
}

func (r *{{ .CamelName }}Repository) SelectByPK(ctx context.Context{{ range $pk := .PKColumns }}, {{ .GoName }}_ {{ .Type }}{{ end }}) (*mysqlentity.{{ .GoName }}, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM `{{ $tableName }}` WHERE
	  {{- range $i, $pk := .PKColumns -}}
	    {{ if $i }} AND{{ end }} `{{ $pk.SnakeName }}`=?
	  {{- end -}}
	"
	  {{- range $pk := .PKColumns -}}
	    , {{ .GoName }}_
	  {{- end -}}
	)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	entity := &mysqlentity.{{ .GoName }}{}
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

func (r *{{ .CamelName }}Repository) SelectByTx(ctx context.Context, _tx database.ROTx{{ range $pk := .PKColumns }}, {{ .GoName }}_ {{ .Type }}{{ end }}) (*mysqlentity.{{ .GoName }}, error) {
	tx, err := mysql.ExtractTx(_tx)
	if err != nil {
		return nil, cerrors.Stack(err)
	}
	rows, err := tx.QueryContext(ctx, "SELECT * FROM `{{ $tableName }}` WHERE
	  {{- range $i, $pk := .PKColumns -}}
	    {{ if $i }} AND{{ end }} `{{ $pk.SnakeName }}`=?
	  {{- end -}}
	"
	  {{- range $pk := .PKColumns -}}
	    , {{ .GoName }}_
	  {{- end -}}
	)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	entity := &mysqlentity.{{ .GoName }}{}
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

{{ range $i, $_ := slice .PKColumns 0 (sub (len .PKColumns) 1) -}}
	{{ $cols := slice $pkColumns 0 (add1 $i) -}}
func (r *{{ $camelName }}Repository) SelectBy{{ range $j, $col := $cols }}{{ if $j }}And{{ end }}{{ $col.GoName }}{{ end -}}
  (ctx context.Context{{ range $cols }}, {{ .GoName }} {{ .Type }}{{ end }}) (mysqlentity.{{ $name }}Slice, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM `{{ $tableName }}` WHERE
	  {{- range $i, $col := $cols -}}
	    {{ if $i }} AND{{ end }} `{{ $col.SnakeName }}`=?
	  {{- end -}}
	"
	  {{- range $cols -}}
	    , {{ .GoName }}
	  {{- end -}}
	)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	slice := make(mysqlentity.{{ $name }}Slice, 0)
	for rows.Next() {
		entity := &mysqlentity.{{ $name }}{}
		ptrs := entity.PtrFromMapping(cols)
		if err := rows.Scan(ptrs...); err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
		slice = append(slice, entity)
	}
	return slice, nil
}
{{ end }}

{{ range .Indexes -}}
  {{ $keys := .Keys -}}
	{{ range $i, $_ := .Keys -}}
		{{ $cols := slice $keys 0 (add1 $i) -}}
			{{ $colNames := list -}}
			{{ range $cols }}{{ $colNames = append $colNames .GoName }}{{ end -}}
			{{ $method := $colNames | join "And" -}}
			{{ if not (has $method $indexMethods) -}}
				{{ $indexMethods = append $indexMethods $method -}}
func (r *{{ $camelName }}Repository) SelectBy{{ $method -}}
  (ctx context.Context{{ range $cols }}, {{ .GoName }} {{ .Type }}{{ end }}) (mysqlentity.{{ $name }}Slice, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM `{{ $tableName }}` WHERE
	  {{- range $i, $col := $cols -}}
	    {{ if $i }} AND{{ end }} `{{ $col.SnakeName }}`=?
	  {{- end -}}
	"
	  {{- range $cols -}}
	    , {{ .GoName }}
	  {{- end -}}
	)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	slice := make(mysqlentity.{{ $name }}Slice, 0)
	for rows.Next() {
		entity := &mysqlentity.{{ $name }}{}
		ptrs := entity.PtrFromMapping(cols)
		if err := rows.Scan(ptrs...); err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
		slice = append(slice, entity)
	}
	return slice, nil
}

{{ end }}
{{ end }}
{{ end }}

{{ range .PKColumns -}}
func (r *{{ $camelName }}Repository) SearchBy{{ .GoName }}(ctx context.Context, searchText string, limit int) ([]{{ .Type }}, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT `{{ .SnakeName }}` FROM `{{ $tableName }}` WHERE `{{ .SnakeName }}` LIKE ? LIMIT ?", fmt.Sprintf("%%%s%%", searchText), limit)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	defer rows.Close()

	slice := make([]{{ .Type }}, 0)
	for rows.Next() {
		var col {{ .Type }}
		if err := rows.Scan(&col); err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
		slice = append(slice, col)
	}
	return slice, nil
}
{{ end }}
func (r *{{ .CamelName }}Repository) Insert(ctx context.Context, _tx database.RWTx, entity *mysqlentity.{{ .GoName }}) error {
	tx, err := mysql.ExtractTx(_tx)
	if err != nil {
		return cerrors.Stack(err)
	}
	query := "INSERT INTO `{{ $tableName }}` (
	  {{- range $i, $column := .Columns -}}
		  {{- if $i }},{{ end -}}
		  `{{ $column.SnakeName }}`
	  {{- end -}}
	) VALUES (
	  {{- range $i, $column := .Columns -}}
		  {{- if $i -}}, {{- end -}}
		  ?
	  {{- end -}}
	)"
	vals := entity.PtrFromMapping(mysqlentity.{{ .GoName }}Cols)
	_, err = tx.ExecContext(ctx, query, vals...)
	if err != nil {
		return cerrors.Wrap(err, cerrors.Internal)
	}
	return nil
}

func (r *{{ .CamelName }}Repository) BulkInsert(ctx context.Context, _tx database.RWTx, entities mysqlentity.{{ .GoName }}Slice, replace bool) error {
	tx, err := mysql.ExtractTx(_tx)
	if err != nil {
		return cerrors.Stack(err)
	}

	columnNames := mysqlentity.{{ .GoName }}Cols
	recordCount := len(entities)
	columnCount := len(columnNames)

	var sqlColumns []string
	var paramPlaces []string
	for _, columnName := range columnNames {
		sqlColumns = append(sqlColumns, fmt.Sprintf("`%s`", columnName))
		paramPlaces = append(paramPlaces, "?")
	}

	recordsList := make([]mysqlentity.{{ .GoName }}Slice, 0, int64(math.Ceil(float64(recordCount)/float64(mysql.PlaceholderLimit))))
	tempRecordsCap := recordCount
	if recordCount > mysql.PlaceholderLimit {
		tempRecordsCap = mysql.PlaceholderLimit
	}
	tempRecords := make(mysqlentity.{{ .GoName }}Slice, 0, tempRecordsCap)
	placeholderCount := 0
	// プリペアドステートメントの上限に合わせてクエリを分割する
	for i, record := range entities {
		placeholderCount += columnCount

		if placeholderCount > mysql.PlaceholderLimit { // 上限を超えた場合
			placeholderCount = columnCount
			recordsList = append(recordsList, tempRecords)
			tempRecords = make(mysqlentity.{{ .GoName }}Slice, 0, mysql.PlaceholderLimit) // 上限を超えないかも知れないが一応上限までCapを確保して初期化
		}

		tempRecords = append(tempRecords, record)
		if i+1 >= recordCount { // 最後のループの場合
			recordsList = append(recordsList, tempRecords)
		}
	}

	var baseQuery string
	if replace {
		baseQuery = fmt.Sprintf("REPLACE INTO `{{ .SnakeName }}` (%s) VALUES ", strings.Join(sqlColumns, ","))
	} else {
		baseQuery = fmt.Sprintf("INSERT INTO `{{ .SnakeName }}` (%s) VALUES ", strings.Join(sqlColumns, ","))
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
			return cerrors.Wrapf(err, cerrors.Internal, "書き込み中にエラーが発生しました。 tableName = {{ .SnakeName }}")
		}
	}

	return nil
}

func (r *{{ .CamelName }}Repository) Update(ctx context.Context, _tx database.RWTx, entity *mysqlentity.{{ .GoName }}) error {
	tx, err := mysql.ExtractTx(_tx)
	if err != nil {
		return cerrors.Stack(err)
	}
	query := "UPDATE `{{ $tableName }}` SET
		{{- range $i, $column := .Columns -}}
			{{- if $i }},{{ else }} {{ end -}}
			`{{ $column.SnakeName }}`=?
		{{- end }} WHERE
		{{- range $i, $pk := .PKColumns -}}
			{{ if $i }} AND{{ end }} `{{ $pk.SnakeName }}`=?
		{{- end -}}
	"
	vals := entity.PtrFromMapping([]string{
	  {{- range $i, $column := .Columns -}}
		  {{- if $i }}, {{ end -}}
		  "{{ $column.SnakeName }}"
	  {{- end -}}
	  {{- range $pk := .PKColumns -}}
	    , "{{ $pk.SnakeName }}"
	  {{- end -}}
	})
	_, err = tx.ExecContext(ctx, query, vals...)
	if err != nil {
		return cerrors.Wrap(err, cerrors.Internal)
	}
	return nil
}

func (r *{{ .CamelName }}Repository) Delete(ctx context.Context, _tx database.RWTx, entity *mysqlentity.{{ .GoName }}) error {
	tx, err := mysql.ExtractTx(_tx)
	if err != nil {
		return cerrors.Stack(err)
	}
	query := "DELETE FROM `{{ $tableName }}` WHERE
	  {{- range $i, $pk := .PKColumns -}}
	    {{ if $i }} AND{{ end }} `{{ $pk.SnakeName }}`=?
	  {{- end -}}
	"
	vals := entity.PtrFromMapping([]string{
	  {{- range $i, $pk := .PKColumns -}}
		  {{- if $i }}, {{ end -}}
		  "{{ $pk.SnakeName }}"
	  {{- end -}}
	})
	_, err = tx.ExecContext(ctx, query, vals...)
	if err != nil {
		return cerrors.Wrap(err, cerrors.Internal)
	}
	return nil
}

func (r *{{ .CamelName }}Repository) BulkDelete(ctx context.Context, _tx database.RWTx, entities mysqlentity.{{ .GoName }}Slice) error {
	if len(entities) == 0 {
		return nil
	}

	tx, err := mysql.ExtractTx(_tx)
	if err != nil {
		return cerrors.Stack(err)
	}

	pkNames := []string{
	{{- range $i, $pk := .PKColumns -}}
		"{{ $pk.SnakeName }}",
	{{- end -}}
	}

	recordCount := len(entities)
	pkCount := len(pkNames)

	recordsList := make([]mysqlentity.{{ .GoName }}Slice, 0, int64(math.Ceil(float64(recordCount)/float64(mysql.PlaceholderLimit))))
	tempRecordsCap := recordCount
	if recordCount > mysql.PlaceholderLimit {
		tempRecordsCap = mysql.PlaceholderLimit
	}
	tempRecords := make(mysqlentity.{{ .GoName }}Slice, 0, tempRecordsCap)
	placeholderCount := 0
	// プリペアドステートメントの上限に合わせてクエリを分割する
	for i, record := range entities {
		placeholderCount += pkCount

		if placeholderCount > mysql.PlaceholderLimit { // 上限を超えた場合
			placeholderCount = pkCount
			recordsList = append(recordsList, tempRecords)
			tempRecords = make(mysqlentity.{{ .GoName }}Slice, 0, mysql.PlaceholderLimit) // 上限を超えないかも知れないが一応上限までCapを確保して初期化
		}

		tempRecords = append(tempRecords, record)
		if i+1 >= recordCount { // 最後のループの場合
			recordsList = append(recordsList, tempRecords)
		}
	}

	baseQuery := "DELETE FROM `{{ $tableName }}` WHERE (
	{{- range $i, $pk := .PKColumns -}}
		{{ if $i }}, {{ end }}`{{ $pk.SnakeName }}`
	{{- end -}}
	) IN ("
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

			queryStringBuilder.WriteString("(
			{{- range $i, $pk := .PKColumns -}}
				{{ if $i }}, {{ end }}?
			{{- end -}}
			)")
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
			return cerrors.Wrapf(err, cerrors.Internal, "削除中にエラーが発生しました。 tableName = {{ .SnakeName }}")
		}
	}

	return nil
}

func (r *{{ .CamelName }}Repository) DeleteAll(ctx context.Context, _tx database.RWTx) error {
	tx, err := mysql.ExtractTx(_tx)
	if err != nil {
		return cerrors.Stack(err)
	}
	query := "DELETE FROM `{{ $tableName }}`"
	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		return cerrors.Wrap(err, cerrors.Internal)
	}
	return nil
}
