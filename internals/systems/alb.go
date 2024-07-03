package systems

type RemoteALBManger interface {
	CreateContainerWithALB(ContainerPayload) SystemResponse
	GetLoadBalancerDNS(map[string]interface{}) SystemResponse
	GetCreateContainerWithALBKey() string
	GetLoadBalancerDNSKey() string
}
