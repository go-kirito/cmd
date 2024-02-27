package tpl

import (
	"fmt"
	"html/template"
	"os"
	"unicode"
)

var usecaseTemplate = `
{{- /* delete empty line */ -}}
package usecase

import (
    "context"
    pb "{{ .Package }}"
	"{{.GoPackage}}/internal/{{.Domain}}/domain/service"
	"{{.GoPackage}}/utils"
	"{{.GoPackage}}/internal/{{.Domain}}/domain/entity"

)

type {{ .Service }}UseCase struct {
	{{.ToLowService }}Service service.I{{ .Service }}Service
}

// @wire
func New{{ .Service }}UseCase(
	a service.I{{ .Service }}Service,
) pb.I{{ .Service }}UseCase {
	return &{{ .Service }}UseCase{
		{{.ToLowService }}Service: a,
	}
}

func (s *{{ .Service }}UseCase) Get(ctx context.Context, req *pb.Get{{.Service}}Request) (*pb.Get{{.Service}}Reply, error) {
	//TODO 权限校验
	platformId := utils.GetPlatformIdFromContext(ctx)


	{{.ToLowService}}, err := s.{{.ToLowService }}Service.Get(ctx, platformId, req.Id)

	if err != nil {
		return nil, err
	}

	return &pb.Get{{.Service}}Reply{
	 	Item: s.to{{.Service}}Item({{.ToLowService}}),
	}, nil
}

func (s *{{ .Service }}UseCase) Create(ctx context.Context, req *pb.Create{{.Service}}Request) (*pb.Create{{.Service}}Reply, error) {
	//TODO 权限校验
	platformId := utils.GetPlatformIdFromContext(ctx)


	var err error
	{{.ToLowService}} := &entity.{{.Service}}{
		//TODO
	}

	{{.ToLowService}}, err = s.{{.ToLowService }}Service.Create(ctx, platformId, {{.ToLowService}})

	if err != nil {
		return nil, err
	}

	item := s.to{{.Service}}Item({{.ToLowService}})

	return &pb.Create{{.Service}}Reply{
	 	Item: item,
	}, nil
}

func (s *{{ .Service }}UseCase) Update(ctx context.Context, req *pb.Update{{.Service}}Request) (*pb.Update{{.Service}}Reply, error) {
	//TODO 权限校验
	platformId := utils.GetPlatformIdFromContext(ctx)


	var err error

	updates := map[string]interface{}{
		//TODO
	}

	err = s.{{.ToLowService }}Service.Update(ctx, platformId, req.Id, updates)

	if err != nil {
		return nil, err
	}

	return &pb.Update{{.Service}}Reply{
	 	Result: "ok",
	}, nil
}

func (s *{{ .Service }}UseCase) Delete(ctx context.Context, req *pb.Delete{{.Service}}Request) (*pb.Delete{{.Service}}Reply, error) {
	//TODO 权限校验
	platformId := utils.GetPlatformIdFromContext(ctx)


	err := s.{{.ToLowService }}Service.Delete(ctx, platformId, req.Id)

	if err != nil {
		return nil, err
	}

	return &pb.Delete{{.Service}}Reply{
	 	Result: "ok",
	}, nil
}

func (s *{{ .Service }}UseCase) List(ctx context.Context, req *pb.List{{.Service}}Request) (*pb.List{{.Service}}Reply, error) {
	//TODO 权限校验
	platformId := utils.GetPlatformIdFromContext(ctx)


	if req.Limit == 0 {
		req.Limit = utils.DefaultLimit
	}

	if req.Limit > 100 {
		req.Limit = 100
	}

	filters, err := utils.BuildFilter(ctx)
	if err != nil {
		return nil, err
	}

	total, {{.ToLowService}}s, err := s.{{.ToLowService }}Service.List(ctx, platformId, filters, int(req.Offset), int(req.Limit))

	if err != nil {
		return nil, err
	}

	var items []*pb.{{.Service}}Item

	for _, {{.ToLowService}} := range {{.ToLowService}}s {
		items = append(items, s.to{{.Service}}Item({{.ToLowService}}))
	}

	return &pb.List{{.Service}}Reply{
	 	Items: items,
		Total: int32(total),
	}, nil
}

func (s *{{ .Service }}UseCase) to{{.Service}}Item({{.ToLowService}} *entity.{{.Service}}) *pb.{{.Service}}Item {
	return &pb.{{.Service}}Item{
		Id: {{.ToLowService}}.Id,
	}
}

`

// GenerateUsecase 生成usecase文件
func GenerateUsecase(path string, goPackage string, protoName, domainName, service string) error {
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
	tmpl, err := template.New("usecase").Parse(usecaseTemplate)
	if err != nil {
		return err
	}

	toLowService, toUpperService := convertServiceName(service)

	data := map[string]string{
		"Service":      toUpperService,
		"ToLowService": toLowService,
		"GoPackage":    goPackage,
		"Domain":       domainName,
		"Package":      fmt.Sprintf("%s/api/%s/v1", goPackage, protoName),
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}

	return nil
}

func convertServiceName(service string) (toLowerService string, toUpperService string) {
	var lowOutput []rune
	service = CamelString(service)
	for i, r := range service {
		if i == 0 {
			lowOutput = append(lowOutput, unicode.ToLower(r))
		} else {
			lowOutput = append(lowOutput, r)
		}
	}

	return string(lowOutput), service
}
