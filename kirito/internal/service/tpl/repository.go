package tpl

import (
	"html/template"
	"log"
	"os"
	"strings"
)

var modelTemplate = `
{{- /* delete empty line */ -}}
package model

import (
	"time"
)

type {{.Service}} struct {
	Id           int64     ` + "`" + `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"` + "`" + `
	Name         string    ` + "`" + `gorm:"column:name;type:varchar(255);comment:名称;NOT NULL"` + "`" + ` // 名称
	PlatformId   int64     ` + "`" + `gorm:"column:platform_id;type:bigint(20);NOT NULL"` + "`" + ` // 平台ID
	CreatedAt    time.Time ` + "`" + `gorm:"column:created_at;type:datetime;default:current_timestamp;comment:'创建时间';NOT NULL"` + "`" + `  // '创建时间'
	UpdatedAt    time.Time ` + "`" + `gorm:"column:updated_at;type:datetime;default:current_timestamp ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间';NOT NULL"` + "`" + ` // '更新时间'
}

func (m *{{.Service}}) TableName() string {
	return "{{.TableName}}"
}
`

func GenerateRepositoryModel(path string, service string) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0700); err != nil {
			return err
		}
	}

	fileName := path + "/" + service + ".go"

	// 1.检查文件是否存在, 如果存在则不生成
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		log.Println("filename:", fileName, err)
		return nil
	}

	// 2.生成model文件
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	// 3. write file
	tmpl, err := template.New("modelTemplate").Parse(modelTemplate)
	if err != nil {
		return err
	}

	_, toUpperService := convertServiceName(service)

	err = tmpl.Execute(file, map[string]string{
		"Service":   toUpperService,
		"TableName": snakeString(service),
	})

	return err

}

var repositoryTemplate = `
{{- /* delete empty line */ -}}
package repository

import (
	"context"
	"{{.GoPackage}}/internal/{{.Domain}}/domain/entity"
	"{{.GoPackage}}/internal/{{.Domain}}/repository/model"

	"github.com/go-kirito/pkg/zdb"

)

type I{{.Service}}Repo interface {
	Get(ctx context.Context, platformId int64, id int64) (*entity.{{.Service}}, error)
	Create(ctx context.Context, platformId int64, {{.ToLowService}} *entity.{{.Service}}) (*entity.{{.Service}}, error)
	Update(ctx context.Context, platformId int64, id int64, updates map[string]interface{}) error
	FindByName(ctx context.Context, platformId int64, name string) (*entity.{{.Service}}, error)
	List(ctx context.Context, platformId int64, filters map[string]interface{}, offset int, limit int) (int64, []*entity.{{.Service}}, error)
	Delete(ctx context.Context, platformId int64, id int64) error
}

type {{.ToLowService}}Repo struct {
}

//@wire
func New{{.Service}}Repo() I{{.Service}}Repo {
	return &{{.ToLowService}}Repo{}
}

func (r {{.ToLowService}}Repo) Get(ctx context.Context, platformId int64, id int64) (*entity.{{.Service}}, error) {
	orm := zdb.NewOrm(ctx)
	
	var m model.{{.Service}}

	if err := orm.Where("id = ? and platform_id = ?", id, platformId).First(&m).Error(); err != nil {
		return nil, err
	}

	return r.toEntity(m), nil
}

func (r {{.ToLowService}}Repo) Create(ctx context.Context, platformId int64, {{.ToLowService}} *entity.{{.Service}}) (*entity.{{.Service}}, error) {
	orm := zdb.NewOrm(ctx)

	m := r.toModel({{.ToLowService}})
	m.PlatformId = platformId

	if err := orm.Create(&m).Error(); err != nil {
		return nil, err
	}

	{{.ToLowService}}.Id = m.Id

	return {{.ToLowService}}, nil
}

func (r {{.ToLowService}}Repo) Update(ctx context.Context, platformId int64, id int64, updates map[string]interface{}) error {
	orm := zdb.NewOrm(ctx)

	return orm.Model(&model.{{.Service}}{}).Where("id = ? and platform_id = ?", id, platformId).Updates(updates).Error()
}

func (r {{.ToLowService}}Repo) FindByName(ctx context.Context, platformId int64, name string) (*entity.{{.Service}}, error) {
	orm := zdb.NewOrm(ctx)

	var m model.{{.Service}}

	if err := orm.Where("name = ? and platform_id = ?", name, platformId).Find(&m).Error(); err != nil {
		return nil, err
	}

	if m.Id == 0 {
		return nil, nil
	}

	return r.toEntity(m), nil
}

func (r {{.ToLowService}}Repo) List(ctx context.Context, platformId int64, filters map[string]interface{}, offset int, limit int) (int64, []*entity.{{.Service}}, error) {
	orm := zdb.NewOrm(ctx)

	var ms []*model.{{.Service}}

	orm = orm.Where("platform_id = ?", platformId)

	if filters != nil {
		if v, ok := filters["ids"]; ok {
			orm = orm.Where("id in (?)", v)
		}

		if v, ok := filters["name"]; ok {
			orm = orm.Where("name like ?", "%"+v.(string)+"%")
		}

		if v, ok := filters["status"]; ok {
			orm = orm.Where("status = ?", v)
		}
	}


	var total int64

	if err := orm.Model(&model.{{.Service}}{}).Count(&total).Error(); err != nil {
		return 0, nil, err
	}

	if offset > 0 {
		orm = orm.Offset(offset)
	}

    if limit > 0 {
		orm = orm.Limit(limit)
    }

	if err := orm.Find(&ms).Error(); err != nil {
		return 0, nil, err
	}

	var es []*entity.{{.Service}}

	for _, m := range ms {
		es = append(es, r.toEntity(*m))
	}

	return total, es, nil
}

func (r {{.ToLowService}}Repo) Delete(ctx context.Context, platformId int64, id int64) error {
	orm := zdb.NewOrm(ctx)

	return orm.Model(&model.{{.Service}}{}).Delete("id = ? and platform_id = ?", id, platformId).Error()
}

func (r {{.ToLowService}}Repo) toEntity(m model.{{.Service}}) *entity.{{.Service}} {
	return &entity.{{.Service}}{
		Id:        m.Id,
		Name:      m.Name,
	}
}

func (r {{.ToLowService}}Repo) toModel(e *entity.{{.Service}}) *model.{{.Service}} {
	return &model.{{.Service}}{
		Id:        e.Id,
		Name:      e.Name,
	}
}
`

func GenerateRepository(path string, goPackage string, domainName, service string) error {
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
	tmpl, err := template.New("repoTemplate").Parse(repositoryTemplate)
	if err != nil {
		return err
	}

	toLowService, toUpperService := convertServiceName(service)

	err = tmpl.Execute(file, map[string]string{
		"Service":      toUpperService,
		"ToLowService": toLowService,
		"Domain":       domainName,
		"GoPackage":    goPackage,
	})

	return nil
}

func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data))
}

func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	flag, num := true, len(s)-1
	for i := 0; i <= num; i++ {
		d := s[i]
		if d == '_' {
			flag = true
			continue
		} else if flag {
			if d >= 'a' && d <= 'z' {
				d = d - 32
			}
			flag = false
		}
		data = append(data, d)
	}
	return string(data)
}
