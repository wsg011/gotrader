package event

import (
	"fmt"
	"github.com/wsg011/gotrader/trader/constant"
	"reflect"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type Event struct {
	Type string
	Data interface{}
}

type HandlerType func(interface{})

// EventEngine is the main structure for the event engine.
// It manages event handlers and provides methods to trigger
// and handle events.
type EventEngine struct {
	sync.Mutex
	handlers    map[string][]HandlerType
	recoverer   func(interface{}, error)
	publicChan  chan Event
	privateChan chan Event
}

// NewEventEngine creates a new instance of EventEngine.
func NewEventEngine() *EventEngine {
	return &EventEngine{
		handlers:    make(map[string][]HandlerType),
		publicChan:  make(chan Event, 100),
		privateChan: make(chan Event, 100),
	}
}

// Push triggers an event with the given type and data.
// It calls all registered handlers for that event type
// concurrently.
func (e *EventEngine) Push(eventType string, data interface{}) {
	switch eventType {
	case constant.EVENT_BOOKTICKER:
		e.publicChan <- Event{Type: eventType, Data: data}
	default:
		logrus.Warnf("Unknow event type %s", eventType)
	}
	// e.Lock()
	// handlers := make([]HandlerType, 0)
	// for handlerType := range e.handlers {
	// 	if strings.HasPrefix(eventType, handlerType) {
	// 		handler := e.handlers[handlerType]
	// 		handlers = append(handlers, handler...)
	// 	}
	// }
	// e.Unlock()

	// if len(handlers) > 0 {
	// 	for _, handler := range handlers {
	// 		go func(h HandlerType) {
	// 			defer e.handleRecovery(eventType)
	// 			h(data)
	// 		}(handler)
	// 	}
	// }

}

// Register adds a new handler for a specific event type.
func (e *EventEngine) Register(eventType string, handler func(interface{})) {
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

func (e *EventEngine) Start() {
	e.Lock()
	defer e.Unlock()

	go func() {
		for event := range e.publicChan {
			logrus.Infof("event %s data %s", event.Type, event.Data)

			eventType := event.Type
			data := event.Data

			e.Lock()
			handlers := make([]HandlerType, 0)
			for handlerType := range e.handlers {
				if strings.HasPrefix(eventType, handlerType) {
					handler := e.handlers[handlerType]
					handlers = append(handlers, handler...)
				}
			}
			e.Unlock()

			if len(handlers) > 0 {
				for _, handler := range handlers {
					go func(h HandlerType) {
						defer e.handleRecovery(eventType)
						h(data)
					}(handler)
				}
			}
		}
	}()
}
