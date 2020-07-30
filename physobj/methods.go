package physobj

import (
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/event"
	"github.com/lolbinarycat/hookshot-oak/direction"
)

const Gravity float64 = 0.35

func (o *PhysObject) DoGravity() {
	o.Body.Delta.ShiftY(Gravity)
}

func (o *PhysObject) DoCustomGravity(grav float64) {
	o.Body.Delta.ShiftY(grav)
}

func (o *PhysObject) GetLastHitObj(Horis bool) *entities.Moving {
	_, iface := event.ScanForEntity(func(ent interface{}) bool {
		mov, ok := ent.(*entities.Moving)
		if !ok {
			return false
		}
		if (Horis && mov.Space.CID == o.ActiColls.LastHitH) ||
			(!Horis && mov.Space.CID == o.ActiColls.LastHitV) {
			return true
		}
		return false
	})
	if iface == nil {
		return nil
	}
	return iface.(*entities.Moving)
}

func (o *PhysObject) IsWallHit() bool {
	if o.ActiColls.LeftWallHit || o.ActiColls.RightWallHit {
		return true
	}
	return false
}

func (o *PhysObject) HasHitInDir(d direction.Dir) bool {
	return (d.IsRight() && o.ActiColls.RightWallHit) ||
		(d.IsLeft() && o.ActiColls.LeftWallHit) ||
		(d.IsUp() && o.ActiColls.CeilingHit) ||
		(d.IsDown() && o.ActiColls.GroundHit)
}
