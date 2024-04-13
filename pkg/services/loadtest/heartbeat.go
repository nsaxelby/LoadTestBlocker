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

type Heartbeat struct {
	heartbeatInterval int
	heartbeatTimeout  int
	heartbeatCount    int
	heartbeatRunning  bool
	timeOfLastSuccess time.Time
	timeOfLastBlock   time.Time
	hub               *website.Hub
}

func NewHeartbeat(hub *website.Hub) *Heartbeat {
	return &Heartbeat{
		heartbeatInterval: 0,
		heartbeatTimeout:  0,
		heartbeatCount:    0,
		heartbeatRunning:  false,
		hub:               hub,
	}
}

func (h *Heartbeat) StartHeartbeat(config models.LoadTestConfig) {
	if h.heartbeatRunning {
		return
	} else {
		// Starts the goroutine
		h.heartbeatRunning = true
		go heartbeat(h, config)
		log.Println("heartbeat started")
	}
}

func (h *Heartbeat) StopHeartbeat() {
	h.heartbeatRunning = false
	log.Println("heartbeat stopped")
}

func heartbeat(h *Heartbeat, config models.LoadTestConfig) {
	client := &http.Client{
		Timeout: time.Second * 1,
	}
	for {
		watch := stopwatch.Start()
		resp, err := client.Get(config.Url)
		watch.Stop()
		if err != nil {
			h.timeOfLastBlock = time.Now()
			failureMessage := "Heartbeat:" + strconv.Itoa(h.heartbeatCount) + "  heartbeat failed : " + err.Error()
			log.Println(failureMessage)
			h.hub.Broadcast <- []byte(failureMessage)
		}

		if resp != nil {
			h.timeOfLastSuccess = time.Now()
			outputMessage := "Heartbeat:" + strconv.Itoa(h.heartbeatCount) + "  heartbeat response : " + resp.Status + fmt.Sprintf(" Milliseconds elapsed: %v", watch.Milliseconds()*time.Millisecond)
			log.Println(outputMessage)
			h.hub.Broadcast <- []byte(outputMessage)
			resp.Body.Close()
		}

		h.heartbeatCount++
		time.Sleep(time.Second)

		if h.heartbeatRunning == false {
			return
		}
	}
}
