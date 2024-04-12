{{ template "autogen_comment" }}
package validator

import (
	"context"
	"fmt"
	"strconv"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"github.com/xhayamix/proto-gen-golang/pkg/cerrors"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/component/masterdata"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/entity/master"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/enum"
	"github.com/xhayamix/proto-gen-golang/pkg/logging/app"
	"github.com/xhayamix/proto-gen-golang/pkg/util/strings"
)

type {{ .CamelName }}Validator struct{}

func new{{ .PascalName }}Validator() Validator {
	return &{{ .CamelName }}Validator{}
}

func (v *{{ .CamelName }}Validator) Register(validate *validator.Validate, uTranslator *ut.UniversalTranslator) error {
	validate.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		value := sl.Current().Interface().(master.{{ .PascalName }})
		switch value.SettingType {
		{{ $name := .PascalName -}}
		{{ range .Elements -}}
		case enum.{{ $name }}Type_{{ .PascalName }}:
			v.report{{ .PascalSettingType }}(sl, value.Value)
		{{ end -}}
		}
		v.extraValidate(ctx, sl, value)
	}, master.{{ .PascalName }}{})

	// translator
	if trans, ok := uTranslator.GetTranslator("en"); ok {
		if err := validate.RegisterTranslation("{{ .Name }}", trans, func(u ut.Translator) error {
			if err := u.Add("{{ .Name }}-type", "'{0}' must be a valid {1}", false); err != nil {
				return err
			}
			if err := u.Add("{{ .Name }}-format", "'{0}' dose not match the '{1}' format", false); err != nil {
				return err
			}
			return nil
		}, v.translationFn); err != nil {
			return cerrors.Wrap(err, cerrors.Internal)
		}
	}
	if trans, ok := uTranslator.GetTranslator("ja"); ok {
		if err := validate.RegisterTranslation("{{ .Name }}", trans, func(u ut.Translator) error {
			if err := u.Add("{{ .Name }}-type", "'{0}'は正しい{1}でなければなりません", false); err != nil {
				return err
			}
			if err := u.Add("{{ .Name }}-format", "'{0}'はフォーマット'{1}'と一致しません", false); err != nil {
				return err
			}
			return nil
		}, v.translationFn); err != nil {
			return cerrors.Wrap(err, cerrors.Internal)
		}
	}
	return nil
}

func (v *{{ .CamelName }}Validator) reportBool(sl validator.StructLevel, field string) {
	if _, err := strconv.ParseBool(field); err != nil {
		sl.ReportError(field, master.SettingColumnName.Value, master.SettingColumnName.Value, "setting", "type$bool")
	}
}

func (v *{{ .CamelName }}Validator) reportBoolSlice(sl validator.StructLevel, field string) {
	if _, err := strings.SplitCommaToBool(field); err != nil {
		sl.ReportError(field, master.SettingColumnName.Value, master.SettingColumnName.Value, "setting", "type$[]bool")
	}
}

func (v *{{ .CamelName }}Validator) reportInt32(sl validator.StructLevel, field string) {
	if _, err := strconv.ParseInt(field, 10, 32); err != nil {
		sl.ReportError(field, master.{{ .PascalName }}ColumnName.Value, master.{{ .PascalName }}ColumnName.Value, "{{ .Name }}", "type$int32")
		return
	}
}

func (v *{{ .CamelName }}Validator) reportInt64(sl validator.StructLevel, field string) {
	if _, err := strconv.ParseInt(field, 10, 64); err != nil {
		sl.ReportError(field, master.{{ .PascalName }}ColumnName.Value, master.{{ .PascalName }}ColumnName.Value, "{{ .Name }}", "type$int64")
		return
	}
}

func (v *{{ .CamelName }}Validator) reportInt32Slice(sl validator.StructLevel, field string) {
	if _, err := strings.SplitCommaToInt32(field); err != nil {
		sl.ReportError(field, master.{{ .PascalName }}ColumnName.Value, master.{{ .PascalName }}ColumnName.Value, "{{ .Name }}", "type$[]int32")
	}
}

func (v *{{ .CamelName }}Validator) reportInt64Slice(sl validator.StructLevel, field string) {
	if _, err := strings.SplitCommaToInt64(field); err != nil {
		sl.ReportError(field, master.{{ .PascalName }}ColumnName.Value, master.{{ .PascalName }}ColumnName.Value, "{{ .Name }}", "type$[]int64")
	}
}

func (v *{{ .CamelName }}Validator) reportString(sl validator.StructLevel, field string) {}

func (v *{{ .CamelName }}Validator) reportStringSlice(sl validator.StructLevel, field string) {}

func (v *{{ .CamelName }}Validator) reportNotDefined(sl validator.StructLevel, field string) {
	sl.ReportError(field, master.{{ .PascalName }}ColumnName.Value, master.{{ .PascalName }}ColumnName.Value, "enum", "enum.{{ .PascalName }}Type")
}

func (v *{{ .CamelName }}Validator) translationFn(u ut.Translator, fe validator.FieldError) string {
	params := strings.SplitN(fe.Param(), "$", 2)
	t, err := u.T("{{ .Name }}-"+params[0], fe.Field(), params[1])
	if err != nil {
		app.GetLogger().Warn(context.Background(), fmt.Sprintf("warning: error translating FieldError: %+v", fe))
		return fe.Error()
	}
	return t
}

type {{ .CamelName }}CustomValidator struct{}

func new{{ .PascalName }}CustomValidator() CustomValidator {
	return &{{ .CamelName }}CustomValidator{}
}

func (v *{{ .CamelName }}CustomValidator) Validate(ctx context.Context, masterCache masterdata.Cache, locale string) (Results, error) {
	// ソースコードのconstantをクライアントに渡しているので設定不要
	ignoreSettingTypeSet := enum.{{ .PascalName }}TypeSlice{
	{{ range .Elements }}{{ if .IsServerConstant -}}
		enum.{{ $.PascalName }}Type_{{ .PascalName }},
	{{ end }}{{ end -}}
	}.ToSet()

	results := make(Results, 0, len(enum.{{ .PascalName }}TypeValues))
	{{if len .PascalPrefix -}}for id := range masterCache.Get{{ .PascalName }}MapByPK() {{ "{" }}{{ end }}
	for _, typ := range enum.{{ .PascalName }}TypeValues {
		if ignoreSettingTypeSet.Has(typ) {
			continue
		}

		if !masterCache.Get{{ .PascalName }}MapByPK().Has({{if len .PascalPrefix }}id, {{ end }}typ) {
			results = append(results, &Result{
				Table: master.{{ .PascalName }}TableName,
				Key: (&master.{{ .PascalName }}PK{
					{{- if len .PascalPrefix }}
					ID: id,
					{{- end }}
					SettingType: typ,
				}).Key(),
				Tag:   ResultTagTypeRequired,
			})
		}
	}
	{{ if len .PascalPrefix }}{{ "}" }}{{ end }}

	rs, err := v.extraValidate(ctx, masterCache, locale)
	if err != nil {
		return nil, cerrors.Stack(err)
	}
	results = append(results, rs...)

	return results, nil
}
