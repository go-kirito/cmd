package tpl

import (
	"html/template"
	"os"
)

var serviceTemplate = `
{{- /* delete empty line */ -}}
package service

import (
	"context"
	"{{.GoPackage}}/internal/{{.Domain}}/domain/entity"
	"{{.GoPackage}}/internal/{{.Domain}}/repository"
	"{{.GoPackage}}/ecode"
)

type I{{.Service}}Service interface {
	Get(ctx context.Context, platformId int64, id int64) (*entity.{{.Service}}, error)
	Create(ctx context.Context, platformId int64, {{.ToLowerService}} *entity.{{.Service}}) (*entity.{{.Service}}, error)
	List(ctx context.Context,platformId int64, filters map[string]interface{}, offset, limit int) (int64, []*entity.{{.Service}}, error)
	Update(ctx context.Context, platformId int64, id int64, updates map[string]interface{}) error
	Delete(ctx context.Context,platformId int64, id int64) error
}

type {{.ToLowerService}}Service struct {
	{{.ToLowerService}}Repo repository.I{{.Service}}Repo
}

//@wire
func New{{.Service}}Service(
	a repository.I{{.Service}}Repo,
) I{{.Service}}Service {
	return &{{.ToLowerService}}Service{
		{{.ToLowerService}}Repo: a,
	}
}

func (s *{{.ToLowerService}}Service) Get(ctx context.Context, platformId int64, id int64) (*entity.{{.Service}}, error) {
	return s.{{.ToLowerService}}Repo.Get(ctx, platformId, id)
}

func (s *{{.ToLowerService}}Service) Create(ctx context.Context, platformId int64, {{.ToLowerService}} *entity.{{.Service}}) (*entity.{{.Service}}, error) {
	if err := {{.ToLowerService}}.Verify(); err != nil {
		return nil, err
	}

	_{{.ToLowerService}}, err := s.{{.ToLowerService}}Repo.FindByName(ctx, platformId, {{.ToLowerService}}.Name)
	if err != nil {
		return nil, err
	}

	if _{{.ToLowerService}} != nil {
		return nil, ecode.Err{{.Service}}AlreadyExists
	}

	return s.{{.ToLowerService}}Repo.Create(ctx, platformId, {{.ToLowerService}})
}

func (s *{{.ToLowerService}}Service) List(ctx context.Context, platformId int64, filters map[string]interface{}, offset, limit int) (int64, []*entity.{{.Service}}, error) {
	validFilterKeys := []string{"ids", "name"}
	newFilters := map[string]interface{}{}
	for _, filterKey := range validFilterKeys {
		if filters != nil {
			if _, ok := newFilters[filterKey]; ok {
				continue
			}

			filterValue, ok := filters[filterKey]
			if ok {
				newFilters[filterKey] = filterValue
			}
		}
	}

	return  s.{{.ToLowerService}}Repo.List(ctx, platformId, newFilters, offset, limit)
}

func (s *{{.ToLowerService}}Service) Update(ctx context.Context, platformId int64, id int64, updates map[string]interface{}) error {
	_, err := s.{{.ToLowerService}}Repo.Get(ctx, platformId, id)
	if err != nil {
		return err
	}

	validFilterKeys := []string{"name"}

	newUpdates := map[string]interface{}{}

	for _, filterKey := range validFilterKeys {
		if updates != nil {
			if _, ok := newUpdates[filterKey]; ok {
				continue
			}

			filterValue, ok := updates[filterKey]
			if ok {
				newUpdates[filterKey] = filterValue
			}
		}
	}

	//处理更改以后是否存在重复的情况
	if name, ok := newUpdates["name"]; ok {
		_{{.ToLowerService}}, err := s.{{.ToLowerService}}Repo.FindByName(ctx, platformId, name.(string))
		if err != nil {
			return err
		}

		if _{{.ToLowerService}} != nil && _{{.ToLowerService}}.Id != id {
			return ecode.Err{{.Service}}AlreadyExists
		}
	}

	return s.{{.ToLowerService}}Repo.Update(ctx, platformId, id, newUpdates)
}


func (s *{{.ToLowerService}}Service) Delete(ctx context.Context, platformId int64, id int64) error {
	_, err := s.{{.ToLowerService}}Repo.Get(ctx, platformId, id)
	if err != nil {
		return err
	}

	return s.{{.ToLowerService}}Repo.Delete(ctx, platformId, id)
}

`

// GenerateService 生成service文件
func GenerateService(path string, goPackage string, domainName, service string) error {

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

	// 2.生成service文件
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// 3. write file
	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return err
	}

	toLowService, toUpperService := convertServiceName(service)

	data := map[string]string{
		"Service":        toUpperService,
		"ToLowerService": toLowService,
		"GoPackage":      goPackage,
		"Domain":         domainName,
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}

	return nil
}
