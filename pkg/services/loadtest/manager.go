package loadtest

import (
	"github.com/nsaxelby/loadtestblocker/pkg/services/models"
)

type LoadTestManager struct {
	heartbeat *Heartbeat
}

func NewLoadTestManager() *LoadTestManager {
	return &LoadTestManager{
		heartbeat: NewHeartbeat(),
	}
}

func (l *LoadTestManager) StartLoadTest(config models.LoadTestConfig) {
	l.heartbeat.StartHeartbeat(config)
}

func (l *LoadTestManager) StopLoadTest() {
	l.heartbeat.StopHeartbeat()
}
