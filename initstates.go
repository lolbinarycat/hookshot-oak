package main

// denil returns a modified version of a playerstate with nil functions replaced
// with empty functions, preventing segfaults from happening if they are called.
// It is designed to be called on a struct literal when setting a value
func (s PlayerState) denil() PlayerState {
	if s.Start == nil {
		s.Start = func(p *Player) {}
	}
	if s.Loop == nil {
		s.Loop = func(p *Player) {}
	}
	if s.End == nil {
		s.End = func(p *Player) {}
	}
	return s
}

// initStates is called at the start of main().
// this is to stop an initialization error.
func initStates() {
	AirState = PlayerState{
		Loop:AirStateLoop,
		Start:func(p *Player) {},
		End:func(p *Player) {},
	}.denil()
	GroundState = PlayerState{
		Loop: GroundStateLoop,
	}.denil()
}
