package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection

func main() {
	// Rabbit
	var err error
	if conn, err = amqp.Dial("amqp://guest:guest@192.168.63.29:5672/"); err != nil {
		panic("failled to initialize amqp: " + err.Error())
	}
	defer conn.Close()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/generate", generate)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

type Request struct {
	Title string `json:"title"`
	Prompt string `json:"prompt"`
}

// Handler
func generate(c echo.Context) error {
	req := new(Request)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := publish(req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
  	return c.String(http.StatusOK, "OK")
}

func publish(req *Request) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare("text", false, false, false, false, nil)
	if err != nil {
		return err
	}

	reqJSON, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if err := ch.PublishWithContext(context.Background(),
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: []byte(reqJSON),
		},
	); err != nil {
		return err
	}
	return nil
}