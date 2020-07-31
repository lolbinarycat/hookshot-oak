package physobj

import (


	"github.com/lolbinarycat/hookshot-oak/labels"
	"github.com/oakmound/oak/v2/collision"
)

func (object *PhysObject) DoCollision(updater func()) {
	const stepThreshold = 8
	oldX, oldY := object.Body.GetPos()
	updater()
	object.ActiColls = ActiveCollisions{} //reset the struct to be all false

	object.Body.ShiftX(object.Body.Delta.X())
	hit := collision.HitLabel(object.Body.Space,
		append(labels.Solids, object.ExtraSolids...)...)
	if hit != nil {
		if object.Body.Delta.X() > 0 { //Right Wall
			object.ActiColls.RightWallHit = true
			object.Body.SetX(oldX)
		} else if object.Body.Delta.X() < 0 { //Left Wall
			object.ActiColls.LeftWallHit = true
			object.Body.SetX(oldX)
		} else {
			// we do this to prevent values being set when they shouldn't be
			goto CheckY
		}
		object.Body.Delta.SetX(0)
		object.ActiColls.HLabel = hit.Label
		object.ActiColls.LastHitH = hit.CID
	}
CheckY:
	object.Body.ShiftY(object.Body.Delta.Y())

	if hit := collision.HitLabel(object.Body.Space,
		append(object.ExtraSolids, labels.Solids...)...); hit != nil {

		if object.Body.Delta.Y() > 0 { //Ground
			object.Body.SetY(oldY)
			object.ActiColls.GroundHit = true
		} else if object.Body.Delta.Y() < 0 { //Ceiling
			object.ActiColls.CeilingHit = true
			object.Body.SetY(oldY)
		} else {
			return
		}
		object.Body.Delta.SetY(0)
		object.ActiColls.VLabel = hit.Label
		object.ActiColls.LastHitV = hit.CID
	}
	//}
	//EndYMovementStep:
}
