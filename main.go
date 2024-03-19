package main

import (
	"jaeger-example/client"
	"jaeger-example/server"
)

func main() {
	go server.Run()
	print(client.Get(nil, "jaeger-example-1", "server1", 5678))
}
