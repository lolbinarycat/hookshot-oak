package oak

// ShinySend calls Send(event) on the underlying shiny.Window
func ShinySend(event interface{}) {
	windowControl.Send(event)
}

// ShinySendFirst is the same as ShinySend, but uses SendFirst instead of Send
func ShinySendFirst(event interface{}) {
	windownControl.SendFirst(event)
}
