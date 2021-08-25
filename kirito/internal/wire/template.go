/**
 * @Author : nopsky
 * @Email : cnnopsky@gmail.com
 * @Date : 2021/8/25 16:49
 */
package wire

import (
	"bytes"
	"fmt"
	"html/template"
)

var diTemplate = `
{{- /* delete empty line */ -}}
// Code generated by protoc-gen-go-kirito. DO NOT EDIT.
package di

import (
	{{ range .Services }}
	{{.PackageName}} "{{.Path}}"
	{{- end }}
	"github.com/go-kirito/pkg/application"
)

type UseCases struct {
	{{ range .Services}}
	{{ .VariableName }} {{.PackageName}}.{{.ParamType}}
	{{- end }}
}

func RegisterService(app application.Application) error {
	uc, err := MakeUseCase()
	if err != nil {
		return err
	}
	{{ range .Services }}
	{{.PackageName}}.{{.Func}}(app, uc.{{.VariableName}})
	{{- end }}
	return nil
}
`

var wireTemplate = `
{{- /* delete empty line */ -}}
// Code generated by protoc-gen-go-kirito. DO NOT EDIT.
//+build wireinject

package di

import (
	{{ range .Services }}
	{{.PackageName}} "{{.Path}}"
	{{- end }}

	"github.com/google/wire"

)

func MakeUseCase() (*UseCases, error) {
	panic(wire.Build(
		wire.Struct(new(UseCases), "*"),
		{{ range .Services }}
		{{.PackageName}}.{{.Func}},
		{{- end }}
	))
}
`

var m = map[string][]*service{}
var w = map[string][]*service{}

type service struct {
	Path         string //文件路径
	PackageName  string //包名
	Func         string //函数名
	ParamType    string //参数类型
	VariableName string //参数变量名称
}

type serviceDesc struct {
	Services []*service
}

func execute(tpl string, m map[string][]*service) ([]byte, error) {
	sd := new(serviceDesc)
	sd.Services = make([]*service, 0)
	for packageName, v := range m {
		for k, s := range v {
			//重新定义包名
			rePackageName := fmt.Sprintf("%s%d", packageName, k)
			s.PackageName = rePackageName
			s.Path = fmt.Sprintf("%s/%s", mod, s.Path)
			sd.Services = append(sd.Services, s)
		}
	}

	buf := new(bytes.Buffer)
	tmpl, err := template.New("service").Parse(tpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, sd); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
