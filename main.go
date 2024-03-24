package main

import (
	"github.com/stuff-ai/api/internal/server"
)

func main() {
	svc := server.New()
	svc.Start(":1323")
}
