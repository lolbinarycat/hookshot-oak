package mouse

import (
	"github.com/oakmound/oak/v2/alg/floatgeom"
	"github.com/oakmound/oak/v2/collision"
)

var (
	// LastEvent is the last triggered mouse event,
	// tracked for continuous mouse responsiveness on events
	// that don't take in a mouse event
	LastEvent = NewZeroEvent(0, 0)
	// LastPress is the last triggered mouse event,
	// where the mouse event was a press.
	// If TrackMouseClicks is set to false then this will not be tracked
	LastPress = NewZeroEvent(0, 0)
)

// An Event is passed in through all Mouse related event bindings to
// indicate what type of mouse event was triggered, where it was triggered,
// and which mouse button it concerns.
type Event struct {
	floatgeom.Point2
	Button string
	Event  string
}

// NewEvent creates and returns an Event
func NewEvent(x, y float64, button, event string) Event {
	return Event{
		Point2: floatgeom.Point2{x, y},
		Button: button,
		Event:  event,
	}
}

// NewZeroEvent creates an event with no button or event string.
func NewZeroEvent(x, y float64) Event {
	return NewEvent(x, y, "", "")
}

// ToSpace converts a mouse event into a collision space
func (e Event) ToSpace() *collision.Space {
	return collision.NewUnassignedSpace(e.X(), e.Y(), 0.1, 0.1)
}
