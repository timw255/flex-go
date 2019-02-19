package flex

import (
	"net/http"
)

// GroupStoreModule ...
type GroupStoreModule struct {
	appMetadata    kinveyAppMetadata
	requestContext RequestMetadata
	taskMetadata   TaskMetadata
}

func newGroupStoreModule(appMetadata kinveyAppMetadata, requestMetadata RequestMetadata, taskMetadata TaskMetadata) GroupStoreModule {
	return GroupStoreModule{
		appMetadata:    appMetadata,
		requestContext: requestMetadata,
		taskMetadata:   taskMetadata,
	}
}

// NewGroupStore ...
func (m GroupStoreModule) NewGroupStore(useUserContext bool) GroupStore {
	s := GroupStore{
		baseRoute: "group",
	}

	s.appMetadata = m.appMetadata
	s.requestContext = m.requestContext
	s.taskMetadata = m.taskMetadata
	s.useBL = true
	s.useUserContext = useUserContext

	s.client = &http.Client{}

	return s
}

// GroupStore ...
type GroupStore struct {
	baseStore
	baseRoute string
}
