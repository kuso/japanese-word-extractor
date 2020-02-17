package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adjust/rmq"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/kuso/japanese-word-extractor/extractor"
	"github.com/lithammer/shortuuid/v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	unackedLimit = 1000
	numConsumers = 2
)

func NewServer() *Server {
	server := Server{}

	server.Pool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				log.Printf("ERROR: fail to init redis: %v", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	defer func() {
		if err := server.Pool.Close(); err != nil {
			panic(err)
		}
	}()

	server.SetupRouter()

	server.HttpServer = &http.Server{
		Addr:           ":8081",
		Handler:        server.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	server.StatusMonitor = &StatusMonitor{}
	server.StatusMonitor.StartMonitoring()

	log.Println("new monitor started")
	log.Println("new server started")

	dict, err := extractor.NewJLPTDictionary()
	if err != nil {
		panic(err)
	}
	server.JLPTDict = dict
	return &server
}

func (server *Server) GracefulShutdown(timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Printf("\nshutdown with timeout: %s\n", timeout)

	if err := server.HttpServer.Shutdown(ctx); err != nil {
		log.Printf("error: %v\n", err)
	} else {
		log.Println("server gracefully stopped")
	}
}

func (server *Server) SetupRouter() {
	server.Router = gin.Default()

	server.Router.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//server.Router.GET("/", server.Hello)
	v1 := server.Router.Group("/v1")
	{
		v1.GET("/hello", server.Hello)
		v1.POST("/job", server.NewJob)
		v1.POST("/jobasync", server.NewJobAsync)
		v1.GET("/job/:id", server.GetJobStatus)
	}
}

func (server *Server) SetupMQ(service string, queue string) {
	server.RmqConn = rmq.OpenConnection(service, "tcp", "localhost:6379", 1)
	server.RmqQueue = server.RmqConn.OpenQueue(queue)
	server.RmqQueue.StartConsuming(unackedLimit, 500*time.Millisecond)
	for i := 0; i < numConsumers; i++ {
		name := fmt.Sprintf("consumer-%d", i)
		server.RmqQueue.AddConsumer(name, NewConsumer(i, 1, server, server.StatusMonitor.ResultChan))
	}
}

/*
func initStore(server *Server) {
	conn := server.Pool.Get()
	defer conn.Close()
}
 */

func (server *Server) Set(key string, val string) error {
	conn := server.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, val)
	if err != nil {
		log.Printf("ERROR: fail set key %s, val %s, error %s", key, val, err.Error())
		return err
	}
	return nil
}

func (server *Server) get(key string) (string, error) {
	conn := server.Pool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Printf("ERROR: fail get key %s, error %s", key, err.Error())
		return "", err
	}
	return s, nil
}

func (server *Server) Hello(c *gin.Context) {
	result := gin.H{"hello": "world",}
	c.JSON(http.StatusOK, result)
}

func (server *Server) NewJob(c *gin.Context) {
	var req JobRequest
	err := c.ShouldBind(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadGateway)
		log.Println(err)
	}

	result := Result{}
	result.Id = shortuuid.New()
	textSections := strings.Split(req.QueryText, "\n")
	for _, textSection := range textSections {
		textSection = strings.TrimSpace(textSection)
		if textSection == "" {
			continue
		}
		tokens := extractor.GetTokens(textSection, server.JLPTDict)
		section := Section{}
		section.Tokens = tokens
		result.Sections = append(result.Sections, &section)
	}

	tmp, err := json.Marshal(result)
	if err != nil {
		c.AbortWithStatus(http.StatusBadGateway)
		log.Println(err)
	}
	jobStr := string(tmp)
	log.Println(jobStr)
	c.String(http.StatusOK, jobStr)
	return
}

func (server *Server) NewJobAsync(c *gin.Context) {
	var req JobRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadGateway)
		log.Println(err)
	}

	uid := shortuuid.New()
	req.Id = uid
	tmp, err := json.Marshal(req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadGateway)
		log.Println(err)
	}

	jobStr := string(tmp)
	result := server.RmqQueue.Publish(jobStr)
	if !result {
		c.AbortWithStatus(http.StatusBadGateway)
		log.Println(result)
	}

	c.String(http.StatusOK, jobStr)
}

func (server *Server) GetJobStatus(c *gin.Context) {
	id := c.Param("id")

	conn := server.Pool.Get()
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()

	if val, err := server.get(id); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		log.Println(err)
	} else {
		var result = JobResult{Id:id, Value:val}
		c.JSON(http.StatusOK, result)
	}
}
