package nettool

import (
	judgemanager "TMManager/src/JudgeManager"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

var Signal_Stopd = false

type StatusModal struct {
	Judgers []judgemanager.Manager
	Statu   string
}

var Api = map[string]func(http.ResponseWriter, *http.Request){
	"/status": Status,
	"/stop":   Stop,
	"/pid":    Pid,
}

func Status(w http.ResponseWriter, r *http.Request) {
	res := StatusModal{
		Judgers: judgemanager.Status(),
		Statu:   judgemanager.Info(),
	}
	buf, _ := json.Marshal(&res)
	w.Write(buf)
}

func Stop(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(judgemanager.Stop()))
	judgemanager.Wg.Done()
}

func Pid(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprint(os.Getpid())))
}

func Run() {
	judgemanager.Wg.Add(1)
	r := chi.NewRouter()
	for k, v := range Api {
		r.Get(k, v)
	}
	http.ListenAndServe("localhost:13001", r)
}

func Stopd(w http.ResponseWriter, r *http.Request) {
	Signal_Stopd = true
	http.Get("http://localhost:13001/stop")
}

func Pidd(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprint(os.Getpid())))
}
func Rund() {
	fmt.Println("配置网络")
	r := chi.NewRouter()
	r.Get("/stop", Stopd)
	r.Get("/pid", Pidd)
	fmt.Println("路由成功，开始监听端口13002")
	go http.ListenAndServe("localhost:13002", r)
}
