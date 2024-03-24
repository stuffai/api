package main

import (
	"github.com/stuff-ai/api/internal/rmq"
	"github.com/stuff-ai/api/internal/server"
)

func main() {
	defer rmq.Shutdown()

	svc := server.New()
	svc.Start(":1323")
}
