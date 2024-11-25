package poteto

type PotetoOption struct {
	WithRequestId bool `yaml:"with_request_id"`
}

var DefaultPotetoOption = PotetoOption{
	WithRequestId: true,
}
