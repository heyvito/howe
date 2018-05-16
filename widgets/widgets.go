package widgets

import "sync"

// HandlerFunc defines a common function for all Widgets to use. It receives an
// arbitrary payload, returning through `output` an string or error object, and
// immediately notifying the `wait` WaitGroup.
type HandlerFunc func(payload map[string]interface{}, output chan interface{}, wait *sync.WaitGroup)

// Handlers holds a list of all known widgets
var Handlers = map[string]HandlerFunc{}

// Register registers a Widget
func Register(name string, handler HandlerFunc) {
	Handlers[name] = handler
}
