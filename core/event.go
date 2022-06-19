package core

type Event interface {
	Name() EventName
	Data() interface{}
}
