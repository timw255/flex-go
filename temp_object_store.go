package flex

// TempObjectStoreModule ...
type TempObjectStoreModule struct {
	tempObjectStore map[string]string
}

func newTempObjectStoreModule() TempObjectStoreModule {
	return TempObjectStoreModule{
		tempObjectStore: make(map[string]string),
	}
}

// Set ...
func (m TempObjectStoreModule) Set(key string, value string) string {
	m.tempObjectStore[key] = value
	return m.tempObjectStore[key]
}

// Get ...
func (m TempObjectStoreModule) Get(key string) string {
	return m.tempObjectStore[key]
}

// GetAll ...
func (m TempObjectStoreModule) GetAll() map[string]string {
	return m.tempObjectStore
}
