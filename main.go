package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	var broker *Broker = &Broker{
		Clients:  make(map[string]chan EventPayload),
		Notifier: make(chan EventPayload),
		Mu:       sync.Mutex{},
		ServErr:  make(chan error),
	}

	go broker.Broadcast()

	r.HandleFunc("/sse", broker.SimpleSSEHandler)
	r.HandleFunc("/notify", broker.NotifierHandler).Methods(http.MethodPost)

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: r,
	}

	log.Println("Server started at port 3000")

	if err := server.ListenAndServe(); err != nil {
		broker.ServErr <- err
		log.Println(err)
		os.Exit(1)
	}
}
