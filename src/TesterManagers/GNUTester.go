package testermanagers

import (
	filetool "TMManager/src/FileTool"
	properties "TMManager/src/Properties"
	sqltool "TMManager/src/SqlTool"
	testcaller "TMManager/src/TestCaller"
	"encoding/json"
	"fmt"
	"strings"
)

func Run(judge *sqltool.Judge, manager *Manager) {
	switch judge.Problem.Judge_mode {
	case 0:
		DefualtRun(judge, manager)
	}
}

var exts = map[string]string{
	"c99":   ".c",
	"c11":   ".c",
	"cpp11": ".cpp",
	"cpp14": ".cpp",
	"cpp17": ".cpp",
	"cpp20": ".cpp",
}

func DefualtRun(judge *sqltool.Judge, manager *Manager) []sqltool.Judge_case {
	judge.Status = "Compiling"
	sqltool.SaveJudge(judge)

	workSpace := manager.WorkSpace + "workspace/"

	strs := strings.Split(judge.Language, " ")
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
		judge.Error_message = res.Info
		judge.Time_use = int(res.RunTime)
		judge.Memory_use = int(res.Memory)
		sqltool.SaveJudge(judge)
		return []sqltool.Judge_case{}
	} else if res.Statu == "017" {
		judge.Status = "SE"
		judge.Error_message = res.Info
		judge.Time_use = int(res.RunTime)
		judge.Memory_use = int(res.Memory)
		sqltool.SaveJudge(judge)
		return []sqltool.Judge_case{}
	}
	cases, err := filetool.ReadCases(casePath)
	if err != nil {
		judge.Status = "SE"
		judge.Error_message = err.Error()
		sqltool.SaveJudge(judge)
		return []sqltool.Judge_case{}
	}

	judge.Status = "Running"
	sqltool.SaveJudge(judge)

	config.SetMount(casePath, "/test/cases/", true)

	judge_cases := make([]sqltool.Judge_case, len(cases))

	for i, ca := range cases {
		config.SetEntry("cugtm", language, version, "run", "/test/workspace/main", fmt.Sprint(judge.Problem.Time_limit), fmt.Sprint(judge.Problem.Memory_limit), "/test/cases/"+ca)
		buf, _ = json.Marshal(&config)
		filetool.WriteFileB(configPath, buf)
		testRes := testcaller.Test(manager.WorkSpace, manager.Name)
		judge_cases[i].JID = judge.ID
		judge_cases[i].Status = testRes.Statu
		judge_cases[i].Time_use = int(testRes.RunTime)
		judge_cases[i].Memory_use = int(testRes.Memory)
		if testRes.RunTime > judge.Time_use {
			judge.Time_use = int(testRes.RunTime)
		}
		if testRes.Memory > judge.Memory_use {
			judge.Memory_use = testRes.Memory
		}
	}
	return judge_cases
}

// "cugtm",
// "gnu",
// "c99",
// "compile",
// "/test/workspace/main",
// "10000",
// "256000"
