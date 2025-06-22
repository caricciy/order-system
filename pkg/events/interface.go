package events

import (
	"sync"
	"time"
)

type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() interface{}
	SetPayload(payload interface{})
}

type EventHandlerInterface interface {
	// Handle is responsible for processing the event.
	Handle(event EventInterface, wg *sync.WaitGroup)
}

type EventDispatcherInterface interface {
	// Register is responsible for registering a handler for a specific event.
	Register(eventName string, handler EventHandlerInterface) error

	// Dispatch is responsible for dispatching the event to all registered handlers.
	Dispatch(event EventInterface) error

	// Remove is responsible for removing a handler for a specific event.
	Remove(eventName string, handler EventHandlerInterface) error

	// Has checks if a handler is registered for a specific event.
	Has(eventName string, handler EventHandlerInterface) bool

	// Clear is responsible for clearing all registered handlers.
	Clear()
}
