package check

import (
	"bytes"
	_ "embed"
	"fmt"
	"sort"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/input"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/generator/api/output"
	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

//go:embed options.gen.go.tpl
var templateFileBytes []byte

type Data struct {
	Name        string
	PackageName string
	Type        string
	Methods     []string
}

//go:embed feature_maintenance.gen.go.tpl
var featureMaintenanceTemplateFileBytes []byte

type FeatureMaintenanceData struct {
	Method string
	Types  []string
}

type Creator struct{}

func (c *Creator) Create(files []*input.File) ([]*output.TemplateInfo, error) {
	disableCheckMaintenanceMethods := make([]string, 0)
	disableCheckAppVersionMethods := make([]string, 0)
	disableCheckLoginTodayMethods := make([]string, 0)
	disableCheckAuthTokenMethods := make([]string, 0)
	disableCheckMasterVersionMethods := make([]string, 0)
	enableCheckRequestSignature := make([]string, 0)
	enableFeatureMaintenanceMapByMethod := make(map[string][]string)
	for _, file := range files {
		if file.Service == nil {
			continue
		}
		serviceName := core.ToPascalCase(file.Service.SnakeName)
		serviceFeatureMaintenanceTypes := file.Service.FeatureMaintenanceTypes
		for _, method := range file.Service.Methods {
			methodName := fmt.Sprintf("/%s.%s/%s", file.PackageName, serviceName, core.ToPascalCase(method.SnakeName))
			if method.DisableCheckMaintenance {
				disableCheckMaintenanceMethods = append(disableCheckMaintenanceMethods, methodName)
			}
			if method.DisableCheckAppVersion {
				disableCheckAppVersionMethods = append(disableCheckAppVersionMethods, methodName)
			}
			if method.DisableCheckLoginToday {
				disableCheckLoginTodayMethods = append(disableCheckLoginTodayMethods, methodName)
			}
			if method.DisableAuthToken {
				disableCheckAuthTokenMethods = append(disableCheckAuthTokenMethods, methodName)
			}
			if method.DisableMasterVersion {
				disableCheckMasterVersionMethods = append(disableCheckMasterVersionMethods, methodName)
			}
			if method.EnableRequestSignature {
				enableCheckRequestSignature = append(enableCheckRequestSignature, methodName)
			}

			// 機能メンテ
			// メソッド単位で無効な場合はスキップ
			if method.DisableFeatureMaintenance {
				continue
			}
			// メソッド単位での指定で上書き
			if len(method.FeatureMaintenanceTypes) > 0 {
				enableFeatureMaintenanceMapByMethod[methodName] = method.FeatureMaintenanceTypes
				continue
			}
			// サービス単位での指定
			if len(serviceFeatureMaintenanceTypes) > 0 {
				enableFeatureMaintenanceMapByMethod[methodName] = serviceFeatureMaintenanceTypes
			}
		}
	}

	tpl, err := core.GetBaseTemplate().Parse(string(templateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}

	infos := make([]*output.TemplateInfo, 0, 4)

	info, err := c.template(tpl, "Maintenance", "disable", disableCheckMaintenanceMethods)
	if err != nil {
		return nil, perrors.Stack(err)
	}
	infos = append(infos, info)

	info, err = c.template(tpl, "AppVersion", "disable", disableCheckAppVersionMethods)
	if err != nil {
		return nil, perrors.Stack(err)
	}
	infos = append(infos, info)

	info, err = c.template(tpl, "LoginToday", "disable", disableCheckLoginTodayMethods)
	if err != nil {
		return nil, perrors.Stack(err)
	}
	infos = append(infos, info)

	info, err = c.template(tpl, "Auth", "disable", disableCheckAuthTokenMethods)
	if err != nil {
		return nil, perrors.Stack(err)
	}
	infos = append(infos, info)

	info, err = c.template(tpl, "MasterVersion", "disable", disableCheckMasterVersionMethods)
	if err != nil {
		return nil, perrors.Stack(err)
	}
	infos = append(infos, info)

	info, err = c.template(tpl, "RequestSignature", "enable", enableCheckRequestSignature)
	if err != nil {
		return nil, perrors.Stack(err)
	}
	infos = append(infos, info)

	dataList := make([]*FeatureMaintenanceData, 0, len(enableFeatureMaintenanceMapByMethod))
	for method, types := range enableFeatureMaintenanceMapByMethod {
		dataList = append(dataList, &FeatureMaintenanceData{
			Method: method,
			Types:  types,
		})
	}
	info, err = c.featureMaintenanceTemplate(dataList)
	if err != nil {
		return nil, perrors.Stack(err)
	}
	infos = append(infos, info)

	return infos, nil
}

func (c *Creator) template(tpl *template.Template, name, typ string, methods []string) (*output.TemplateInfo, error) {
	packageName := core.ToPkgName(name)
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, &Data{
		Name:        name,
		PackageName: packageName,
		Type:        typ,
		Methods:     methods,
	}); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/cmd/api/middleware", packageName, "options.gen.go"),
	}, nil
}

func (c *Creator) featureMaintenanceTemplate(dataList []*FeatureMaintenanceData) (*output.TemplateInfo, error) {
	featureMaintenanceTpl, err := core.GetBaseTemplate().Funcs(sprig.TxtFuncMap()).Parse(string(featureMaintenanceTemplateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}

	sort.Slice(dataList, func(i, j int) bool {
		return dataList[i].Method < dataList[j].Method
	})

	packageName := "maintenance"
	buf := &bytes.Buffer{}
	if err := featureMaintenanceTpl.Execute(buf, struct {
		DataList    []*FeatureMaintenanceData
		PackageName string
	}{DataList: dataList}); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/cmd/api/middleware", packageName, "feature_maintenance_options.gen.go"),
	}, nil
}
