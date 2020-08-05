package ui

import "image/draw"

type ToggleableOptionList struct {
	*OptionList
	Active bool
}

func (tol *ToggleableOptionList) DrawOffset(buff draw.Image, xOff, yOff float64) {
	if tol.Active {
		tol.OptionList.DrawOffset(buff,xOff,yOff)
	}
}

func (tol *ToggleableOptionList) Draw(buff draw.Image) {
	tol.DrawOffset(buff, 0, 0)
}

func (tol *ToggleableOptionList) Toggle() {
	tol.Active = !tol.Active
}

type SubOptionList struct {
	ToggleableOptionList
	Name string
}

func NewSubOptionList(name string,x, y float64, opts ...*Option) *SubOptionList {
	sl := SubOptionList{
		ToggleableOptionList: ToggleableOptionList{
			OptionList: NewOptionList(x,y,opts...),
			Active: false,
		},
		Name: name,
	}
	return &sl
}

// BuildOption builds the option that should be included in the parent option list
func (sl *SubOptionList) BuildOption() *Option {
	return &Option{
		Name: sl.Name,
		Action: func () {
			if sl.Active == false {
				sl.Active = true
			} else {
				sl.ActivateSelected()
			}
		},
		Extras: &OptionExtras{
			OnCycle: func() bool {
				if sl.Active == false {
					// is sublist is inactive, cycle as normal
					return true
				}
				sl.Cycle()
				return false
			},
		},
	}
}
