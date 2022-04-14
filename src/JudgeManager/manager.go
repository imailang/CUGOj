package judegemanager

import (
	filetool "TMManager/src/FileTool"
	queuetool "TMManager/src/QueueTool"
	sqltool "TMManager/src/SqlTool"
	testcaller "TMManager/src/TestCaller"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Thread interface {
	run()
}

type Manager struct {
	Name              string                  `gotable:"评测机名"`
	SignalKill        bool                    `gotable:"是否预备停止"`
	WorkSpace         string                  `gotable:"工作路径"`
	BundleConfig      testcaller.BundleConfig ``
	Msgs              <-chan amqp.Delivery    `json:"-"`
	Conn              *amqp.Connection        `json:"-"`
	Ch                *amqp.Channel           `json:"-"`
	LastWorkBeginTime time.Time               `gotable:"上次评测开始时间"`
	WorkTime          time.Duration           `gotable:"上周期工作时长"`
	CreateTime        time.Time               `gotable:"评测机创建时间"`
	Lock              *sync.Mutex             `json:"-"`
}

var Judgers = map[string]*Manager{}
var JudgersLock sync.Mutex
var AssessDelay time.Duration

var CoreNum int
var MaxUseCore int
var Wg = sync.WaitGroup{}

func Initial() {
	CoreNum = runtime.NumCPU()
	MaxUseCore = CoreNum - 1
}

func NewManager(name string) *Manager {
	fmt.Println(filetool.Home() + "img/config.json")
	buf, err := filetool.ReadFile(filetool.Home() + "img/config.json")
	if err != nil {
		fmt.Printf("镜像文件丢失")
		return nil
	}
	config := testcaller.BundleConfig{}
	err = json.Unmarshal(buf, &config)
	if err != nil {
		fmt.Println("镜像配置文件存在错误")
	}
	filetool.Clear(filetool.Home() + "workspace/" + name)
	err = os.Mkdir(filetool.Home()+"workspace/"+name+"/workspace", 0777)
	if err != nil {
		fmt.Println(err)
	}
	return &Manager{
		SignalKill:   false,
		WorkSpace:    filetool.Home() + "workspace/" + name + "/",
		Name:         name,
		BundleConfig: config,
		CreateTime:   time.Now(),
		Lock:         &sync.Mutex{},
	}
}

func (m *Manager) Initial() error {
	conn, err := queuetool.RabbitMQConn()
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}
	q, err := ch.QueueDeclare(
		"judge",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return err
	}
	err = ch.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return err
	}
	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return err
	}
	m.Conn = conn
	m.Ch = ch
	m.Msgs = msgs
	return nil
}

func Status() []Manager {
	JudgersLock.Lock()
	defer JudgersLock.Unlock()
	res := make([]Manager, len(Judgers))
	pos := 0
	for _, j := range Judgers {
		res[pos] = *j
		pos++
	}
	return res
}

func Info() string {
	return "正在运行的评测机：" + fmt.Sprint(len(Judgers))
}

func Stop() string {
	JudgersLock.Lock()
	pos := 0
	keys := make([]string, len(Judgers))
	for k := range Judgers {
		keys[pos] = k
		pos++
	}
	JudgersLock.Unlock()
	log := ""
	for _, key := range keys {
		log += KillJudger(key) + "\n"
	}
	if len(Judgers) != 0 {
		return "评测机关闭出错:\n" + log
	}
	return "所有评测机已关闭"
}

func Assess() {
	JudgersLock.Lock()
	worktime := time.Duration(0)
	for _, j := range Judgers {
		j.Lock.Lock()

		if !j.LastWorkBeginTime.Equal(time.Time{}) {
			j.WorkTime += time.Since(j.LastWorkBeginTime)
			j.LastWorkBeginTime = time.Now()
		}
		worktime += j.WorkTime
		j.WorkTime = time.Duration(0)
		j.Lock.Unlock()
	}
	rate := float64(worktime) / float64(len(Judgers)) / float64(5*time.Second)
	JudgersLock.Unlock()
	if rate < 0.5 && len(Judgers) > 2 {
		StopJudger(GetLastName())
	} else if rate > 0.8 && len(Judgers) <= MaxUseCore {
		for rate > 0.8 && len(Judgers) <= MaxUseCore {
			CreateJudger()
			rate = float64(worktime) / float64(len(Judgers)) / float64(5*time.Second)
		}

	}
}

func GetLegalName() string {
	for i := 1; ; i++ {
		_, ok := Judgers["judger"+fmt.Sprint(i)]
		if !ok {
			return "judger" + fmt.Sprint(i)
		}
	}
}

func GetLastName() string {
	JudgersLock.Lock()
	defer JudgersLock.Unlock()
	Max := 0
	for k, _ := range Judgers {
		tmp, _ := strconv.Atoi(k[6:])
		if tmp > Max {
			Max = tmp
		}
	}
	return "judger" + fmt.Sprint(Max)
}

func CreateJudger() string {
	JudgersLock.Lock()
	defer JudgersLock.Unlock()
	name := GetLegalName()
	m := NewManager(name)
	m.Initial()
	Judgers[name] = m
	go m.Run()
	fmt.Println("创建新评测机：" + name)
	return name
}

func StopJudger(name string) {
	JudgersLock.Lock()
	m, ok := Judgers[name]
	if !ok {
		return
	}
	JudgersLock.Unlock()
	if m.SignalKill {
		KillJudger(m.Name)
		fmt.Println("关闭评测机：" + name + "工作队列")
		return
	}
	m.SignalKill = true
	fmt.Println("向评测机：" + name + "发送停止工作信号")
}

func KillJudger(name string) string {
	JudgersLock.Lock()
	defer JudgersLock.Unlock()
	m, ok := Judgers[name]
	if !ok {
		return "评测机" + name + "不存在"
	}
	err := m.Ch.Close()
	if err != nil {
		return err.Error()
	}
	err = m.Conn.Close()
	if err != nil {
		return err.Error()
	}
	delete(Judgers, name)
	fmt.Println("评测机" + name + "已关闭")
	return "评测机" + name + "已关闭"
}

func AssessRun() {
	tick := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-tick.C:
			Assess()
		}
	}
}

func (m *Manager) Run() {
	fmt.Println("评测机" + m.Name + "启动")
	Wg.Add(1)
	defer Wg.Done()
	for msg := range m.Msgs {
		m.LastWorkBeginTime = time.Now()
		judge := sqltool.QueryJudge(string(msg.Body))
		Run(&judge, m)
		msg.Ack(true)

		m.Lock.Lock()
		m.WorkTime += time.Since(m.LastWorkBeginTime)
		m.LastWorkBeginTime = time.Time{}
		m.Lock.Unlock()

		if m.SignalKill {
			KillJudger(m.Name)
		}
	}
}
