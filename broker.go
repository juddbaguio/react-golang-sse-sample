package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
)

type EventPayload struct {
	Event *string     `json:"event,omitempty"`
	Data  interface{} `json:"data"`
}

type Broker struct {
	Clients  map[string]chan EventPayload
	Notifier chan EventPayload
	ServErr  chan error
	Mu       sync.Mutex
}

func (b *Broker) Broadcast() {
	for {
		select {
		case data := <-b.Notifier:
			for _, client := range b.Clients {
				client <- data
			}
		case <-b.ServErr:
			return
		}
	}
}

func (b *Broker) SimpleSSEHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Connection does not support streaming", http.StatusBadRequest)
		return
	}

	client := make(chan EventPayload)
	clientId := fmt.Sprintf("%v", rand.Int())
	log.Printf("client %v has connected", clientId)
	defer func() {
		delete(b.Clients, clientId)
		close(client)
		client = nil
	}()
	defer log.Printf("client %v has disconnected", clientId)

	b.Mu.Lock()
	b.Clients[clientId] = client
	b.Mu.Unlock()

	for {
		select {
		case <-r.Context().Done():
			return
		case data := <-client:
			jsonResp, _ := json.Marshal(data.Data)
			if data.Event != nil {
				fmt.Fprintf(w, "event: %v\n", *data.Event)
			}
			fmt.Fprintf(w, "data: %v \n\n", string(jsonResp))
			flusher.Flush()
		}
	}
}

func (b *Broker) NotifierHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var event *EventPayload = &EventPayload{}

	err := json.NewDecoder(r.Body).Decode(event)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	b.Notifier <- *event

	w.Write([]byte("notify success"))
}
