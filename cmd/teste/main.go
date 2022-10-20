package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/devfullcycle/fcutils/pkg/events"
)

func main() {
	chErr := make(chan error)
	ed := events.NewEventDispatcher()
	handler := Handler{ID: 1}
	event := ClientCreateEvent{Name: "test", Payload: "test"}
	ed.Register("client.create", &handler)

	ed.Dispatch(&event, chErr)

	for {
		select {
		case err := <-chErr:
			fmt.Println(err)
		}
	}
}

type ClientCreateEvent struct {
	DateTime time.Time
	Payload  interface{}
	Name     string
}

func (e *ClientCreateEvent) GetDateTime() time.Time {
	return time.Now()
}

func (e *ClientCreateEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *ClientCreateEvent) GetName() string {
	return e.Name
}

type Handler struct {
	ID int
}

func (h *Handler) Handle(event events.EventInterface, errs chan error) {
	errs <- errors.New("error")
	fmt.Println("handle")
}
