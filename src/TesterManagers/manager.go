package testermanagers

import (
	filetool "TMManager/src/FileTool"
	queuetool "TMManager/src/QueueTool"
	sqltool "TMManager/src/SqlTool"
	testcaller "TMManager/src/TestCaller"
	"encoding/json"
	"fmt"
	"os"
)

type Thread interface {
	run()
}

type Manager struct {
	Signal_kill  bool
	WorkSpace    string
	Name         string
	BundleConfig testcaller.BundleConfig
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
		Signal_kill:  false,
		WorkSpace:    filetool.Home() + "workspace/" + name + "/",
		Name:         name,
		BundleConfig: config,
	}
}

func (m *Manager) showError(err error) {
	fmt.Println(err)
}

func (m *Manager) run() {
	conn, err := queuetool.RabbitMQConn()
	if err != nil {
		m.showError(err)
		return
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		m.showError(err)
		return
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"test",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		m.showError(err)
		return
	}
	err = ch.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		m.showError(err)
		return
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
		m.showError(err)
		return
	}
	for msg := range msgs {
		if m.Signal_kill {
			break
		}
		judge := sqltool.QueryJudge(string(msg.Body))
		Run(&judge, m)
	}

}
