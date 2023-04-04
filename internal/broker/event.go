package broker

type EventTopic interface {
	Validate() error
	String() string
}

type Event interface {
	Topic() EventTopic
	Value() []byte
	String() string
}
