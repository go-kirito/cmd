package ioc

import (
	"fmt"
	"github.com/go-kirito/cmd/kirito/internal/ioc/model"
	"github.com/spf13/cobra"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

/**
路径配置
*/
var (
	EnvDir   string = "/di"
	WirePath string = "/di/wire.go"
)

//IOC cmd
// IOC cmd run project command.

var CmdIoc = &cobra.Command{
	Use:   "wire",
	Short: "New|Update wire service",
	Long:  "wire project. Example: kirito wire",
	Run:   Run,
}

var SUrl string
var TUrl string

func init() {
	CmdIoc.Flags().StringVarP(&SUrl, "SUrl", "s", "./", "目标目录")
	CmdIoc.Flags().StringVarP(&TUrl, "TUrl", "t", "./", "存放目录")
}

func Run(cmd *cobra.Command, args []string) {

	File(strings.TrimSpace(SUrl), TUrl)

}

func File(GetPath string, TUrl string) {
	//转换绝对路径
	GetPath = AsbPath(GetPath)
	TUrl = AsbPath(TUrl)
	//校验go.mod 及初始化配置文件是否存在 不存在则创建
	CheckMod(GetPath, TUrl)

	//递归存放符合条件文件路径
	list := ListFiles(GetPath)

	pack_func := make([]*model.PackFunc, 0)

	for _, file := range list {
		build := FileAnnotation(file) //返回wire.Build
		//存放当前路径
		if build != nil {
			//获取import
			if imp := GetImport(GetPath, file); imp != "" {
				for _, v := range build {
					v.Url = imp
					if imp != "" {
						pack_func = append(pack_func, v)
					}
				}
			}
		}

	}
	PackFuncDate := FuncImpDate(pack_func)

	imports := fmt.Sprintf(`import (
	"github.com/google/wire"
	"github.com/go-kirito/pkg/application"
	 %s
)`, PackFuncDate.ImportDate)

	pack := fmt.Sprintf(`//+build wireinject

package di

%s

func RegisterService(app application.Application) error {
	    wire.Build(%s)

	return nil
}
`, imports, PackFuncDate.FuncDate)

	err := WriteToFile(TUrl+WirePath, pack)
	if err != nil {
		log.Fatal("生成wire失败")
	}

}

//递归调用全文件并保存该地址
var FilePath = make([]string, 0)

func ListFiles(GetPath string) []string {

	files := FileDir(GetPath)
	//存放路径
	//递归循环所有目录
	for _, v := range files {
		PthSep := string(os.PathSeparator)
		filename := GetPath + PthSep + v.Name()
		//拼接当前文件夹下所有文件地址
		if v.Mode().IsDir() == false && path.Ext(v.Name()) == ".go" {
			FilePath = append(FilePath, filename)
		}
		if v.Mode().IsDir() == true {
			//查找下级
			ListFiles(filename) //递归调用
		}
	}
	return FilePath
}

/**
go文件执行 注解
返回：函数方法 ,逗号分割
*/

func FileAnnotation(filePath string) (data []*model.PackFunc) {

	fset := token.NewFileSet() //token值
	//// 这里取绝对路径，方便打印出来的语法树可以转跳到编辑器
	path, _ := filepath.Abs(filePath) //处理的地址
	ast_file, err := parser.ParseFile(fset, path, nil, 4)
	if err != nil {
		log.Println(err)
		return nil
	}

	for _, decl := range ast_file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if ShouldGen(fn.Doc.Text()) {
				data = append(data, &model.PackFunc{PackName: ast_file.Name.Name, FuncName: fn.Name.Name})
			}
		}
	}

	return

}
