package main

import (
	"flag"
	"fmt"

	"github.com/go-kirito/cmd/protoc-gen-go-kirito/grpc"
	"github.com/go-kirito/cmd/protoc-gen-go-kirito/http"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const version = "0.0.1"

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	omitempty := flag.Bool("omitempty", true, "omit if google.api is empty")

	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-http %v\n", http.Version)
		fmt.Printf("protoc-gen-go-grpc %v\n", grpc.Version)

		return
	}

	var flags flag.FlagSet
	grpc.RequireUnimplemented = flags.Bool("require_unimplemented_servers", false, "set to false to match legacy behavior")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			hasHttp := true
			if http.GenerateFile(gen, f, *omitempty) == nil {
				hasHttp = false
			}
			grpc.GenerateFile(gen, f)
			generateFile(gen, f, hasHttp)
		}
		return nil
	})
}
