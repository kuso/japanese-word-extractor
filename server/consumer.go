package server

import (
	"encoding/json"
	"fmt"
	"github.com/adjust/rmq"
	"log"
	"time"
)

func NewConsumer(tag int, batchSize int, server *Server, out chan<- *Result) *Consumer {
	return &Consumer{
		Server: server,
		Name:   fmt.Sprintf("consumer-%d", tag),
		Count:  0,
		BatchSize: batchSize,
		Before: time.Now(),
		ResultChan: out,
	}
}

func (consumer *Consumer) Consume(delivery rmq.Delivery) {
	var job JobRequest
	if err := json.Unmarshal([]byte(delivery.Payload()), &job); err != nil {
		delivery.Reject()
		return
	}

	result := Result{}
	//result.ConvertedHTML = job.QueryText
	bytes, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
	}
	consumer.Server.Set(job.Id, string(bytes))
	delivery.Ack()
	consumer.ResultChan <- &result
}
