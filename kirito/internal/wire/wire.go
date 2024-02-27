/**
 * @Author : nopsky
 * @Email : cnnopsky@gmail.com
 * @Date : 2021/8/25 14:42
 */
package wire

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-kirito/cmd/kirito/internal/base"
	"github.com/spf13/cobra"
)

var (
	defaultPath         = "di"
	defaultDIFileName   = "di.go"
	defaultWireFileName = "wire.go"
)

var CmdWire = &cobra.Command{
	Use:   "wire",
	Short: "生成wire需要的依赖关系文件",
	Long:  "生成wire需要的依赖关系文件. Example: kirito wire",
	Run:   run,
}

var mod = ""

func run(cmd *cobra.Command, args []string) {

	var err error

	wd, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
		return
	}

	//获取当前项目的mod名称
	mod, err = base.ModulePath(path.Join(wd, "go.mod"))
	if err != nil {
		fmt.Println("go.mod文件不存在，请在项目根路径执行")
		return
	}

	//遍历目录
	dir := args[0]
	if dir == "" {
		dir = "."
	} else {
		dir = strings.TrimSpace(dir)
	}

	err = walk(dir)

	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		if err := os.MkdirAll(defaultPath, 0700); err != nil {
			log.Fatal(err)
		}
	}

	//生成wire文件
	wireContent, err := execute()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = ioutil.WriteFile(defaultPath+"/"+defaultWireFileName, wireContent, 0644); err != nil {
		fmt.Println(err.Error())
		return
	}
}

// 遍历目录
func walk(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if ext := filepath.Ext(path); ext != ".go" || strings.HasPrefix(path, "third_party") || strings.HasPrefix(path, "vendor") {
			return nil
		}
		//解析文件
		return extract(path)
	})
}

// 提取信息
func extract(file string) error {
	fSet := token.NewFileSet()

	f, err := parser.ParseFile(fSet, file, nil, 4)

	if err != nil {
		return err
	}

	for _, i := range f.Decls {
		if fn, ok := i.(*ast.FuncDecl); ok {
			desc := fn.Doc.Text()
			if isContainKeyword(desc, "@di") {
				//提取信息(路径，包名，函数名，参数)
				path := filepath.Dir(file)
				packageName := f.Name.Name
				funcName := fn.Name.Name
				var paramType string
				if fn.Type.Params != nil && len(fn.Type.Params.List) == 2 {
					paramType = fn.Type.Params.List[1].Type.(*ast.Ident).Name
				}

				s := &service{
					Path:         path,
					PackageName:  packageName,
					Func:         funcName,
					ParamType:    paramType,
					VariableName: paramType,
				}
				//判断是否存在同名的package
				v, ok := m[packageName]
				if !ok {
					v = make([]*service, 0)
				}
				v = append(v, s)
				m[packageName] = v
			}

			if isContainKeyword(desc, "@wire") {
				//提取信息(路径，包名，函数名，参数)
				path := filepath.Dir(file)
				packageName := f.Name.Name
				funcName := fn.Name.Name

				s := &service{
					Path:        path,
					PackageName: packageName,
					Func:        funcName,
				}
				//判断是否存在同名的package
				v, ok := w[packageName]
				if !ok {
					v = make([]*service, 0)
				}
				v = append(v, s)
				w[packageName] = v
			}
		}
	}

	return nil
}

// 关键字匹配提取
func isContainKeyword(desc string, keyword string) bool {
	if keyword == "" {
		return true
	}

	if desc == "" {
		return false
	}

	reg, err := regexp.Compile(`^` + keyword)
	if err != nil {
		fmt.Println("keyword err:", err.Error())
		return false
	}

	return reg.MatchString(desc)
}
