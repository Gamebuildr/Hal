package compose

// Engine is the interface to implement container systems
type Engine interface {
	ActiveContainers() (int, error)
	RunContainer(image string) error
	getContainerID(image string) (string, error)
}

// ContainerData is the data the container system expects
// in a POST request
type ContainerData struct {
	Image string `json:"image"`
}

// Container is the base container system
type Container struct {
	Engine Engine
}
