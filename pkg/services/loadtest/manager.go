package loadtest

import (
	"github.com/nsaxelby/loadtestblocker/pkg/services/models"
	"github.com/nsaxelby/loadtestblocker/pkg/services/website"
)

type LoadTestManager struct {
	heartbeat  *Heartbeat
	simpleLoad *SimpleLoad
	hub        *website.Hub
}

func NewLoadTestManager(hub *website.Hub) *LoadTestManager {
	return &LoadTestManager{
		heartbeat:  NewHeartbeat(hub),
		hub:        hub,
		simpleLoad: NewSimpleLoad(hub),
	}
}

func (l *LoadTestManager) StartLoadTest(config models.LoadTestConfig) {
	l.hub.Broadcast <- []byte("Starting load test")
	l.heartbeat.StartHeartbeat(config)
	l.simpleLoad.StartSimpleLoadTest(config)
}

func (l *LoadTestManager) StopLoadTest() {
	l.hub.Broadcast <- []byte("Stopping load test")
	l.heartbeat.StopHeartbeat()
	l.simpleLoad.StopSimpleLoadTest()
}
