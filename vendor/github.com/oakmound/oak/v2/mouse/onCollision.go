package mouse

import (
	"errors"

	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/event"
)

// CollisionPhase is a component that can be placed into another struct to
// enable PhaseCollision on the struct. See PhaseCollision.
type CollisionPhase struct {
	OnCollisionS *collision.Space

	wasTouching bool
}

func (cp *CollisionPhase) getCollisionPhase() *CollisionPhase {
	return cp
}

type collisionPhase interface {
	getCollisionPhase() *CollisionPhase
}

// PhaseCollision binds to the entity behind the space's CID so that it will
// receive MouseCollisionStart and MouseCollisionStop events, appropriately when
// the mouse begins to hover or stops hovering over the input space.
func PhaseCollision(s *collision.Space) error {
	en := event.GetEntity(int(s.CID))
	if cp, ok := en.(collisionPhase); ok {
		oc := cp.getCollisionPhase()
		oc.OnCollisionS = s
		s.CID.Bind(phaseCollisionEnter, event.Enter)
		return nil
	}
	return errors.New("This space's entity does not implement collisionPhase")
}

// MouseCollisionStart/Stop: see collision Start/Stop, for mouse collision
// Payload: (mouse.Event)
const (
	Start = "MouseCollisionStart"
	Stop  = "MouseCollisionStop"
)

func phaseCollisionEnter(id int, nothing interface{}) int {
	e := event.GetEntity(id).(collisionPhase)
	oc := e.getCollisionPhase()

	if oc.OnCollisionS.Contains(LastEvent.ToSpace()) {
		if !oc.wasTouching {
			event.CID(id).Trigger(Start, LastEvent)
			oc.wasTouching = true
		}
	} else {
		if oc.wasTouching {
			event.CID(id).Trigger(Stop, LastEvent)
			oc.wasTouching = false
		}

	}
	return 0
}
