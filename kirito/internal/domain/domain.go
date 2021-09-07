/**
 * @Author : nopsky
 * @Email : cnnopsky@gmail.com
 * @Date : 2021/9/7 18:39
 */
package domain

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var CmdDomain = &cobra.Command{
	Use:   "domain",
	Short: "生成领域结构目录",
	Long:  "生成DDD领域结构目录. Example: kirito domain",
	Run:   run,
}

var dirAll = []string{
	"domain/entity",
	"domain/service",
	"domain/valobj",
	"usecase",
	"repository",
	"facade",
}

var targetDir string

func init() {
	CmdDomain.Flags().StringVarP(&targetDir, "target-dir", "t", "internal/demo", "generate target directory")
}

func run(cmd *cobra.Command, args []string) {

	if targetDir == "" {
		fmt.Fprintln(os.Stderr, "Please specify directory. Example: kirito domain -t internal/demo")
		return
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Printf("Target directory: %s does not exsits\n", targetDir)
		return
	}

	wd, err := os.Getwd()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Get pwd error: %s\n", err.Error())
		return
	}

	for _, dir := range dirAll {
		dir = path.Join(wd, targetDir, dir)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "Create directory error: %s\n", err.Error())
			return
		}
	}
}
