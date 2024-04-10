package loadtest

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bradhe/stopwatch"
	"github.com/nsaxelby/loadtestblocker/pkg/services/models"
	"github.com/nsaxelby/loadtestblocker/pkg/services/website"
)

type SimpleLoad struct {
	requestTimeout          int
	requestCount            int
	simpleLoadTestRunning   bool
	hub                     *website.Hub
	requestsInCurrentSecond int
	currentSecondBenchmark  int
}

func NewSimpleLoad(hub *website.Hub) *SimpleLoad {
	return &SimpleLoad{
		requestTimeout:          0,
		requestCount:            0,
		simpleLoadTestRunning:   false,
		hub:                     hub,
		requestsInCurrentSecond: 0,
		currentSecondBenchmark:  0,
	}
}

func (l *SimpleLoad) StartSimpleLoadTest(config models.LoadTestConfig) {
	if l.simpleLoadTestRunning {
		return
	} else {
		// Starts the goroutine
		l.simpleLoadTestRunning = true
		go loadtest(l, config)
		log.Println("heartbeat started")
	}
}

func (l *SimpleLoad) StopSimpleLoadTest() {
	l.simpleLoadTestRunning = false
	log.Println("heartbeat stopped")
}

func loadtest(l *SimpleLoad, config models.LoadTestConfig) {
	for {
		currentSecond, err := strconv.Atoi(time.Now().Format("05"))
		if err != nil {
			log.Println(err)
		}

		if currentSecond != l.currentSecondBenchmark {
			l.currentSecondBenchmark = currentSecond
			l.hub.Broadcast <- []byte("RPS: " + strconv.Itoa(l.requestsInCurrentSecond))
			l.requestsInCurrentSecond = 0
		} else {
			l.requestsInCurrentSecond++
		}

		client := &http.Client{
			Timeout: time.Second * 1,
		}
		watch := stopwatch.Start()
		resp, err := client.Get(config.Url)
		watch.Stop()
		if err != nil {
			failureMessage := strconv.Itoa(l.requestCount) + "  load test failed : " + err.Error()
			log.Println(failureMessage)
			//l.hub.Broadcast <- []byte(failureMessage)
		}

		if resp != nil {
			outputMessage := strconv.Itoa(l.requestCount) + "  load test response : " + resp.Status + fmt.Sprintf(" Milliseconds elapsed: %v", watch.Milliseconds()*time.Millisecond)
			log.Println(outputMessage)
			//l.hub.Broadcast <- []byte(outputMessage)
			resp.Body.Close()
		}

		l.requestCount++

		if l.simpleLoadTestRunning == false {
			return
		}
	}
}
