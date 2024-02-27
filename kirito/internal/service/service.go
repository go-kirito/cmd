package service

import (
	"fmt"
	"github.com/go-kirito/cmd/kirito/internal/proto/add"
	"github.com/go-kirito/cmd/kirito/internal/proto/client"
	"log"
	"os"
	"strings"

	"github.com/go-kirito/cmd/kirito/internal/service/tpl"

	"github.com/go-kirito/cmd/kirito/internal/domain"
	"github.com/spf13/cobra"
)

var CmdService = &cobra.Command{
	Use:   "service",
	Short: "快速开发新的服务",
	Long:  "快速开发新的服务，自动创建proto, domain, entity, service, repository, usecase 等目录. Example: kirito service -d user -e user",
	Run:   run,
}

var domainName string
var entityName string
var repoName string
var all string
var protoName string

func init() {
	CmdService.Flags().StringVarP(&domainName, "domain", "d", "", "领域名称")
	CmdService.Flags().StringVarP(&entityName, "entity", "e", "", "生成实体对象及服务处理代码")
	CmdService.Flags().StringVarP(&repoName, "repo", "r", "", "生成存储层处理代码")
	CmdService.Flags().StringVarP(&all, "all", "a", "", "生成所有代码")
	CmdService.Flags().StringVarP(&protoName, "proto", "p", "", "proto文件路径")
}

func run(cmd *cobra.Command, args []string) {

	wd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	path := strings.Replace(wd, "\\", "/", -1)
	path = strings.TrimRight(path, "/")

	goPackages := strings.Split(path, "/")
	goPackage := goPackages[len(goPackages)-1]

	if domainName == "" && all == "" {
		fmt.Println("请指定领域名称")
		return
	}

	var domainPath string

	if domainName == "" {
		domainName = all
	}

	domainPath = fmt.Sprintf("internal/%s", domainName)

	if all != "" {
		if protoName == "" {
			protoName = domainName
		}

		entityName = all
	}

	if entityName == "" {
		fmt.Println("请指定实体名称")
		return
	}

	if protoName != "" && entityName != "" {
		//1. 创建proto文件
		protoPath := fmt.Sprintf("api/%s/v1/%s.proto", protoName, entityName)
		add.CmdAdd.Run(cmd, []string{protoPath})

		//2. 解析proto文件
		client.CmdClient.Run(cmd, []string{protoPath})

		//7. 生成usecase文件
		usecasePath := fmt.Sprintf("%s/usecase", domainPath)
		err = tpl.GenerateUsecase(usecasePath, goPackage, protoName, domainName, entityName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if domainPath != "" && entityName != "" {
		//3. 生成领域结构目录
		domain.GenerateDomainDirectory(domainPath)

		//4. 生成entity文件
		entityPath := fmt.Sprintf("%s/domain/entity", domainPath)
		err = tpl.GenerateEntity(entityPath, entityName)
		if err != nil {
			fmt.Println(err)
			return
		}

		//5. 生成service文件
		servicePath := fmt.Sprintf("%s/domain/service", domainPath)
		err = tpl.GenerateService(servicePath, goPackage, domainName, entityName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if domainPath != "" && repoName != "" {

		//6. 生成repository文件
		//6.1 生成repository model文件
		repositoryModelPath := fmt.Sprintf("%s/repository/model", domainPath)
		err = tpl.GenerateRepositoryModel(repositoryModelPath, repoName)
		if err != nil {
			fmt.Println(err)
			return
		}

		//6.2 生成repository文件
		repositoryPath := fmt.Sprintf("%s/repository", domainPath)
		err = tpl.GenerateRepository(repositoryPath, goPackage, domainName, repoName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	//8. 生成ecode文件
	err = tpl.GenerateEcode("ecode", entityName)
	if err != nil {
		fmt.Println(err)
		return
	}

}
