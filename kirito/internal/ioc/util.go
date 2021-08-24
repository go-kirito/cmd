package ioc

import (
	"fmt"
	"github.com/go-kirito/cmd/kirito/internal/ioc/model"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

/**
校验mod是否存在
*/

func CheckMod(GetPath string, TUrl string) {

	switchs := false
	if GetPath == "" {
		log.Fatal("文件路径不可为空")
	}
	file := FileDir(GetPath)
	for _, v := range file {
		if v.Name() == "go.mod" {
			switchs = true
		}
	}
	if !switchs {
		log.Fatal("亲～请在go.mod目录下执行")
	}
	//判断目录是否存在
	if !IsDir(TUrl + EnvDir) {
		//fmt.Println(GetPath+EnvDir)
		err := os.Mkdir(TUrl+EnvDir, 0777)
		if err != nil {
			log.Fatal(fmt.Sprintf("创建%s文件夹失败", EnvDir))
		}
	}
	////判断初始化文件是否存在
	//if !FileExist(GetPath+EventPath){
	//	//创建Event配置文件
	//	file,err:=os.Create(GetPath+EventPath)
	//	if err!=nil {
	//		log.Fatal(fmt.Sprintf("创建%s文件失败",EventPath))
	//	}
	//	file.Write([]byte(tplevent))
	//	defer file.Close()
	//}
	//判断初始化文件是否存在
	if !FileExist(TUrl + WirePath) {
		//创建Event配置文件
		file, err := os.Create(TUrl + WirePath)
		if err != nil {
			log.Fatal(fmt.Sprintf("创建%s文件失败", WirePath))
		}
		file.Write([]byte(tplwire))
		defer file.Close()
	}

}

/**
返回文件列表
*/

func FileDir(GetPath string) []fs.FileInfo {
	GetPath = strings.Replace(GetPath, "\\", "/", -1)
	files, err := ioutil.ReadDir(GetPath)
	if err != nil {
		log.Fatal(err)
	}

	return files
}

/**
返回GoMod模块名
*/

func GetModName(GetPath string) string {
	//fmt.Println(GetPath)
	command := "cd " + GetPath + "; " + "go list -m"
	cmd := exec.Command("/bin/sh", "-c", command)
	//cmd := exec.Command( "go","list","-m")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("combined out:\n%s\n", string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	//fmt.Printf("combined out:\n%s\n", string(out))
	return strings.Replace(string(out), "\n", "", -1)
}

func SetWires(path string) {
	command := "cd " + path + "; " + "wire ."
	cmd := exec.Command("/bin/sh", "-c", command)
	cmdErr := cmd.Run()
	if cmdErr != nil {
		log.Fatal(cmdErr, "编译wire失败")
	}
	//command
	//cmd:=exec.Command("cd"+path,"wire .")
	//cmd.Run()

}

/**
覆盖文件
*/

func WriteToFile(fileName string, content string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("file create failed. err: " + err.Error())
	} else {
		// offset
		//os.Truncate(filename, 0) //clear
		n, _ := f.Seek(0, os.SEEK_END)
		_, err = f.WriteAt([]byte(content), n)
		fmt.Println("write succeed!")
		defer f.Close()
	}
	return err
}

/**
注解规则
*/

func ShouldGen(comment string) bool {
	reg, err := regexp.Compile(`^@wire`)
	if err != nil {
		log.Println(err)
	}
	return reg.MatchString(comment)
}

/**
处理import
*/

func GetImport(RootPath, GetPath string) string {
	rootp := GetPackName(RootPath+"/", 1)

	GetPath, _ = filepath.Split(GetPath)

	ex := fmt.Sprintf("%s(.*)%s", rootp, GetPackName(GetPath, 1))
	reg := regexp.MustCompile(ex)

	strHaiCoder := reg.FindAllString(GetPath, -1)

	if strHaiCoder != nil {

		re3 := regexp.MustCompile("^[^/]*/")
		strHaiCoder1 := re3.ReplaceAllString(strHaiCoder[0], "")
		//fmt.Println(GetModName(GetPath),GetPath)

		return fmt.Sprintf(`"%s"`, GetModName(GetPath)+"/"+strHaiCoder1)
	}
	return ""
}

