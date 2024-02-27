package add

import (
	"bytes"
	"strings"
	"text/template"
)

const protoTemplate = `
syntax = "proto3";

package {{.Package}};
import "google/api/annotations.proto";

option go_package = "{{.GoPackage}}";
option java_multiple_files = true;
option java_package = "{{.JavaPackage}}";

service {{.Service}} {
    rpc Create (Create{{.Service}}Request) returns (Create{{.Service}}Reply){
		option(google.api.http) = {
			post:"/{{.RouteName}}s",
			body:"*",
		};
	}
    rpc Update (Update{{.Service}}Request) returns (Update{{.Service}}Reply){
		option(google.api.http) = {
			put:"/{{.RouteName}}s/{id}",
			body:"*",
		};
	}
    rpc Delete (Delete{{.Service}}Request) returns (Delete{{.Service}}Reply){
		option(google.api.http) = {
			delete:"/{{.RouteName}}s/{id}",
		};
	}
    rpc Get (Get{{.Service}}Request) returns (Get{{.Service}}Reply){
		option(google.api.http) = {
			get:"/{{.RouteName}}s/{id}",
		};
	}
    rpc List (List{{.Service}}Request) returns (List{{.Service}}Reply){
		option(google.api.http) = {
			get:"/{{.RouteName}}s",
		};
	}
}

message Create{{.Service}}Request {
	string name = 1;
}

message Create{{.Service}}Reply {
	{{.Service}}Item item = 1;
}

message Update{{.Service}}Request {
	int64 id = 1;
	string name = 2;
	string status = 3;
}
message Update{{.Service}}Reply {
	string result = 1;
}

message Delete{{.Service}}Request {
	int64 id = 1;
}
message Delete{{.Service}}Reply {
	string result = 1;
}

message Get{{.Service}}Request {
	int64 id = 1;
}
message Get{{.Service}}Reply {
	{{.Service}}Item item = 1;
}

message List{{.Service}}Request {
	int32 offset = 1;
	int32 limit = 2;
}
message List{{.Service}}Reply {
	repeated {{.Service}}Item items = 1;
	int32 total = 2;
}

message {{.Service}}Item {
	int64 id = 1;
	string name = 2;
	string status = 3;
}
`

func (p *Proto) execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("proto").Parse(strings.TrimSpace(protoTemplate))
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
