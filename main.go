package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")

func loop(hub *Hub) {
	var tick = 0
	for {
		log.Println("tick", tick)
		hub.broadcast <- []byte(strconv.Itoa(tick))
		tick++
		time.Sleep(time.Second)
	}
}

func main() {
	hub := newHub()
	go hub.run()

	fs := http.FileServer(http.Dir("./web"))
	// Start web server to show basic info, just a text box is fine
	http.Handle("/", fs)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}

	go loop(hub)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	// Setup signalr, or some sort of web socket library to stream results to ther web site, just a console/log is fine for now
}
