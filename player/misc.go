package player

import (
	"time"
	//"math"

	dir "github.com/lolbinarycat/hookshot-oak/direction"
	"github.com/lolbinarycat/hookshot-oak/fginput"
	"github.com/oakmound/oak/v2"
	//"github.com/oakmound/oak/v2/event"
	//"github.com/oakmound/oak/v2/collision"
	//"github.com/oakmound/oak/v2/entities"
	//"github.com/oakmound/oak/v2/dlog"
)

const JumpHeight int = 5
const WallJumpHeight float64 = 6
const WallJumpWidth float64 = 3
const ClimbSpeed float64 = 3

const (
	AirAccel    float64 = 0.4
	AirMaxSpeed float64 = 3
)


func (p *Player) WallJump(dir dir.Dir, EnterLaunch bool) {
	p.Body.Delta.SetY(-WallJumpHeight)

	if dir.IsLeft() {
		p.Body.Delta.SetX(-WallJumpWidth)
	} else if dir.IsRight() {
		p.Body.Delta.SetX(WallJumpWidth)
	} else {
		panic("invalid direction to WallJump functon")
	}

	if EnterLaunch {
		p.SetState(WallJumpLaunchState)
	} else {
		p.SetState(AirState)
	}
}

// DoCliming is the function for shared procceses between
// ClimbRightState and ClimbLeft state
func (p *Player) DoCliming() {
	// this is a hack, and should problebly be fixed
	if int(p.TimeFromStateStart())%2 == 0 && p.ActiColls.LeftWallHit == false && p.ActiColls.RightWallHit == false {
		p.SetState(AirState)
	}
	if p.Mods["jump"].JustActivated() {
		p.SetState(AirState)
	}
	if oak.IsDown(p.Ctrls.Up) {
		p.Body.Delta.SetY(-ClimbSpeed)
	} else if oak.IsDown(p.Ctrls.Down) {
		p.Body.Delta.SetY(ClimbSpeed)
	} else {
		p.Body.Delta.SetY(0)
	}

	p.StateCommon()
}

func (p *Player) IsJumpInput() bool {
	return p.Mods["jump"].JustActivated()
}

func (p *Player) IsButtonPressedWithin(button string, dur time.Duration) bool {
	if k, d := oak.IsHeld(button); k && (d <= dur) {
		return true
	} else {
		return false
	}
}

func (p *Player) IsHsInput() bool {
	return p.Mods["hs"].JustActivated()
}

func (p *Player) IfHsPressedStartHs() {
	if p.IsHsInput() {
		p.SetState(HsStartState)
	}
}

func (p *Player) Jump() {
	p.Body.Delta.ShiftY(-p.Body.Speed.Y())
	p.Body.ShiftY(p.Body.Delta.Y())
	p.SetState(JumpHeightDecState)
}



func (p *Player) DoAirControls() {
	if p.HeldDir.IsLeft() && p.Body.Delta.X() > -AirMaxSpeed {
		// check to prevent inconsistant top speeds
		//(e.g. if you are half a AirAccel away from AirMaxSpeed)
		if p.Body.Delta.X()-AirAccel > -AirMaxSpeed {
			p.Body.Delta.ShiftX(-AirAccel)
		} else {
			p.Body.Delta.SetX(-AirMaxSpeed)
		}
	} else if p.HeldDir.IsRight() && p.Body.Delta.X() < AirMaxSpeed {
		//second verse, same as the first
		if p.Body.Delta.X()+AirAccel < AirMaxSpeed {
			p.Body.Delta.ShiftX(AirAccel)
		} else {
			p.Body.Delta.SetX(AirMaxSpeed)
		}
	}
}

//TimeFromStateStart gets how long it has been since the last state transition
func (p *Player) TimeFromStateStart() time.Duration {
	return time.Now().Sub(p.StateStartTime)
}

func (p *Player) Die() {
	//TODO: death animation
	p.Respawn()
}

func (p *Player) Respawn() {
	p.SetState(RespawnFallState)
	p.Body.Delta.SetPos(0, 0)
	p.Body.SetPos(p.RespawnPos.X, p.RespawnPos.Y)
}




