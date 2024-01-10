package event

import (
	"fmt"
	"reflect"
	"sync"
)

type Event struct {
	Type string
	Data interface{}
}

type HandlerType func(event Event)

// EventEngine is the main structure for the event engine.
// It manages event handlers and provides methods to trigger
// and handle events.
type EventEngine struct {
	sync.Mutex
	handlers  map[string][]HandlerType
	recoverer func(interface{}, error)
}

// NewEventEngine creates a new instance of EventEngine.
func NewEventEngine() *EventEngine {
	return &EventEngine{
		handlers: make(map[string][]HandlerType),
	}
}

// Push triggers an event with the given type and data.
// It calls all registered handlers for that event type
// concurrently.
func (e *EventEngine) Push(eventType string, data interface{}) {
	e.Lock()
	handlers, ok := e.handlers[eventType]
	e.Unlock()

	if ok {
		for _, handler := range handlers {
			go func(h HandlerType) {
				defer e.handleRecovery(eventType)
				h(Event{Type: eventType, Data: data})
			}(handler)
		}
	}

}

// Register adds a new handler for a specific event type.
func (e *EventEngine) Register(eventType string, handler HandlerType) {
	e.Lock()
	defer e.Unlock()

	if _, ok := e.handlers[eventType]; !ok {
		e.handlers[eventType] = make([]HandlerType, 0)
	}

	e.handlers[eventType] = append(e.handlers[eventType], handler)
}

// Unregister removes a handler from a specific event type.
func (e *EventEngine) Unregister(eventType string, handler HandlerType) {
	e.Lock()
	defer e.Unlock()

	if handlers, ok := e.handlers[eventType]; ok {
		for i, h := range handlers {
			if reflect.ValueOf(h).Pointer() == reflect.ValueOf(handler).Pointer() {
				e.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

// SetRecoverer sets the recovery function to handle panics.
func (e *EventEngine) SetRecoverer(recoverer func(interface{}, error)) {
	e.recoverer = recoverer
}

// handleRecovery recovers from panics and calls the recoverer function.
func (e *EventEngine) handleRecovery(eventType interface{}) {
	if r := recover(); r != nil {
		if e.recoverer != nil {
			e.recoverer(eventType, fmt.Errorf("%v", r))
		} else {
			fmt.Printf("Recovered from panic in event `%v`: %v\n", eventType, r)
		}
	}
}
