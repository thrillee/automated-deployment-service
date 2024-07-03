package systems

import (
	"context"

	"github.com/thrillee/automated-deployment-service/internals/processor"
)

type System struct {
	RemoteALBManger        RemoteALBManger
	RemoteDNSManager       RemoteDNSManager
	RemoteSSLManager       RemoteSSLManager
	RemoteContainerService RemoteContainerService
	SystemName             string
}

func convertResultToProcessorResponse(r SystemResponse) processor.ProcessorResponse {
	return processor.ProcessorResponse{
		Success: r.Success,
		Message: r.Message,
		Err:     r.Err,
		Result:  r.Result,
	}
}

func (s System) RegisterSystems(p *processor.ProcessorHandlerFactory) error {
	// Application Load Balancer Service
	var err error
	err = p.Register(s.RemoteALBManger.GetCreateContainerWithALBKey(), func(ctx context.Context, d interface{}) processor.ProcessorResponse {
		res := s.RemoteALBManger.CreateContainerWithALB(d.(ContainerPayload))
		return convertResultToProcessorResponse(res)
	})

	err = p.Register(s.RemoteALBManger.GetLoadBalancerDNSKey(), func(ctx context.Context, d interface{}) processor.ProcessorResponse {
		res := s.RemoteALBManger.GetLoadBalancerDNS(d.(map[string]interface{}))
		return convertResultToProcessorResponse(res)
	})

	// Container Service
	err = p.Register(s.RemoteContainerService.GetCreateContainerKey(), func(ctx context.Context, d interface{}) processor.ProcessorResponse {
		res := s.RemoteContainerService.CreateContainer(d.(ContainerPayload))
		return convertResultToProcessorResponse(res)
	})

	err = p.Register(s.RemoteContainerService.GetUpdateContainerKey(), func(ctx context.Context, d interface{}) processor.ProcessorResponse {
		res := s.RemoteContainerService.UpdateContainer(d.(map[string]interface{}))
		return convertResultToProcessorResponse(res)
	})

	err = p.Register(s.RemoteContainerService.GetStartContainerKey(), func(ctx context.Context, d interface{}) processor.ProcessorResponse {
		res := s.RemoteContainerService.StartContainer(d.(map[string]interface{}))
		return convertResultToProcessorResponse(res)
	})

	err = p.Register(s.RemoteContainerService.GetSetupAutoScalingKey(), func(ctx context.Context, d interface{}) processor.ProcessorResponse {
		res := s.RemoteContainerService.SetupAutoScaling(d.(map[string]interface{}))
		return convertResultToProcessorResponse(res)
	})

	// DNS Service
	err = p.Register(s.RemoteDNSManager.GetCreateAliasRecordKey(), func(ctx context.Context, d interface{}) processor.ProcessorResponse {
		res := s.RemoteDNSManager.CreateAliasRecord(d.(map[string]interface{}))
		return convertResultToProcessorResponse(res)
	})

	err = p.Register(s.RemoteDNSManager.GetUpdateDNSSettingsKey(), func(ctx context.Context, d interface{}) processor.ProcessorResponse {
		res := s.RemoteDNSManager.UpdateDNSSettings(d.(map[string]interface{}))
		return convertResultToProcessorResponse(res)
	})

	// SSL Service
	err = p.Register(s.RemoteSSLManager.GetRequestSSLCertificateKey(), func(ctx context.Context, d interface{}) processor.ProcessorResponse {
		res := s.RemoteSSLManager.RequestSSLCertificate(d.(map[string]interface{}))
		return convertResultToProcessorResponse(res)
	})

	err = p.Register(s.RemoteSSLManager.GetWaitSSLCertificateKey(), func(ctx context.Context, d interface{}) processor.ProcessorResponse {
		res := s.RemoteSSLManager.WaitSSLCertificate(d.(map[string]interface{}))
		return convertResultToProcessorResponse(res)
	})

	err = p.Register(s.RemoteSSLManager.GetCreateALBSSLCertKey(), func(ctx context.Context, d interface{}) processor.ProcessorResponse {
		res := s.RemoteSSLManager.CreateALBSSLCert(d.(map[string]interface{}))
		return convertResultToProcessorResponse(res)
	})

	return err
}
