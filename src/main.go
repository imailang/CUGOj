package main

import (
	judgemanager "CUGOj/src/JudgeManager"
	nettool "CUGOj/src/NetTool"
	properties "CUGOj/src/Properties"
	sqltool "CUGOj/src/SqlTool"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/liushuochen/gotable"
)

var pid = -1
var pidd = -1

func Getpid() {
	res, err := http.Get("http://localhost:13001/pid")
	if err != nil {
		//fmt.Println(err)
		return
	}
	buf, err := ioutil.ReadAll(res.Body)

	if err != nil {
		//fmt.Println(err)
		return
	}
	res.Body.Close()
	pid, _ = strconv.Atoi(string(buf))

	//fmt.Println(pid)
}

func Getpidd() {
	res, err := http.Get("http://localhost:13002/pid")
	if err != nil {
		//fmt.Println(err)
		return
	}
	buf, err := ioutil.ReadAll(res.Body)

	if err != nil {
		//fmt.Println(err)
		return
	}
	res.Body.Close()
	pidd, _ = strconv.Atoi(string(buf))

}

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

	Getpid()
	Getpidd()

	if argv >= 2 {
		if args[1] == "init" {
			err = sqltool.CreateTables()
			if err != nil {
				fmt.Print(err)
			}
		} else if args[1] == "start" {
			if pid != -1 {
				fmt.Println("评测机调度机已存在")
				return
			} else if pidd != -1 {
				fmt.Println("评测调度机守护进程存在，但是评测调度机未找到，请检查日志文件")
				return
			}
			cmd := exec.Command(args[0], "__startd")
			err = cmd.Start()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("CUGOJ 评测调度机守护进程启动成功 进程号：" + fmt.Sprint(cmd.Process.Pid))
			}

		} else if args[1] == "pid" {
			fmt.Println(pid)
		} else if args[1] == "pidd" {
			fmt.Println(pidd)
		} else if args[1] == "__startd" {
			nettool.Rund()
			nettool.Signal_Stopd = false
			for !nettool.Signal_Stopd {
				fmt.Println("启动评测调度机进程")
				cmd := exec.Command(args[0], "__start")
				cmd.Stdout = os.Stdout
				err = cmd.Run()
				if err != nil {
					fmt.Print(err)
				}
			}

		} else if args[1] == "__start" {
			fmt.Println("评测调度机启动")
			judgemanager.Initial()
			fmt.Println("评测调度机初始化完成")
			go nettool.Run()
			fmt.Println("评测调度机开始监听13001端口")
			judgemanager.CreateJudger()
			if judgemanager.MaxUseCore > 1 {
				judgemanager.CreateJudger()
			} else {
				fmt.Println("核心数过少，建议在多核环境运行评测机")
			}

			go judgemanager.AssessRun()
			judgemanager.Wg.Wait()
		} else if args[1] == "stop" {
			if pid == -1 {
				if pidd != -1 {
					fmt.Println("评测调度机守护进程存在，但是评测调度机未找到，请检查日志文件")
					return
				}
				fmt.Println("评测机未启动")
				return
			}
			_, err := http.Get("http://localhost:13002/stop")
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("评测机关闭")

		} else if args[1] == "status" {
			if pid == -1 {
				if pidd != -1 {
					fmt.Println("评测调度机守护进程存在，但是评测调度机未找到，请检查日志文件")
					return
				}
				fmt.Println("评测机未启动")
				return
			}
			res, err := http.Get("http://localhost:13001/status")

			if err != nil {
				fmt.Println(err)
				return
			}
			buf, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			res.Body.Close()
			modal := nettool.StatusModal{}
			err = json.Unmarshal(buf, &modal)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(modal.Statu)
			if argv >= 3 && args[2] == "-l" {
				tab, _ := gotable.Create("评测机名", "停止信号", "工作路径", "上次评测开始时间", "上周期工作时长", "评测机创建时间")
				for _, judge := range modal.Judgers {
					tab.AddRow([]string{
						judge.Name,
						fmt.Sprint(judge.SignalKill),
						judge.WorkSpace,
						fmt.Sprint(judge.LastWorkBeginTime),
						fmt.Sprint(judge.WorkTime),
						fmt.Sprint(judge.CreateTime),
					})
				}
				fmt.Println(tab)
			}
		} else {
			fmt.Println("参数" + args[1] + "错误，无对应操作")
		}
	} else {
		fmt.Println("参数过少")
		fmt.Println("start  启动评测机")
		fmt.Println("stop   停止评测机工作")
		fmt.Println("status 查看评测机状态，可选参数-l 显示每个评测机的状态")
	}

}
