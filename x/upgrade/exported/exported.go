package exported

// ProtocolVersionManager defines the interface which allows managing the appVersion field.
type ProtocolVersionManager interface {
	GetAppVersion() (uint64, error)
	SetAppVersion(version uint64) error
}
