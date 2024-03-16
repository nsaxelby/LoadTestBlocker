package loadtest

import (
	"log"
	"strconv"
	"time"
)

type Heartbeat struct {
	heartbeatInterval int
	heartbeatTimeout  int
	heartbeatCount    int
	heartbeatRunning  bool
}

func NewHeartbeat() *Heartbeat {
	return &Heartbeat{
		heartbeatInterval: 0,
		heartbeatTimeout:  0,
		heartbeatCount:    0,
		heartbeatRunning:  false,
	}
}

func (h *Heartbeat) StartHeartbeat() {
	if h.heartbeatRunning {
		return
	} else {
		// Starts the goroutine
		h.heartbeatRunning = true
		go heartbeat(h)
		log.Println("heartbeat started")
	}
}

func (h *Heartbeat) StopHeartbeat() {
	h.heartbeatRunning = false
	log.Println("heartbeat stopped")
}

func heartbeat(h *Heartbeat) {
	for {
		if h.heartbeatRunning == false {
			return
		}
		log.Println("tick : " + strconv.Itoa(h.heartbeatCount))
		//hub.Broadcast <- []byte(strconv.Itoa(tick))
		h.heartbeatCount++
		time.Sleep(time.Second)
	}
}
