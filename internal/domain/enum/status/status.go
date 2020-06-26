package status

type Status uint8

const (
	Nil Status = iota
	Initialised
	Started
	Stopped
)
