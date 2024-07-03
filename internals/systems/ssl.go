package systems

type RemoteSSLManager interface {
	RequestSSLCertificate(map[string]interface{}) SystemResponse
	WaitSSLCertificate(map[string]interface{}) SystemResponse
	CreateALBSSLCert(map[string]interface{}) SystemResponse

	GetRequestSSLCertificateKey() string
	GetWaitSSLCertificateKey() string
	GetCreateALBSSLCertKey() string
}