/**
获取包名
*/

func GetPackName(GetPath string, num int) string {
	pathList := strings.Split(GetPath, "/")
	pack := pathList[len(pathList)-2] //包名
	return pack
}

/**
绝对路径返回
*/

func AsbPath(GetPath string) string {

	//_,file:=filepath.Split(GetPath)
	//if file=="" {
	//	GetPath = GetPackName(GetPath,2)
	//	fmt.Println(GetPath)
	//}
	//path:=""
	//if !filepath.IsAbs(GetPath) {
	//相对路径转绝对路径
	//_,err:= filepath.Abs(GetPath)
	//if err!=nil {
	//	log.Fatal(err)
	//}
	//return strings.Replace(GetPath,"\\","/",-1)
	//return filepath.Split(GetPath)
	//}
	gpath, err := filepath.Abs(GetPath)
	if err != nil {
		log.Fatal(err)
	}
	if !IsDir(gpath) {
		log.Fatal("文件夹路径有问题请检查")
	}
	return gpath
}

/**
判断所给的路径是否为文件夹
*/

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

/**
切片字符串去重复
*/

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}

var tplevent = `package wires

type Event struct {
	// ... TBD
}

//@bin()
func NewEvent() Event {
	return Event{}
}

`

/**
初始化 wire文件信息
*/
var tplwire string = `//+build wireinject

package wires

import (
	"github.com/google/wire"
)

func InitializeEvent() Event {
	wire.Build()

	return Event{}
}
`

func FuncImpDate(data []*model.PackFunc) *model.FuncImport {
	//获取数量和包名
	num := AstPackNum(data)
	impdate := make([]string, 0)
	funcimp := &model.FuncImport{}
	for _, v := range num {
		for k1, v1 := range v.ImpUrl {
			for _, v2 := range data {
				if v2.PackName == v.PackName && v2.Url == v1 {
					//k1+1 别名 数量第几个
					//获取函数名 进行拼接
					funcimp.FuncDate += fmt.Sprintf("\n%s", v2.PackName+strconv.Itoa(k1+1)+"."+v2.FuncName+",")
				}
				//进行追加imp
				impdate = append(impdate, fmt.Sprintf("%s\n", v.PackName+strconv.Itoa(k1+1)+" "+v1))
			}
		}

	}
	//处理最后一个,字符
	funcimp.FuncDate = TrimSuffix(funcimp.FuncDate, ",")
	impdate = RemoveRepeatedElement(impdate)
	for _, v := range impdate {
		funcimp.ImportDate += v
	}

	return funcimp
}

/**
重复包处理及数量
*/

func AstPackNum(data []*model.PackFunc) []*model.PackFuncNum {
	packname := make([]string, 0)
	ImpUrl := make([]string, 0)
	pack_num := make([]*model.PackFuncNum, 0)
	for _, v := range data {

		packname = append(packname, v.PackName)
	}
	//pack name去重
	for _, v := range RemoveRepeatedElement(packname) {
		//遍历当前全部数据
		for _, v1 := range data {
			//数据与pack 匹配
			if v == v1.PackName {
				//存放当前url 数量

				ImpUrl = append(ImpUrl, v1.Url)

			}
		}

		pack_num = append(pack_num, &model.PackFuncNum{PackName: v,
			Num:    len(RemoveRepeatedElement(ImpUrl)),
			ImpUrl: RemoveRepeatedElement(ImpUrl),
		})
		ImpUrl = make([]string, 0) //重新初始化
	}
	return pack_num
}

/**
删除字符串最后一个字符
*/

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {

		s = s[:len(s)-len(suffix)]

	}

	return s

}
