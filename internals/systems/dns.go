package systems

type DNSResponse struct {
	Success bool
	Message string
	Err     error
	Result  map[string]interface{}
}

type RemoteDNSManager interface {
	CreateAliasRecord(map[string]interface{}) DNSResponse
	UpdateDNSSettings(map[string]interface{}) DNSResponse
	GetCreateAliasRecordKey() string
	GetUpdateDNSSettingsKey() string
}
