package model

/**
包+函数+路径
*/

type PackFunc struct {
	PackName string
	FuncName string
	Url      string
}

type PackFuncNum struct {
	PackName string
	Num      int
	ImpUrl   []string
}

type PackUri struct {
	PackName string
	Url      []string
}

/**
数据组装
*/

type FuncImportAssembly struct {
	PackName string
	ImpUrl   []string
	FuncName string
}

/**
返回最终数据
*/

type FuncImport struct {
	FuncDate   string
	ImportDate string
}
