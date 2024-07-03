package systems

type SSLResponse struct {
	Success bool
	Message string
	Err     error
	Result  map[string]interface{}
}

type RemoteSSLManager interface {
	RequestSSLCertificate(map[string]interface{}) SSLResponse
	WaitSSLCertificate(map[string]interface{}) SSLResponse
	CreateALBSSLCert(map[string]interface{}) SSLResponse

	GetRequestSSLCertificateKey() string
	GetWaitSSLCertificateKey() string
	GetCreateALBSSLCertKey() string
}
