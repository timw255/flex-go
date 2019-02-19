package flex

// BackendContextModule ...
type BackendContextModule struct {
	appMetadata kinveyAppMetadata
}

func newBackendContextModule(appMetadata kinveyAppMetadata) BackendContextModule {
	s := BackendContextModule{
		appMetadata: appMetadata,
	}
	return s
}

// GetAppKey ...
func (m BackendContextModule) GetAppKey() string {
	return m.appMetadata.ID
}

// GetAppSecret ...
func (m BackendContextModule) GetAppSecret() string {
	return m.appMetadata.AppSecret
}

// GetMasterSecret ...
func (m BackendContextModule) GetMasterSecret() string {
	return m.appMetadata.MasterSecret
}
