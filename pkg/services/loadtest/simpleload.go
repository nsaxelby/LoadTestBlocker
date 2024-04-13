package loadtest

import (
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

func (l *SimpleLoad) Start(config models.LoadTestConfig) {
	if l.simpleLoadTestRunning {
		return
	} else {
		// Starts the goroutine
		l.simpleLoadTestRunning = true
		go simpleloadtest(l, config)
		log.Println("heartbeat started")
	}
}

func (l *SimpleLoad) Stop() {
	l.simpleLoadTestRunning = false
	log.Println("heartbeat stopped")
}

func simpleloadtest(l *SimpleLoad, config models.LoadTestConfig) {
	client := &http.Client{
		Timeout: time.Second * 1,
	}
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

		watch := stopwatch.Start()
		resp, err := client.Get(config.Url)
		watch.Stop()
		if err != nil {
			failureMessage := strconv.Itoa(l.requestCount) + "  load test failed : " + err.Error()
			log.Println(failureMessage)
			//l.hub.Broadcast <- []byte(failureMessage)
		}

		if resp != nil {
			resp.Body.Close()
		}

		l.requestCount++

		if l.simpleLoadTestRunning == false {
			return
		}
	}
}
