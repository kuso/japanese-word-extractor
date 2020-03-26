package server

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

func Test_NewServer(t *testing.T) {
	server := NewServer()
	log.Println(server)
}

func Test_Server_Set_Get(t *testing.T) {
	mockConn := redigomock.NewConn()
	mockPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return mockConn, nil
		},
	}

	mockConn.Command("SET", "TESTKEY", "test_value")
	mockConn.Command("GET", "TESTKEY").Expect("test_value")

	server := &Server{}
	server.Pool = mockPool

	server.Set("TESTKEY", "test_value")
	got, err := server.get("TESTKEY")
	if got != "test_value" {
		t.Errorf("got %v, want test_value", got)
	}
	if err != nil {
		t.Errorf("%v", err)
	}
}

func Test_HelloWorld(t *testing.T) {
	mockConn := redigomock.NewConn()
	mockPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return mockConn, nil
		},
	}
	server := &Server{}
	server.Pool = mockPool
	server.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/hello", nil)
	server.Router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	body := gin.H{
		"hello": "world",
	}
	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	value, exists := response["hello"]
	assert.True(t, exists)
	assert.Equal(t, body["hello"], value)
}

/*
func Test_GetJobStatus_method_1(t *testing.T) {
	mockConn := redigomock.NewConn()
	mockPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return mockConn, nil
		},
	}
	server := &Server{}
	server.Pool = mockPool
	server.SetupRouter()

	mockResult := JobResult{
		Id: "test_job",
		Value: "OK",
	}
	_, err := json.Marshal(mockResult)
	if err != nil {
		log.Fatal(err)
		return
	}
	mockConn.Command("GET", "test_job").Expect(mockResult.Value)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/job/test_job", nil)
	server.Router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)

	assert.Equal(t, mockResult.Id, response["id"])
	assert.Equal(t, mockResult.Value, response["value"])
}

func Test_GetJobStatus(t *testing.T) {
	mockConn := redigomock.NewConn()
	mockPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return mockConn, nil
		},
	}
	server := &Server{}
	server.Pool = mockPool
	server.SetupRouter()

	mockResult := JobResult{
		Id: "test_job",
		Value: "OK",
	}
	str, err := json.Marshal(mockResult)
	if err != nil {
		log.Fatal(err)
		return
	}
	mockConn.Command("GET", "test_job").Expect(mockResult.Value)

	apitest.New().
		Handler(server.Router).
		Get("/v1/job/test_job").
		Expect(t).
		Body(string(str)).
		Assert(jsonpath.Equal(`$.id`, mockResult.Id)).
		Assert(jsonpath.Equal(`$.value`, mockResult.Value)).
		Status(http.StatusOK).
		End()
}

func Test_NewJob(t *testing.T) {
	server := &Server{}
	server.SetupRouter()

	mockConn := redigomock.NewConn()
	mockPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return mockConn, nil
		},
	}
	server.Pool = mockPool

	queueName := "test_tasks"

	mockRmqConn := rmq.NewTestConnection()
	server.RmqConn = mockRmqConn
	server.RmqQueue = server.RmqConn.OpenQueue(queueName)

	mockRequest := JobRequest{
		QueryText: "test_content_text",
	}
	tmp, err := json.Marshal(mockRequest)
	assert.Nil(t, err)
	str := string(tmp)

	apitest.New().
		Handler(server.Router).
		Post("/v1/job").
		JSON(str).
		Expect(t).
		Assert(jsonpath.Present(`$.id`)).
		Assert(jsonpath.Equal(`$.contenttext`, "test_content_text")).
		Status(http.StatusOK).
		End()

	// get mock task id
	content := mockRmqConn.GetDelivery(queueName, 0)
	mockTask :=  JobRequest{}
	err = json.Unmarshal([]byte(content), &mockTask)
	assert.Nil(t, err)
	log.Println(mockTask)

	// set the mock result simulating result has been saved to redis
	mockResult := JobResult{
		Id: mockTask.Id,
		Value: "OK",
	}
	_, err = json.Marshal(mockResult)
	assert.Nil(t, err)
	mockConn.Command("GET", mockTask.Id).Expect(mockResult.Value)

	// query the server for result
	apitest.New().
		Handler(server.Router).
		Get("/v1/job/" + mockTask.Id).
		Expect(t).
		Assert(jsonpath.Equal(`$.id`, mockTask.Id)).
		Assert(jsonpath.Equal(`$.value`, mockResult.Value)).
		Status(http.StatusOK).
		End()
}

func Test_Live_NewJob_Sync(t *testing.T) {
	queueName := "test_jlpt_queue"
	server := NewServer()
	server.SetupRouter()
	server.SetupMQ("test_jlpt_service", queueName)

	mockRequest := JobRequest{
		QueryText: "寿司が食べたい。\n",
	}
	tmp, err := json.Marshal(mockRequest)
	assert.Nil(t, err)
	str := string(tmp)

	apitest.New().
		Handler(server.Router).
		Post("/v1/job").
		JSON(str).
		Expect(t).
		Assert(func(res *http.Response, req *http.Request) error {
			assert.Equal(t, http.StatusOK, res.StatusCode)
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err.Error())
			}
			log.Println(jsonPrettyPrint(string(body)))
			var result = Result{}
			err = json.Unmarshal(body, &result)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, 1, len(result.Sections))
			assert.Equal(t, 5, len(result.Sections[0].Tokens))
			return nil
		}).
		Status(http.StatusOK).
		End()

	server.RmqQueue.StopConsuming()
}

func Test_Live_NewJob_Sync1(t *testing.T) {
	queueName := "test_jlpt_queue"
	server := NewServer()
	server.SetupRouter()
	server.SetupMQ("test_jlpt_service", queueName)

	mockRequest := JobRequest{
		QueryText: "対等\n",
	}
	tmp, err := json.Marshal(mockRequest)
	assert.Nil(t, err)
	str := string(tmp)

	apitest.New().
		Handler(server.Router).
		Post("/v1/job").
		JSON(str).
		Expect(t).
		Assert(func(res *http.Response, req *http.Request) error {
			assert.Equal(t, http.StatusOK, res.StatusCode)
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err.Error())
			}
			log.Println(jsonPrettyPrint(string(body)))
			var result = Result{}
			err = json.Unmarshal(body, &result)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, 1, len(result.Sections))
			return nil
		}).
		Status(http.StatusOK).
		End()

	server.RmqQueue.StopConsuming()
}


func Test_Live_NewJob_Sync_3(t *testing.T) {
	queueName := "test_jlpt_queue"
	server := NewServer()
	server.SetupRouter()
	server.SetupMQ("test_jlpt_service", queueName)

	mockRequest := JobRequest{
		QueryText: "通らずに\n",
	}
	tmp, err := json.Marshal(mockRequest)
	assert.Nil(t, err)
	str := string(tmp)

	apitest.New().
		Handler(server.Router).
		Post("/v1/job").
		JSON(str).
		Expect(t).
		Assert(func(res *http.Response, req *http.Request) error {
			assert.Equal(t, http.StatusOK, res.StatusCode)
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err.Error())
			}
			log.Println(jsonPrettyPrint(string(body)))
			var result = Result{}
			err = json.Unmarshal(body, &result)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, 1, len(result.Sections))
			return nil
		}).
		Status(http.StatusOK).
		End()

	server.RmqQueue.StopConsuming()
}

*/
