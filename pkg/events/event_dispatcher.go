package events

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) {
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
}

func (ed *EventDispatcher) Dispatch(event EventInterface) {
	for _, handler := range ed.handlers[event.GetName()] {
		handler.Handle(event)
	}
}

// remove event from dispatecher
func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) {
	for _, h := range ed.handlers[eventName] {
		if h == handler {
			ed.handlers[eventName] = append(ed.handlers[eventName][:0], ed.handlers[eventName][1:]...)
		}
	}
}

func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	for _, h := range ed.handlers[eventName] {
		if h == handler {
			return true
		}
	}
	return false
}

// remove all
func (ed *EventDispatcher) Clear() {
	ed.handlers = make(map[string][]EventHandlerInterface)
}
