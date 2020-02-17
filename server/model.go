package server

import (
	"github.com/adjust/rmq"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/kuso/japanese-word-extractor/extractor"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	RmqConn    rmq.Connection
	RmqQueue   rmq.Queue
	Pool       *redis.Pool
	Router     *gin.Engine
	HttpServer *http.Server
	StatusMonitor *StatusMonitor
	JLPTDict  *extractor.JLPTDictionary
}

type StatusMonitor struct {
	sync.RWMutex
	OKResults     map[string]int
	ErrResults    map[string]int
	Results       []Result
	ResultChan    chan *Result
}

type JobRequest struct {
	Id          string `json:"id"`
	QueryText   string `json:"querytext"`
}

type JobResult struct {
	Id     string   `json:"id"`
	Value  string   `json:"value"`
}

type Section struct {
	Tokens []*extractor.JLPTToken `json:"tokens"`
}

type Result struct {
	Id string `json:"id"`
	Sections  []*Section `json:"sections"`
}

type Consumer struct {
	Server    *Server
	Name      string
	BatchSize int
	Count     int
	Before    time.Time
	ResultChan chan<- *Result
}

