package queuetool

import (
	properties "TMManager/src/Properties"

	"github.com/streadway/amqp"
)

func RabbitMQConn() (*amqp.Connection, error) {
	user, err := properties.Get("QuUser")
	if err != nil {
		return nil, err
	}
	pwd, err := properties.Get("QuPassword")
	if err != nil {
		return nil, err
	}
	host, err := properties.Get("QuHost")
	if err != nil {
		return nil, err
	}
	port, err := properties.Get("QuPort")
	if err != nil {
		return nil, err
	}
	url := "amqp://" + user + ":" + pwd + "@" + host + ":" + port + "/"
	conn, err := amqp.Dial(url)
	return conn, err
}
