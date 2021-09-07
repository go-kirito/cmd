package server

import (
	"bytes"
	"html/template"
)

var serviceTemplate = `
{{- /* delete empty line */ -}}
package usecase

import (
    "context"

    pb "{{ .Package }}"
    {{- if .GoogleEmpty }}
    "google.golang.org/protobuf/types/known/emptypb"
    {{- end }}
)

type {{ .Service }}UseCase struct {
}

// @wire
func New{{ .Service }}UseCase() pb.I{{ .Service }}UseCase {
	return &{{ .Service }}UseCase{}
}

{{- $s1 := "google.protobuf.Empty" }}
{{ range .Methods }}
func (s *{{ .Service }}UseCase) {{ .Name }}(ctx context.Context, req {{ if eq .Request $s1 }}*emptypb.Empty{{ else }}*pb.{{ .Request }}{{ end }}) ({{ if eq .Reply $s1 }}*emptypb.Empty{{ else }}*pb.{{ .Reply }}{{ end }}, error) {
	return {{ if eq .Reply $s1 }}&emptypb.Empty{}{{ else }}&pb.{{ .Reply }}{}{{ end }}, nil
}
{{- end }}
`

// Service is a proto service.
type Service struct {
	Package     string
	Service     string
	Methods     []*Method
	GoogleEmpty bool
}

// Method is a proto method.
type Method struct {
	Service string
	Name    string
	Request string
	Reply   string
}

func (s *Service) execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, method := range s.Methods {
		if method.Request == "google.protobuf.Empty" || method.Reply == "google.protobuf.Empty" {
			s.GoogleEmpty = true
		}
	}
	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
