package systems

type ContainerResponse struct {
	Success bool
	Message string
	Err     error
	Result  map[string]interface{}
}

type ContainerPayload struct {
	SubscriberName string
	ContainerName  string
	ImageUri       string
	ContainerPort  int
	HostPort       int
	Props          map[string]interface{}
}

type RemoteContainerService interface {
	CreateContainer(ContainerPayload) ContainerResponse
	UpdateContainer(map[string]interface{}) ContainerResponse
	StartContainer(map[string]interface{}) ContainerResponse
	SetupAutoScaling(map[string]interface{}) ContainerResponse
	GetStartContainerKey() string
	GetCreateContainerKey() string
	GetUpdateContainerKey() string
	GetSetupAutoScalingKey() string
}
