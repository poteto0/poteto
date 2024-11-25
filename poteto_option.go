package poteto

type PotetoOption struct {
	WithRequestId   bool   `yaml:"with_request_id"`
	ListenerNetwork string `yaml:"listener_network"`
}

var DefaultPotetoOption = PotetoOption{
	WithRequestId:   true,
	ListenerNetwork: "tcp",
}
