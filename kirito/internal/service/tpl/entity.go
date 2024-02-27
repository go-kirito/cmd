package tpl

import (
	"html/template"
	"log"
	"os"
)

var entityTemplate = `
{{- /* delete empty line */ -}}
package entity

type {{ .Service }} struct {
	Id      int64
	Name    string
}


func(s {{.Service}}) Verify() error {
	return nil
}
`

func GenerateEntity(path string, service string) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0700); err != nil {
			return err
		}
	}

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
	tmpl, err := template.New("entity").Parse(entityTemplate)
	if err != nil {
		return err
	}

	_, toUpperService := convertServiceName(service)

	log.Println("service:", toUpperService)

	data := map[string]string{
		"Service": toUpperService,
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}

	return nil
}
