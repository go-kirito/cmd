/**
 * @Author : nopsky
 * @Email : cnnopsky@gmail.com
 * @Date : 2021/9/18 10:04
 */
package model

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/go-kirito/cmd/kirito/internal/model/parser"
	"github.com/spf13/cobra"
)

var CmdModel = &cobra.Command{
	Use:   "model",
	Short: "sql to gorm",
	Long:  "generate gorm struct from sql struct. Example: kirito model",
	Run:   run,
}

type options struct {
	Dns       string
	TableName string
	TargetDir string
}

var opt options

func init() {
	CmdModel.Flags().StringVarP(&opt.Dns, "dns", "", "root:123456@/demo", "mysql host")
	CmdModel.Flags().StringVarP(&opt.TableName, "table", "", "", "mysql table name")
	CmdModel.Flags().StringVarP(&opt.TargetDir, "target-dir", "t", "internal/app/repository", "generate target directory")
}

func run(cmd *cobra.Command, args []string) {

	var err error

	wd, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
		return
	}

	defaultPath := wd + "/" + opt.TargetDir + "/model"

	log.Println("defaultPath:", defaultPath)

	var tables []string
	if opt.TableName != "" {
		tables = append(tables, opt.TableName)
	}

	sqls, err := parser.GetCreateTableFromDB(opt.Dns, tables)
	if err != nil {
		log.Fatal(err)
	}

	o := parser.WithPackage("model")

	for _, sql := range sqls {
		codes, err := parser.ParseSql(sql, o)
		if err != nil {
			log.Fatal(err)
		}

		buf := new(bytes.Buffer)
		tmpl, err := template.New("service").Parse(fileTmplRaw)
		if err != nil {
			log.Fatal(err)
		}

		for _, code := range codes {
			if err := tmpl.Execute(buf, code); err != nil {
				log.Fatal(err)
			}

			if err = ioutil.WriteFile(defaultPath+"/"+strings.ToLower(code.StructCode.RawTableName)+".go", buf.Bytes(), 0644); err != nil {
				fmt.Println(err.Error())
				return
			}

			buf.Reset()
		}
	}

	fd := exec.Command("gofmt", "-s", "-w", defaultPath)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	if err := fd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
		return
	}
	return

	//生成wire文件
	//, err := parser.ParseSql(sql)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//

}
