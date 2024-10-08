{{ template "autogen_comment" }}
//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_{{ .SnakeName }}.go
//go:generate goimports -w --local "github.com/xhayamix/proto-gen-golang" mock_$GOPACKAGE/mock_{{ .SnakeName }}.go
{{ $name := .GoName -}}
{{ $pkColumns := .PKColumns -}}
{{ $indexMethods := list -}}
{{ $pkgName := .PkgName }}
package transaction

import (
	"context"

	"github.com/xhayamix/proto-gen-golang/pkg/domain/database"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/entity/transaction"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/enum"
)

type {{ $name }}Repository interface {
	SelectAll(ctx context.Context) (transaction.{{ .GoName }}Slice, error)
	SelectAllOffset(ctx context.Context, offset, limit int) (transaction.{{ .GoName }}Slice, error)
	SelectAllByTx(ctx context.Context, tx database.ROTx) (transaction.{{ .GoName }}Slice, error)
	SelectByPK(ctx context.Context{{ range .PKColumns }}, {{ .GoName }}_ {{ .Type }}{{ end }}) (*transaction.{{ .GoName }}, error)
	SelectByTx(ctx context.Context, tx database.ROTx{{ range .PKColumns }}, {{ .GoName }}_ {{ .Type }}{{ end }}) (*transaction.{{ .GoName }}, error)
	SelectByPKs(ctx context.Context, pks transaction.{{ .GoName }}PKs) (transaction.{{ .GoName }}Slice, error)
	{{ range $i, $_ := slice .PKColumns 0 (sub (len .PKColumns) 1) -}}
		{{ $cols := slice $pkColumns 0 (add1 $i) -}}
		SelectBy{{ range $j, $col := $cols }}{{ if $j }}And{{ end }}{{ $col.GoName }}{{ end -}}
			(ctx context.Context{{ range $cols }}, {{ .GoName }} {{ .Type }}{{ end }}) (transaction.{{ $name }}Slice, error)
	{{ end -}}
	{{ range .Indexes -}}
		{{ $keys := .Keys -}}
		{{ range $i, $_ := .Keys -}}
			{{ $cols := slice $keys 0 (add1 $i) -}}
			{{ $colNames := list -}}
			{{ range $cols }}{{ $colNames = append $colNames .GoName }}{{ end -}}
			{{ $method := $colNames | join "And" -}}
			{{ if not (has $method $indexMethods) -}}
				{{ $indexMethods = append $indexMethods $method -}}
				SelectBy{{ $method -}}
					(ctx context.Context{{ range $cols }}, {{ .GoName }} {{ .Type }}{{ end }}) (transaction.{{ $name }}Slice, error)
			{{ end -}}
		{{ end -}}
	{{ end -}}
	{{ range slice .PKColumns -}}
		SearchBy{{ .GoName }}(ctx context.Context, searchText string, limit int) ([]{{ .Type }}, error)
	{{ end -}}
	Insert(ctx context.Context, tx database.RWTx, entity *transaction.{{ .GoName }}) error
	BulkInsert(ctx context.Context, tx database.RWTx, entities transaction.{{ .GoName }}Slice, replace bool) error
	Update(ctx context.Context, tx database.RWTx, entity *transaction.{{ .GoName }}) error
	Delete(ctx context.Context, tx database.RWTx, entity *transaction.{{ .GoName }}) error
	BulkDelete(ctx context.Context, tx database.RWTx, entities transaction.{{ .GoName }}Slice) error
	DeleteAll(ctx context.Context, tx database.RWTx) error
}
