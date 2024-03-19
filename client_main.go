package main

import "jaeger-example/client"

func main() {
	println(client.Get(nil, "Start", "server1", 5678))
}
