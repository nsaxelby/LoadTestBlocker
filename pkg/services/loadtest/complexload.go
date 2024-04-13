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

type ComplexLoad struct {
	requestTimeout          int
	requestCount            int
	complexLoadTestRunning  bool
	hub                     *website.Hub
	requestsInCurrentSecond int
	currentSecondBenchmark  int
	sleepTime               int
}

func NewComplexLoad(hub *website.Hub) *ComplexLoad {
	return &ComplexLoad{
		requestTimeout:          0,
		requestCount:            0,
		complexLoadTestRunning:  false,
		hub:                     hub,
		requestsInCurrentSecond: 0,
		currentSecondBenchmark:  0,
		sleepTime:               0,
	}
}

func (l *ComplexLoad) Start(config models.LoadTestConfig) {
	if l.complexLoadTestRunning {
		return
	} else {
		// Starts the goroutine
		l.complexLoadTestRunning = true
		rpsReceiveChan := make(chan int)
		go complexloadtest(l, config, rpsReceiveChan)
		go rateAdjuster(l, config, rpsReceiveChan)
		log.Println("heartbeat started")
	}
}

func (l *ComplexLoad) Stop() {
	l.complexLoadTestRunning = false
	log.Println("heartbeat stopped")
}

func complexloadtest(l *ComplexLoad, config models.LoadTestConfig, rpsReportingChan chan int) {
	for {
		currentSecond, err := strconv.Atoi(time.Now().Format("05"))
		if err != nil {
			log.Println(err)
		}

		if currentSecond != l.currentSecondBenchmark {
			l.currentSecondBenchmark = currentSecond
			l.hub.Broadcast <- []byte("RPS: " + strconv.Itoa(l.requestsInCurrentSecond))
			rpsReportingChan <- l.requestsInCurrentSecond
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
			l.hub.Broadcast <- []byte(failureMessage)
		}

		if resp != nil {
			resp.Body.Close()
		}

		l.requestCount++

		if l.complexLoadTestRunning == false {
			rpsReportingChan <- 0
			return
		}
		time.Sleep(time.Duration(l.sleepTime) * time.Microsecond)
	}
}

func rateAdjuster(l *ComplexLoad, config models.LoadTestConfig, rpsReportingChan chan int) {
	sleepTime := 0
	microsecondsToAdjustBy := 5
	for {
		rpsCurrent := <-rpsReportingChan
		if rpsCurrent == 0 {
			continue
		}
		if l.complexLoadTestRunning == false {
			rpsReportingChan <- 0
			return
		}

		// want: 100 , current = 20, diff = 80
		diff := config.RatePerSecond - rpsCurrent
		// absdiff = 80
		absDiff := max(diff, -diff)
		// 80 / 20 = 4
		percentDiff := float32(float32(absDiff)/float32(config.RatePerSecond)) * 100
		log.Println("percent diff: ", percentDiff)
		// 4 x100 = 400

		if percentDiff > 1000 {
			microsecondsToAdjustBy = 50
		} else if percentDiff > 500 {
			microsecondsToAdjustBy = 30
		} else if percentDiff > 250 {
			microsecondsToAdjustBy = 15
		} else if percentDiff > 100 {
			microsecondsToAdjustBy = 10
		} else if percentDiff > 50 {
			microsecondsToAdjustBy = 5
		} else if percentDiff > 10 {
			microsecondsToAdjustBy = 2
		} else {
			microsecondsToAdjustBy = 0
		}

		if diff >= 1 {
			// diff is behind, so number would be positive, which means we need to decrease the sleep time
			sleepTime = sleepTime - microsecondsToAdjustBy
			if sleepTime < 0 {
				sleepTime = 0
			}
		} else {
			// diff is ahead, so number would be negative, which means we need to increase the sleep time
			sleepTime = sleepTime + microsecondsToAdjustBy

		}
		log.Println("sleep time: ", sleepTime)
		log.Println("microsecs: ", microsecondsToAdjustBy)
		l.sleepTime = sleepTime
	}
}