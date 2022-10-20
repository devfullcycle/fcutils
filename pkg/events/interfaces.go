package events

import (
	"time"
)

type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface)
	Dispatch(EventInterface)
	Remove(eventName string, handler EventHandlerInterface)
	Has(eventName string, handler EventHandlerInterface) bool
}

type EventHandlerInterface interface {
	Handle(event EventInterface, errs chan error)
}

type EventInterface interface {
	GetDateTime() time.Time
	GetPayload() interface{}
	GetName() string
}
