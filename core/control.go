package core

type ControlSignal int

const (
	ControlNone ControlSignal = iota
	ControlPrune
	ControlQuit
)

func MergeControl(c1, c2 ControlSignal) ControlSignal {
	if c1 == ControlQuit || c2 == ControlQuit {
		return ControlQuit
	}

	if c1 == ControlPrune || c2 == ControlPrune {
		return ControlPrune
	}

	return ControlNone
}
