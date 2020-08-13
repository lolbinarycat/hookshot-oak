package modules

import (
	"fmt"
	"github.com/oakmound/oak/v2/joystick"
	"github.com/oakmound/oak/v2"
)

// Input refers to a keyboard key/controller button input.
// An empty string refers to an unbound (to a button/key) input
// Bound refers to whether a module is bound to this input.
type Input struct {
	Key string `json:"key"`
	Button string `json:"button"`
	Bound bool
	Ctrlr *joystick.Joystick
}

func NewInput(k string,ctlr *joystick.Joystick,b string) Input {
	return Input{k,b,false,ctlr}
}

func (i Input) IsDown() bool {
	if oak.IsDown(i.Key) {
		return true
	} else {
		
		if i.Ctrlr != nil {
			st, _ := i.Ctrlr.GetState()
			return st.Buttons[i.Button]
		}
		return false
	}
}

type InputList [8]Input

func (l InputList) String() string {
	var str string
	for i, v := range l {
		str += fmt.Sprintf("%d{key:%v,button:%v}\n",i,v.Key,v.Button)
	}
	return str
}
