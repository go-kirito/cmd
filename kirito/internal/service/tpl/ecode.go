package tpl

import (
	"html/template"
	"os"
)

var ecodeTemplate = `
{{- /* delete empty line */ -}}
package ecode

import "github.com/go-kirito/pkg/errors"

var (
	Err{{.Service}}NotFound  	 = errors.NotFound("{{.Service}}.NotFound", "数据不存在")
	Err{{.Service}}AlreadyExists = errors.BadRequest("{{.Service}}.AlreadyExists", "数据已存在")
)
`

func GenerateEcode(path string, service string) error {

	fileName := path + "/" + service + ".go"

	// 1.检查文件是否存在, 如果存在则不生成
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		return nil
	}

	// 2.生成entity文件
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// 3. write file
	tmpl, err := template.New("ecode").Parse(ecodeTemplate)
	if err != nil {
		return err
	}

	_, toUpperService := convertServiceName(service)

	data := map[string]string{
		"Service": toUpperService,
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}

	return nil
}
