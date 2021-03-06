package compose

// Engine is the interface to implement container systems
type Engine interface {
	ActiveContainers() (int, error)
	PullImage(image string) (string, error)
	RunContainer(message string, image string) error
	getContainerID(image string) (string, error)
}

// Container is the base container system
type Container struct {
	Engine Engine
}
