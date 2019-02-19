package flex

import (
	"net/http"
)

// UserStoreModule ...
type UserStoreModule struct {
	appMetadata    kinveyAppMetadata
	requestContext RequestMetadata
	taskMetadata   TaskMetadata
}

func newUserStoreModule(appMetadata kinveyAppMetadata, requestMetadata RequestMetadata, taskMetadata TaskMetadata) UserStoreModule {
	return UserStoreModule{
		appMetadata:    appMetadata,
		requestContext: requestMetadata,
		taskMetadata:   taskMetadata,
	}
}

// NewUserStore ...
func (m UserStoreModule) NewUserStore(useUserContext bool) UserStore {
	s := UserStore{
		baseRoute: "user",
	}
	s.appMetadata = m.appMetadata
	s.requestContext = m.requestContext
	s.taskMetadata = m.taskMetadata
	s.useBL = true
	s.useUserContext = useUserContext
	s.client = &http.Client{}
	return s
}

// UserStore ...
type UserStore struct {
	baseStore
	baseRoute string
}
