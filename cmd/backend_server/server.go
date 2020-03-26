package main

import (
	"github.com/kuso/japanese-word-extractor/server"
)

func main() {
	server := server.NewServer()
	server.SetupRouter()
	//queueName := "test_jlpt_queue"
	//server.SetupMQ("test_jlpt_service", queueName)

	go server.HttpServer.ListenAndServe()
	server.GracefulShutdown(3000)
}
