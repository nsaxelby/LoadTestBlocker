package loadtest

import (
	"github.com/nsaxelby/loadtestblocker/pkg/services/models"
	"github.com/nsaxelby/loadtestblocker/pkg/services/website"
)

type LoadTestManager struct {
	heartbeat *Heartbeat
	loadtest  LoadTest
	hub       *website.Hub
}

func NewLoadTestManager(hub *website.Hub) *LoadTestManager {
	return &LoadTestManager{
		heartbeat: NewHeartbeat(hub),
		hub:       hub,
	}
}

func (l *LoadTestManager) StartLoadTest(config models.LoadTestConfig) {
	l.hub.Broadcast <- []byte("Starting load test")
	if config.TargetRPSTestCheck {
		l.loadtest = NewComplexLoad(l.hub)
	} else {
		l.loadtest = NewSimpleLoad(l.hub)
	}
	l.heartbeat.StartHeartbeat(config)
	l.loadtest.Start(config)
}

func (l *LoadTestManager) StopLoadTest() {
	l.hub.Broadcast <- []byte("Stopping load test")
	l.heartbeat.StopHeartbeat()
	l.loadtest.Stop()
}
