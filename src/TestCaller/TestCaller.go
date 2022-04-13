package testcaller

import (
	properties "TMManager/src/Properties"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type TestInfo struct {
	Statu   string
	Info    string
	RunTime int
	Memory  int
}

type FileMap struct {
	SrcPath string
	DesPath string
}

func Test(configPath, judgerName string) TestInfo {
	fmt.Println(properties.GetAnyway("ContainerName") + " run -b " + configPath + " " + judgerName)
	cmd := exec.Command(properties.GetAnyway("ContainerName"),
		"run",
		"-b",
		configPath,
		judgerName,
	)
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	fmt.Println(cmd.Stdout.(*bytes.Buffer).String())
	fmt.Println(cmd.Stderr.(*bytes.Buffer).String())
	if err != nil {
		return TestInfo{
			Statu:   "017",
			Info:    "评测机内部错误" + err.Error(),
			RunTime: 0,
			Memory:  0,
		}
	}
	res := TestInfo{}
	err = json.Unmarshal(cmd.Stdout.(*bytes.Buffer).Bytes(), &res)
	if err != nil {
		return TestInfo{
			Statu:   "017",
			Info:    "评测机内部错误" + err.Error(),
			RunTime: 0,
			Memory:  0,
		}
	}
	return res
}

func MapFiles(fileMaps []FileMap) string {
	res := "unshare -m;"
	for _, fm := range fileMaps {
		res += "mount --bind " + fm.SrcPath + " " + fm.DesPath + ";"
	}
	return res
}
