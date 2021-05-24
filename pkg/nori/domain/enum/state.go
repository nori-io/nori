package enum

type State uint8

const (
	Undefined State = iota
	None
	Inited
	Running
	Stopped
)

const (
	undefined = "undefined"
	none      = "none"
	inited    = "inited"
	running   = "running"
	stopped   = "stopped"
)

func (s State) String() string {
	switch s {
	case None:
		return none
	case Inited:
		return inited
	case Running:
		return running
	case Stopped:
		return stopped
	default:
		return undefined
	}
}

func (s State) Value() uint8 {
	return uint8(s)
}

func New(v uint8) State {
	if uint8(Undefined) <= v && v <= uint8(Stopped) {
		return State(v)
	}
	return Undefined
}
