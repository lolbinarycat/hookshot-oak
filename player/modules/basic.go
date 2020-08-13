package modules

type BasicMod struct {
	Equipped bool
	Obtained bool //`default:true`
}


func (l *List) AddBasic(key string) *List {
	(*l)[key] = Module(&BasicMod{})
	return l
}

//this is here to fufill the interface
func (m BasicMod) JustActivated() bool {
	return false
}

func (m *BasicMod) Obtain() {
	m.Obtained = true
}

func (m *BasicMod) Equip() {
	m.Equipped = true
}

func (m *BasicMod) Unequip() {
	m.Equipped = false
}

func (m BasicMod) Active() bool {
	return m.Equipped && m.Obtained
}

func (m *BasicMod) GetBasic() *BasicMod {
	return m
}
