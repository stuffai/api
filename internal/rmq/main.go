package rmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection

func init() {
	var err error
	if conn, err = amqp.Dial("amqp://guest:guest@192.168.63.29:5672/"); err != nil {
		panic("failled to initialize amqp: " + err.Error())
	}
}

func Shutdown() {
	conn.Close()
}

func Publish(ctx context.Context, b []byte) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare("text", false, false, false, false, nil)
	if err != nil {
		return err
	}

	if err := ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{Body: b},
	); err != nil {
		return err
	}
	return nil
}
