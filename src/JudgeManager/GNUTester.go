package judegemanager

import (
	filetool "CUGOj/src/FileTool"
	properties "CUGOj/src/Properties"
	sqltool "CUGOj/src/SqlTool"
	testcaller "CUGOj/src/TestCaller"
	"encoding/json"
	"fmt"
	"strings"
)

func Run(judge *sqltool.Judge, manager *Manager) {
	switch judge.Problem.JudgeMode {
	case 0:
		DefualtRun(judge, manager)
	}
}

var statusMap = map[string]string{
	"010": "AC",
	"011": "TLE",
	"012": "RE",
	"013": "MLE",
	"014": "WA",
	"015": "OLE",
	"017": "SE",
}

var exts = map[string]string{
	"c99":   ".c",
	"c11":   ".c",
	"cpp11": ".cpp",
	"cpp14": ".cpp",
	"cpp17": ".cpp",
	"cpp20": ".cpp",
}

func DefualtRun(judge *sqltool.Judge, manager *Manager) []sqltool.JudgeCase {
	judge.Status = "Compiling"
	judge.Judger = manager.Name
	sqltool.SaveJudge(judge)

	workSpace := manager.WorkSpace + "workspace/"

	strs := strings.Split(judge.Language, " ")
	if len(strs) != 2 {
		judge.Status = "CE"
		judge.ErrorMessage = "语言选择存在问题"
		judge.TimeUse = int(0)
		judge.MemoryUse = int(0)
		sqltool.SaveJudge(judge)
		return []sqltool.JudgeCase{}
	}
	language := strs[0]
	version := strs[1]
	srcPath := workSpace + "main"
	casePath := filetool.Home() + "data/problems/" + fmt.Sprint(judge.PID) + "/cases/"
	configPath := manager.WorkSpace + "config.json"

	config := manager.BundleConfig
	config.SetMount(workSpace, "/test/workspace/", false)
	config.Root.Path = filetool.Home() + "img/rootfs"
	config.SetEntry("cugtm", language, version, "compile", "/test/workspace/main", properties.GetAnyway("CompileTimeLimit"), properties.GetAnyway("CompileMemoryLimit"))
	buf, _ := json.Marshal(&config)

	filetool.Clear(workSpace)

	filetool.WriteFileB(configPath, buf)
	filetool.WriteFile(srcPath+exts[version], &(judge.Code))

	res := testcaller.Test(manager.WorkSpace, manager.Name)

	if res.Statu == "007" {
		judge.Status = "CE"
		judge.ErrorMessage = res.Info
		judge.TimeUse = int(res.RunTime)
		judge.MemoryUse = int(res.Memory)
		sqltool.SaveJudge(judge)
		return []sqltool.JudgeCase{}
	} else if res.Statu == "017" {
		judge.Status = "SE"
		judge.ErrorMessage = res.Info
		judge.TimeUse = int(res.RunTime)
		judge.MemoryUse = int(res.Memory)
		sqltool.SaveJudge(judge)
		return []sqltool.JudgeCase{}
	}
	cases, err := filetool.ReadCases(casePath)
	if err != nil {
		judge.Status = "SE"
		judge.ErrorMessage = err.Error()
		sqltool.SaveJudge(judge)
		return []sqltool.JudgeCase{}
	}

	judge.Status = "Running"
	sqltool.SaveJudge(judge)

	config.SetMount(casePath, "/test/cases/", true)

	judgeCases := make([]sqltool.JudgeCase, len(cases))

	judge.Status = "010"
	for i, ca := range cases {
		config.SetEntry("cugtm", language, version, "run", "/test/workspace/main", fmt.Sprint(judge.Problem.TimeLimit), fmt.Sprint(judge.Problem.MemoryLimit), "/test/cases/"+ca)
		buf, _ = json.Marshal(&config)
		filetool.WriteFileB(configPath, buf)
		testRes := testcaller.Test(manager.WorkSpace, manager.Name)
		judgeCases[i].JID = judge.ID
		judgeCases[i].Status = testRes.Statu
		judgeCases[i].TimeUse = int(testRes.RunTime)
		judgeCases[i].MemoryUse = int(testRes.Memory)
		judgeCases[i].CaseId = i + 1
		if testRes.RunTime > judge.TimeUse {
			judge.TimeUse = int(testRes.RunTime)
		}
		if testRes.Memory > judge.MemoryUse {
			judge.MemoryUse = testRes.Memory
		}
		if judge.Status == "010" && testRes.Statu != "010" {
			judge.Status = testRes.Statu
			judge.ErrorMessage = testRes.Info
		}
	}
	judge.Status = statusMap[judge.Status]
	sqltool.SaveJudge(judge)
	sqltool.CreateJudgeCases(&judgeCases)
	sqltool.AddSubmit(judge.PID, judge.Status == "AC")

	return judgeCases
}

// "cugtm",
// "gnu",
// "c99",
// "compile",
// "/test/workspace/main",
// "10000",
// "256000"
