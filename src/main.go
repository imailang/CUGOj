package main

import (
	filetool "TMManager/src/FileTool"
	properties "TMManager/src/Properties"
	sqltool "TMManager/src/SqlTool"
	testcaller "TMManager/src/TestCaller"
	testermanagers "TMManager/src/TesterManagers"
	"encoding/json"
	"fmt"
	"os"
)

//
//参数列表
//	args[1]:
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

	m := testermanagers.NewManager("test")
	m.BundleConfig.SetMount(m.WorkSpace+"workspace/", "/test/workspace", false)
	m.BundleConfig.SetMount(filetool.Home()+"data/problems/testproblem/", "/test/cases", true)
	m.BundleConfig.Root.Path = filetool.Home() + "img/rootfs"
	buf, _ := json.Marshal(&m.BundleConfig)
	filetool.WriteFileB(m.WorkSpace+"config.json", buf)
	code := "#include<stdio.h>int main(){int a,b;scanf(\"%d%d\",&a,&b);printf(\"%d\",a+b);return 0;}"
	filetool.WriteFile("/code/TMManager/workspace/test/workspace/main.c", &code)

	fmt.Println(testcaller.Test(m.WorkSpace, "test").Info)

}
