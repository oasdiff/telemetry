package main

import "github.com/oasdiff/telemetry/server"

func main() {

	server.SetupRouter().Run(":8080")
}
