package loadtest

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bradhe/stopwatch"
	"github.com/nsaxelby/loadtestblocker/pkg/services/models"
)

type Heartbeat struct {
	heartbeatInterval int
	heartbeatTimeout  int
	heartbeatCount    int
	heartbeatRunning  bool
	timeOfLastSuccess time.Time
	timeOfLastBlock   time.Time
}

func NewHeartbeat() *Heartbeat {
	return &Heartbeat{
		heartbeatInterval: 0,
		heartbeatTimeout:  0,
		heartbeatCount:    0,
		heartbeatRunning:  false,
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
	for {
		client := &http.Client{
			Timeout: time.Second * 1,
		}
		watch := stopwatch.Start()
		resp, err := client.Get(config.Url)
		watch.Stop()
		if err != nil {
			h.timeOfLastBlock = time.Now()
			log.Println(strconv.Itoa(h.heartbeatCount) + " Failed heartbeat : " + err.Error())
		}

		if resp != nil {
			h.timeOfLastSuccess = time.Now()
			log.Println(strconv.Itoa(h.heartbeatCount) + "  heartbeat response : " + resp.Status + fmt.Sprintf(" Milliseconds elapsed: %v", watch.Milliseconds()*time.Millisecond))
			resp.Body.Close()
		}

		if h.heartbeatRunning == false {
			return
		}
		//log.Println("tick : " + strconv.Itoa(h.heartbeatCount))
		//hub.Broadcast <- []byte(strconv.Itoa(tick))
		h.heartbeatCount++
		time.Sleep(time.Second)
	}
}
