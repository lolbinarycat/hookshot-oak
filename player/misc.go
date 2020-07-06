package player

import (
	"time"
	"math"

	dir "github.com/lolbinarycat/hookshot-oak/direction"
	"github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/dlog"
)

const JumpHeight int = 5
const WallJumpHeight float64 = 6
const WallJumpWidth float64 = 3
const ClimbSpeed float64 = 3

const (
	AirAccel    float64 = 0.4
	AirMaxSpeed float64 = 3
)
const Gravity float64 = 0.35

func (b *PhysObject) BlockUpdater(p *Player) func(){
	//b.Body.ApplyFriction(1)
	//b.Body.Delta.

	return func() {
		if p.HeldObj != b.Body {
		b.DoGravity()
		if b.ActiColls.GroundHit {
			b.Body.Delta.SetPos(0, 0)
		}
	}}
}

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

type SetStateOptArgs struct {
	SkipEnd bool
}

func (p *Player) SetStateAdv(state PlayerState, opt SetStateOptArgs) {
	defer func() {
		if r := recover(); r != nil {
			dlog.Error("error while setting state", r)
			p.State = state
		}
	}()

	if opt.SkipEnd == false {
		p.State.End(p)
	}
	p.StateStartTime = time.Now()

	p.State = state
	p.State.Start(p)
}

func (p *Player) SetState(state PlayerState) {
	p.SetStateAdv(state, SetStateOptArgs{})
}

func (p *Player) DoAirControls() {
	if oak.IsDown(currentControls.Left) && p.Body.Delta.X() > -AirMaxSpeed {
		// check to prevent inconsistant top speeds
		//(e.g. if you are half a AirAccel away from AirMaxSpeed)
		if p.Body.Delta.X()-AirAccel > -AirMaxSpeed {
			p.Body.Delta.ShiftX(-AirAccel)
		} else {
			p.Body.Delta.SetX(-AirMaxSpeed)
		}
	} else if oak.IsDown(currentControls.Right) && p.Body.Delta.X() < AirMaxSpeed {
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

func (o *PhysObject) DoGravity() {
	o.Body.Delta.ShiftY(Gravity)
}

func (o *PhysObject) DoCustomGravity(grav float64) {
	o.Body.Delta.ShiftY(grav)
}



func (p *Player) GrabObjRight(targetLabels ...collision.Label) (bool, event.CID) {
	return p.GrabObject(p.Body.W, p.Body.H, p.Body.W, targetLabels...)
}

func (p *Player) GrabObjLeft(targetLabels ...collision.Label) (bool, event.CID) {
	return p.GrabObject(-p.Body.W, -p.Body.H, p.Body.W, targetLabels...)
}

// GetLastHitObj attempts to get an entity from a PhysObject's
// ActiColls.LastHit* attribute. .LastHitH if Horis == true,
// and .LastHitV if false.
// it will return nil if unsucssesful.

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
	return iface.(*entities.Moving)
}


func (p *Player) DoStateLoop() {
	defer func () {
		recover()
	}()
	if p.TimeFromStateStart() > p.State.MaxDuration {
		p.SetState(*p.State.NextState)
		return
	}
	p.State.Loop(p)
}

func (o *PhysObject) IsWallHit() bool {
	if o.ActiColls.LeftWallHit || o.ActiColls.RightWallHit {
		return true
	}
	return false
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

//func (p *Player) SetEyePos() {
//
//}

//func AttachMut(a *vectAttach,attachTo physics.Attachable,offsets... float64) {
//	*a = a.Attach(attachTo,offsets...)
//}

func (o *PhysObject) HasHitInDir(d dir.Dir) bool {
	return (d.IsRight() && o.ActiColls.RightWallHit) ||
		(d.IsLeft() && o.ActiColls.LeftWallHit) ||
		(d.IsUp() && o.ActiColls.CeilingHit) ||
		(d.IsDown() && o.ActiColls.GroundHit)
}

func (p *Player) HeldDir() (d dir.Dir) {
	if oak.IsDown(p.Ctrls.Up) {
		d = d.Add(dir.MaxUp())
	}
	if oak.IsDown(p.Ctrls.Down) {
		d = d.Add(dir.MaxDown())
	}
	if oak.IsDown(p.Ctrls.Left) {
		d = d.Add(dir.MaxLeft())
	}
	if oak.IsDown(p.Ctrls.Right) {
		d = d.Add(dir.MaxRight())
	}
	return d
}


// Depreciated
func (p *Player) GrabObject(xOff, yOff, maxDist float64, targetLabels ...collision.Label) (bool, event.CID) {
	if len(targetLabels) > 1 {
		dlog.Error("muliple labels not implemented yet")
	}

	id, ent := event.ScanForEntity(func(e interface{}) bool {
		if ent, ok := e.(*entities.Moving); ok {

			if ent.Space.Label != targetLabels[0] {
				dlog.Verb("label check failed")
				return false
			}
			if !(ent.Space.CID == p.ActiColls.LastHitH) {
				dlog.Verb("id is equal. id:", ent.CID)
				return false
			}

			if ent.DistanceTo(p.Body.X()+xOff, p.Body.Y()+yOff) <=
				maxDist+(math.Max(ent.W, ent.H)) {

				dlog.Verb("distance condition fufilled")
				// if the entity has the correct label, and is within the max distance:
				return true
			}

			//dlog.Verb("d ==",d)
		} else {
			// if the entity is not a entities.Solid, we cannot grab it
			dlog.Verb("type check failed")
			return false
		}
		//this is just to stop "missing return at end of function"
		return false
	})

	// if id is equal to -1, it means ScanForEntity was unable
	// to find an entity within the given paramaters
	if id == -1 {
		dlog.Verb("ScanForEntity Failed")
		return false, -1
	}
	//p.HeldObjId = event.CID(id)
	if mov, ok := ent.(*entities.Moving); ok {
		p.HeldObj = &*mov
		dlog.Verb("HeldObj set")
	} else {
		dlog.Verb("ent exists, but is not *entities.Moving")
		return false, -1
	}

	return true, event.CID(id)
}
