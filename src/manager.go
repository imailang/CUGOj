package main

import (
	queuetool "TMManager/src/QueueTool"
	sqltool "TMManager/src/SqlTool"
	testermanagers "TMManager/src/TesterManagers"
	"fmt"
)

type Thread interface {
	run()
}

type Manager struct {
	Signal_kill bool
	workSpace   string
}

func NewManager() *Manager {
	return &Manager{
		Signal_kill: false,
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
		testermanagers.Run(&judge, m.workSpace)
	}

}
