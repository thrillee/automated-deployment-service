package systems

type RemoteALBManger interface {
	CreateContainerWithALB(ContainerPayload) ContainerResponse
	GetLoadBalancerDNS(map[string]interface{}) ContainerResponse
	GetCreateContainerWithALBKey() string
	GetLoadBalancerDNSKey() string
}