// func (p *Player) GrabObjRight(targetLabels ...collision.Label) (bool, event.CID) {
// 	return p.GrabObject(p.Body.W, p.Body.H, p.Body.W, targetLabels...)
// }

// func (p *Player) GrabObjLeft(targetLabels ...collision.Label) (bool, event.CID) {
// 	return p.GrabObject(-p.Body.W, -p.Body.H, p.Body.W, targetLabels...)
// }

// GetLastHitObj attempts to get an entity from a PhysObject's
// ActiColls.LastHit* attribute. .LastHitH if Horis == true,
// and .LastHitV if false.
// it will return nil if unsucssesful.




func (p *Player) DoStateLoop() {
	
	// if p.TimeFromStateStart() > p.State.MaxDuration {
	// 	p.SetState(*p.State.NextState)
	// 	return
	// }
	// if DoMap returns true, it means that the state changed.
	// if p.State.DoMap(p) == false {
	// 	nextState := p.State.LLoop(p)
	// 	if nextState == nil {
	// 		p.State.Loop(p)
	// 	} else {
	// 		p.SetState(*nextState)
	// 	}
	// }
	var nextState *PlayerState = nil
	if p.State.LLoop != nil {
		nextState = p.State.LLoop(p)
	}
	if nextState != nil {
		p.SetState(*nextState)
	} else {
		if p.State.Loop != nil {
			p.State.Loop(p)
		}
		p.FramesInState++
	}
}



// Depreciated
func (p *Player) IsHsInPlayer() bool {
	xover, yover := p.Hs.Body.Space.Overlap(p.Body.Space)
	if xover >= p.Hs.Body.W || yover >= p.Hs.Body.H {
		return true
	}
	return false
}

func (p *Player) DoHsCheck() bool {
	if p.IsHsInPlayer() || p.IsWallHit() || p.Hs.X <= 0 {
		p.EndHs()
		return true
	}
	return false
}

func (p *Player) Seq(s string) bool {
	p.DirBuffer.Check(fginput.Seq(s))
}

//func (p *Player) SetEyePos() {
//
//}

//func AttachMut(a *vectAttach,attachTo physics.Attachable,offsets... float64) {
//	*a = a.Attach(attachTo,offsets...)
//}




// Depreciated
// func (p *Player) GrabObject(xOff, yOff, maxDist float64, targetLabels ...collision.Label) (bool, event.CID) {
// 	if len(targetLabels) > 1 {
// 		dlog.Error("muliple labels not implemented yet")
// 	}

// 	id, ent := event.ScanForEntity(func(e interface{}) bool {
// 		if ent, ok := e.(*entities.Moving); ok {

// 			if ent.Space.Label != targetLabels[0] {
// 				dlog.Verb("label check failed")
// 				return false
// 			}
// 			if !(ent.Space.CID == p.ActiColls.LastHitH) {
// 				dlog.Verb("id is equal. id:", ent.CID)
// 				return false
// 			}

// 			if ent.DistanceTo(p.Body.X()+xOff, p.Body.Y()+yOff) <=
// 				maxDist+(math.Max(ent.W, ent.H)) {

// 				dlog.Verb("distance condition fufilled")
// 				// if the entity has the correct label, and is within the max distance:
// 				return true
// 			}

// 			//dlog.Verb("d ==",d)
// 		} else {
// 			// if the entity is not a entities.Solid, we cannot grab it
// 			dlog.Verb("type check failed")
// 			return false
// 		}
// 		//this is just to stop "missing return at end of function"
// 		return false
// 	})

// 	// if id is equal to -1, it means ScanForEntity was unable
// 	// to find an entity within the given paramaters
// 	if id == -1 {
// 		dlog.Verb("ScanForEntity Failed")
// 		return false, -1
// 	}
// 	//p.HeldObjId = event.CID(id)
// 	if mov, ok := ent.(*entities.Moving); ok {
// 		p.HeldObj = &*mov
// 		event.Trigger("holdUpdate",true)
// 		dlog.Verb("HeldObj set")
// 	} else {
// 		dlog.Verb("ent exists, but is not *entities.Moving")
// 		return false, -1
// 	}

// 	return true, event.CID(id)
// }
