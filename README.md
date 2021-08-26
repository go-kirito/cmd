# Kirito

Kirito 来自于 [Kratos](https://github.com/go-kratos/kratos/) 魔改版，新增注入依赖的wire文件，结合wire命令，可自动生成依赖注入文件，开发人员只需要关心业务开发即可



# 环境依赖

- [go](https://golang.org/dl/)
- [protoc](https://github.com/protocolbuffers/protobuf)
- [protoc-gen-go](https://github.com/protocolbuffers/protobuf-go)



# 安装

```shell
go get -u github.com/go-kirito/cmd/kirito
```



# 项目创建

```shel
# create project's layout
kirito new helloworld

cd helloworld
# pull dependencies
go mod download
# generate proto template
kirito proto add api/helloworld/helloworld.proto
# generate proto code
kirito proto client api/helloworld/helloworld.proto
# generate wire
kirito wire .
# generate usecase template
kirito proto server api/helloworld/helloworld.proto -t internal/app/helloworld/usecase

# run
make run
```

