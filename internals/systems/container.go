package systems

type SystemResponse struct {
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
	CreateContainer(ContainerPayload) SystemResponse
	UpdateContainer(map[string]interface{}) SystemResponse
	StartContainer(map[string]interface{}) SystemResponse
	SetupAutoScaling(map[string]interface{}) SystemResponse
	GetStartContainerKey() string
	GetCreateContainerKey() string
	GetUpdateContainerKey() string
	GetSetupAutoScalingKey() string
}
