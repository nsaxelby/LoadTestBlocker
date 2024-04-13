package models

type CommandHandler interface {
	StartLoadTest(config LoadTestConfig)
	StopLoadTest()
}

type LoadTestConfig struct {
	Url                string
	Duration           int
	RatePerSecond      int
	TargetRPSTestCheck bool
}
