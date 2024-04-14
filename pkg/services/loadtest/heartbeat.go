package loadtest

import (
	"encoding/json"
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
	heartbeatIntervalMs int
	heartbeatTimeout    int
	heartbeatCount      int
	heartbeatRunning    bool
	timeOfLastSuccess   time.Time
	timeOfLastBlock     time.Time
	hub                 *website.Hub
}

func NewHeartbeat(hub *website.Hub) *Heartbeat {
	return &Heartbeat{
		heartbeatIntervalMs: 1000,
		heartbeatTimeout:    0,
		heartbeatCount:      0,
		heartbeatRunning:    false,
		hub:                 hub,
	}
}

func (h *Heartbeat) StartHeartbeat(config models.LoadTestConfig) {
	if h.heartbeatRunning {
		return
	} else {
		// Starts the goroutine
		h.heartbeatRunning = true
		h.heartbeatCount = 0
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
			heartbeatObj := &models.ServerHeartbeatEvent{
				Success:   false,
				MSLatency: int(watch.Milliseconds()),
				Timestamp: int64(time.Now().UnixMilli()),
				Message:   err.Error(),
				Count:     h.heartbeatCount,
				Status:    strconv.Itoa(0),
			}

			baseEvent := &models.ServerEvent{
				EventType: "heartbeat",
				Data:      heartbeatObj,
			}

			srvEvent, err := json.Marshal(baseEvent)
			if err != nil {
				fmt.Println(err)
				return
			}

			h.timeOfLastBlock = time.Now()
			failureMessage := "HB Error: " + strconv.Itoa(heartbeatObj.Count) + " - " + heartbeatObj.Status + " msg: " + heartbeatObj.Message
			log.Println(failureMessage)
			h.hub.Broadcast <- []byte(srvEvent)
		}

		if resp != nil {
			h.timeOfLastSuccess = time.Now()

			heartbeatObj := &models.ServerHeartbeatEvent{
				Success:   true,
				MSLatency: int(watch.Milliseconds()),
				Timestamp: int64(time.Now().UnixMilli()),
				Message:   "Status: " + resp.Status,
				Count:     h.heartbeatCount,
				Status:    strconv.Itoa(resp.StatusCode),
			}

			baseEvent := &models.ServerEvent{
				EventType: "heartbeat",
				Data:      heartbeatObj,
			}

			srvEvent, err := json.Marshal(baseEvent)
			if err != nil {
				fmt.Println(err)
				return
			}

			outputMessage := strconv.Itoa(heartbeatObj.Count) + " - " + resp.Status + fmt.Sprintf(" ms: %v", watch.Milliseconds()*time.Millisecond)
			log.Println(outputMessage)
			h.hub.Broadcast <- []byte(srvEvent)
			resp.Body.Close()
		}

		h.heartbeatCount++
		time.Sleep(time.Millisecond * time.Duration(h.heartbeatIntervalMs))

		if h.heartbeatRunning == false {
			return
		}
	}
}
