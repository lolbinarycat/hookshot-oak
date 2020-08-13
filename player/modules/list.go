package modules

import (
	"fmt"

	"github.com/oakmound/oak/v2/dlog"
)

type List map[string]Module

func (l *List) GiveAll(equip bool) {
	if len(*l) == 0 {
		dlog.Error("no modules to give")
	}
	for _, m := range *l {
		m.Obtain()
		if equip {
			m.Equip()
		}
	}
}

func (l List) ListModules() {
	for i, m := range l {
		fmt.Println(i,m)
	}
}
