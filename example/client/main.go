package main

import (
	"fmt"

	"github.com/9d77v/iec104"
	"github.com/9d77v/iec104/example/client/config"
	"github.com/9d77v/iec104/example/client/worker"
)

func main() {
	address := fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort)
	client := iec104.NewClient(address, config.Logger)
	client.Run(worker.Task)
}
