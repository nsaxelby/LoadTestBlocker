package loadtest

import (
	"github.com/nsaxelby/loadtestblocker/pkg/services/models"
)

type LoadTest interface {
	Start(config models.LoadTestConfig)
	Stop()
}
