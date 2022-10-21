package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/devfullcycle/fcutils/pkg/events"
)

func main() {
	ed := events.NewEventDispatcher()
	handler := Handler{ID: 1}
	event := ClientCreatedEvent{Name: "client.created", Payload: "test"}
	ed.Register(event.GetName(), &handler)
	ed.Dispatch(&event)
}

type ClientCreatedEvent struct {
	DateTime time.Time
	Payload  interface{}
	Name     string
}

func (e *ClientCreatedEvent) GetDateTime() time.Time {
	return time.Now()
}

func (e *ClientCreatedEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *ClientCreatedEvent) GetName() string {
	return e.Name
}

type Handler struct {
	ID int
}

func (h *Handler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Handler", h.ID, "received event", event.GetName())
	// err <- errors.New("error in handler")
}
