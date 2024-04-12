package input

import (
	"fmt"
	"net/http"
	"strings"

	validateoptions "github.com/envoyproxy/protoc-gen-validate/validate"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/xhayamix/proto-gen-golang/pkg/domain/proto/server/options/api"
	zapoptions "github.com/xhayamix/proto-gen-golang/pkg/domain/proto/server/options/zap"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

func ConvertFileFromProto(protoFile *protogen.File) (*File, error) {
	if len(protoFile.Services) > 1 {
		return nil, perrors.Newf("このprotoファイルに2つ以上のメッセージを含めることは出来ません。 protoFile = %v", protoFile.Desc.FullName())
	}

	paths := strings.Split(protoFile.GeneratedFilenamePrefix, "/")
	file := &File{
		IsCommon:    string(protoFile.Desc.Name()) == "common",
		SnakeName:   core.ToSnakeCase(paths[len(paths)-1]),
		PackageName: "client.api", // string(protoFile.Desc.Package().Name()), client,serverでpackage名があってないので本来のやり方だとgameになってしまう
		Service:     nil,
		Messages:    make([]*Message, 0, len(protoFile.Messages)),
	}
	for _, message := range protoFile.Messages {
		m, err := createMessage(message)
		if err != nil {
			return nil, perrors.Stack(err)
		}

		file.Messages = append(file.Messages, m)
	}

	if len(protoFile.Services) == 0 {
		return file, nil
	}

	protoService := protoFile.Services[0]
	serviceOption, ok := proto.GetExtension(protoService.Desc.Options(), api.E_ServiceOption).(*api.ServiceOption)
	if !ok {
		return nil, perrors.Newf("型アサーションに失敗しました")
	}

	s := &Service{
		FeatureMaintenanceTypes: serviceOption.GetFeatureMaintenanceTypes(),
		SnakeName:               core.ToSnakeCase(string(protoService.Desc.Name())),
		Comment:                 core.CommentReplacer.Replace(protoService.Comments.Leading.String()),
		Methods:                 make([]*Method, 0, len(protoService.Methods)),
	}
	if s.Comment == "" {
		s.Comment = string(protoService.Desc.Name())
	}
	file.Service = s

	for _, method := range protoService.Methods {
		in, err := createMessage(method.Input)
		if err != nil {
			return nil, perrors.Stack(err)
		}
		out, err := createMessage(method.Output)
		if err != nil {
			return nil, perrors.Stack(err)
		}

		methodOption, ok := proto.GetExtension(method.Desc.Options(), api.E_MethodOption).(*api.MethodOption)
		if !ok {
			return nil, perrors.Newf("型アサーションに失敗しました")
		}

		checkOption, ok := proto.GetExtension(method.Desc.Options(), api.E_CheckOption).(*api.CheckOption)
		if !ok {
			return nil, perrors.Newf("型アサーションに失敗しました")
		}

		errorOption, ok := proto.GetExtension(method.Desc.Options(), api.E_ErrorOption).(*api.ErrorOption)
		if !ok {
			return nil, perrors.Newf("型アサーションに失敗しました")
		}

		httpRule, ok := proto.GetExtension(method.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
		if !ok {
			return nil, perrors.Newf("型アサーションに失敗しました")
		}

		// DisableResponseCacheがついてない場合はEnableResponseCache:trueにする
		if !methodOption.GetDisableResponseCache() {
			if checkOption == nil {
				checkOption = &api.CheckOption{}
			}
			checkOption.EnableResponseCache = true
		}

		// クライアント用にstringからenum変換
		errorCodes := make([]string, 0, len(errorOption.GetErrorCodes()))
		for _, errorCode := range errorOption.GetErrorCodes() {
			errorCodes = append(errorCodes, "ErrorCode_"+errorCode)
		}
		var errorOptionString string
		if len(errorCodes) > 0 {
			errorOptionString = fmt.Sprintf("errorCodes: [ %s ]", strings.Join(errorCodes, ", "))
		}

		httpMethod, httpPath := parseHttpRule(httpRule)
		m := &Method{
			SnakeName:                 core.ToSnakeCase(string(method.Desc.Name())),
			Comment:                   core.CommentReplacer.Replace(method.Comments.Leading.String()),
			InputMessage:              in,
			OutputMessage:             out,
			HttpMethod:                httpMethod,
			HttpPath:                  httpPath,
			DisableCommonResponse:     methodOption.GetDisableCommonResponse(),
			DisableResponseCache:      methodOption.GetDisableResponseCache(),
			DisableCheckMaintenance:   methodOption.GetDisableCheckMaintenance(),
			DisableCheckAppVersion:    methodOption.GetDisableCheckAppVersion(),
			DisableCheckLoginToday:    methodOption.GetDisableCheckLoginToday(),
			DisableFeatureMaintenance: methodOption.GetDisableFeatureMaintenance(),
			FeatureMaintenanceTypes:   methodOption.GetFeatureMaintenanceTypes(),
			DisableGameAuthToken:      checkOption.GetDisableGameAuthToken(),
			DisableMasterVersion:      checkOption.GetDisableMasterVersion(),
			EnableRequestSignature:    checkOption.GetEnableRequestSignature(),
			CheckOption:               checkOption.String(),
			ErrorOption:               errorOptionString,
		}
		if m.Comment == "" {
			m.Comment = string(method.Desc.Name())
		}
		s.Methods = append(s.Methods, m)
	}

	return file, nil
}

func parseHttpRule(httpRule *annotations.HttpRule) (httpMethod, httpPath string) {
	if httpPath = httpRule.GetGet(); httpPath != "" {
		return http.MethodGet, httpPath
	}
	if httpPath = httpRule.GetPut(); httpPath != "" {
		return http.MethodPut, httpPath
	}
	if httpPath = httpRule.GetPost(); httpPath != "" {
		return http.MethodPost, httpPath
	}
	if httpPath = httpRule.GetDelete(); httpPath != "" {
		return http.MethodDelete, httpPath
	}
	if httpPath = httpRule.GetPatch(); httpPath != "" {
		return http.MethodPatch, httpPath
	}
	return "", ""
}

func createMessage(protoMessage *protogen.Message) (*Message, error) {
	message := &Message{
		Messages:  make([]*Message, 0, len(protoMessage.Messages)),
		Fields:    make([]*Field, 0, len(protoMessage.Fields)),
		SnakeName: core.ToSnakeCase(string(protoMessage.Desc.Name())),
		Comment:   core.CommentReplacer.Replace(protoMessage.Comments.Leading.String()),
	}
	if message.Comment == "" {
		message.Comment = string(protoMessage.Desc.Name())
	}

	for _, msg := range protoMessage.Messages {
		// Mapのために生成されたmessageは無視
		if strings.HasSuffix(msg.GoIdent.GoName, "MapEntry") {
			continue
		}
		m, err := createMessage(msg)
		if err != nil {
			return nil, perrors.Stack(err)
		}
		message.Messages = append(message.Messages, m)
	}
	for _, field := range protoMessage.Fields {
		validateOptions, ok := proto.GetExtension(field.Desc.Options(), validateoptions.E_Rules).(fmt.Stringer)
		if !ok {
			return nil, perrors.Newf("型アサーションに失敗しました")
		}
		validateOption := validateOptions.String()
		// https://github.com/golang/go/blob/master/src/fmt/print.go#L20
		if validateOption == "<nil>" {
			validateOption = ""
		}
		isHidden, ok := proto.GetExtension(field.Desc.Options(), zapoptions.E_Hidden).(bool)
		if !ok {
			return nil, perrors.Newf("型アサーションに失敗しました")
		}
		fieldOption, ok := proto.GetExtension(field.Desc.Options(), api.E_FieldOption).(*api.FieldOption)
		if !ok {
			return nil, perrors.Newf("型アサーションに失敗しました")
		}
		var masterRef *MasterRef
		if fieldOption.GetMasterRef() != nil {
			fields := fieldOption.GetMasterRef().GetParentFields()
			parentFieldSnakeNames := make([]string, 0, len(fields))
			for _, f := range fields {
				parentFieldSnakeNames = append(parentFieldSnakeNames, core.ToSnakeCase(f))
			}

			masterRef = &MasterRef{
				TableSnakeName:        fieldOption.GetMasterRef().GetTable(),
				ColumnSnakeName:       fieldOption.GetMasterRef().GetColumn(),
				ParentFieldSnakeNames: parentFieldSnakeNames,
			}
		}

		var typeName string
		var typeKind TypeKind
		var pkgType PkgType
		var importFilePath string
		var parentGoPackage string
		switch field.Desc.Kind() {
		case protoreflect.EnumKind:
			typeName = string(field.Desc.Enum().Name())
			typeKind = TypeKind_Enum
		case protoreflect.MessageKind:
			pkgName := string(field.Desc.Message().ParentFile().Package())
			switch {
			case strings.HasSuffix(pkgName, "game.common"):
				pkgType = PkgType_APICommon
			case strings.HasSuffix(pkgName, "client.common"):
				pkgType = PkgType_ClientCommon
			case strings.HasSuffix(pkgName, "server.common"):
				pkgType = PkgType_ServerCommon
			case strings.HasSuffix(pkgName, "client.transaction"):
				pkgType = PkgType_ClientTransaction
			default:
				pkgType = PkgType_API
			}
			importFilePath = field.Desc.Message().ParentFile().Path()
			parentGoPackage = field.Desc.Message().ParentFile().Options().(*descriptorpb.FileOptions).GetGoPackage()
			typeName = string(field.Desc.Message().Name())
			typeKind = TypeKind_Message
		default:
			typeName = field.Desc.Kind().String()
			typeKind = TypeKind_Primitive
		}
		var mapKeyType, mapValueType string
		if field.Desc.IsMap() {
			typeKind = TypeKind_Map
			typeName = ""
			mapKeyType = field.Desc.MapKey().Kind().String()
			if field.Desc.MapValue().Kind() == protoreflect.MessageKind {
				mapValueType = string(field.Desc.MapValue().Message().Name())
			} else {
				mapValueType = field.Desc.MapValue().Kind().String()
			}
		}

		paths := strings.Split(importFilePath, "/")
		fileName := paths[len(paths)-1]
		if pkgType == PkgType_API &&
			strings.HasSuffix(field.Desc.ParentFile().Path(), fileName) {
			fileName = ""
		}

		f := &Field{
			PkgType:         pkgType,
			ImportFileName:  fileName,
			ParentGoPackage: parentGoPackage,
			SnakeName:       core.ToSnakeCase(string(field.Desc.Name())),
			Comment:         core.CommentReplacer.Replace(field.Comments.Leading.String()),
			Type:            typeName,
			TypeKind:        typeKind,
			IsList:          field.Desc.IsList(),
			IsEnum:          field.Desc.Kind() == protoreflect.EnumKind,
			MapKeyType:      mapKeyType,
			MapValueType:    mapValueType,
			Number:          int32(field.Desc.Number()),
			ValidateOption:  validateOption,
			HiddenOption:    isHidden,
			FieldOption: &FieldOption{
				MasterRef: masterRef,
			},
		}
		if f.Comment == "" {
			f.Comment = string(field.Desc.Name())
		}
		message.Fields = append(message.Fields, f)
	}

	return message, nil
}
