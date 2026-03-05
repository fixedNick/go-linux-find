package core

type Decision struct {
	Match   bool
	Actions []Action
	Control ControlSignal
}
