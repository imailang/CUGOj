package testcaller

import (
	properties "TMManager/src/Properties"
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"
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

func Test(fileMaps []FileMap, judgerName string, args ...string) TestInfo {
	fileMapCmd := MapFiles(fileMaps)
	cmd := exec.Command("sh",
		"-c",
		fileMapCmd,
		"cd ../img/TestMachine",
		properties.GetAnyway("ContainerName")+" create "+judgerName,
		strings.Join(args, " "),
	)
	cmd.Stdout = &bytes.Buffer{}
	err := cmd.Run()
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
