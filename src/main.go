package main

import (
	properties "TMManager/src/Properties"
	sqltool "TMManager/src/SqlTool"
	"fmt"
	"os"
)

//
//参数列表
//args[1]:
//  init:创建数据库
func main() {
	err := properties.LoadProperties()
	if err != nil {
		fmt.Println(err)
	}
	err = sqltool.InitialSql()
	if err != nil {
		fmt.Println(err)
	}
	args := os.Args
	argv := len(args)
	if argv == 2 && args[1] == "init" {
		err = sqltool.CreateTables()
		if err != nil {
			fmt.Print(err)
		}
		return
	}
}
