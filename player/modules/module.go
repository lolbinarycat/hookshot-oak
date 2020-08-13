package modules

type Module interface{
	Equip()
	Unequip()
	Obtain()
	Active() bool
	JustActivated() bool
	GetBasic() *BasicMod
}
