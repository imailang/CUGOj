package testermanagers

import (
	filetool "TMManager/src/FileTool"
	properties "TMManager/src/Properties"
	sqltool "TMManager/src/SqlTool"
	testcaller "TMManager/src/TestCaller"
	"fmt"
	"strings"
)

func Run(judge *sqltool.Judge, workSpace string) {
	switch judge.Problem.Judge_mode {
	case 0:
		DefualtRun(judge, workSpace)
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

func DefualtRun(judge *sqltool.Judge, workSpace string) []sqltool.Judge_case {
	judge.Status = "Compiling"
	sqltool.SaveJudge(judge)

	strs := strings.Split(judge.Language, "")
	language := strs[0]
	version := strs[1]
	srcPath := workSpace + "/main"
	casePath := "./data/problems/" + fmt.Sprint(judge.ID) + "/cases/"

	fileMaps := []testcaller.FileMap{
		{
			SrcPath: workSpace,
			DesPath: "test/workspace/",
		},
		{
			SrcPath: casePath,
			DesPath: "test/cases/",
		},
	}

	filetool.Clear(workSpace)
	filetool.WriteFile(srcPath+exts[version], &(judge.Code))
	res := testcaller.Test(fileMaps, language, version, "compile", srcPath, properties.GetAnyway("CompileTimeLimit"), properties.GetAnyway("CompileMemoryLimit"))
	if res.Statu == "007" {
		judge.Status = "CE"
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
	judge_cases := make([]sqltool.Judge_case, len(cases))
	for i, ca := range cases {
		testRes := testcaller.Test(fileMaps, language, version, "run", srcPath, fmt.Sprint(judge.Problem.Time_limit), fmt.Sprint(judge.Problem.Memory_limit), casePath+ca)
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
