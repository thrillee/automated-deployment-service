package systems

type RemoteDNSManager interface {
	CreateAliasRecord(map[string]interface{}) SystemResponse
	UpdateDNSSettings(map[string]interface{}) SystemResponse
	GetCreateAliasRecordKey() string
	GetUpdateDNSSettingsKey() string
}
