package event

type Event interface {
	Topic() string
	String() string
}
