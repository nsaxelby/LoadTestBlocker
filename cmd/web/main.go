package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/nsaxelby/loadtestblocker/pkg/services/loadtest"
	"github.com/nsaxelby/loadtestblocker/pkg/services/website"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	// Hub manages the websocket connections on the website
	hub := website.NewHub()
	go hub.Run()

	// LoadTestManager is what does the heartbeat, and the actual load test
	manager := loadtest.NewLoadTestManager(hub)

	initWebsite(hub, manager)
}

func initWebsite(hub *website.Hub, manager *loadtest.LoadTestManager) {
	fs := http.FileServer(http.Dir("./web"))

	// Start web server to show basic info, just a text box is fine
	http.Handle("/", fs)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		website.ServeWs(hub, manager, w, r)
	})

	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Println("Server started")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
