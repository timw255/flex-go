package flex

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// RoleStoreModule ...
type RoleStoreModule struct {
	appMetadata    kinveyAppMetadata
	requestContext RequestMetadata
	taskMetadata   TaskMetadata
}

func newRoleStoreModule(appMetadata kinveyAppMetadata, requestMetadata RequestMetadata, taskMetadata TaskMetadata) RoleStoreModule {
	return RoleStoreModule{
		appMetadata:    appMetadata,
		requestContext: requestMetadata,
		taskMetadata:   taskMetadata,
	}
}

// NewRoleStore ...
func (m RoleStoreModule) NewRoleStore(useUserContext bool) RoleStore {
	s := RoleStore{
		baseRoute: "roles",
	}
	s.appMetadata = m.appMetadata
	s.requestContext = m.requestContext
	s.taskMetadata = m.taskMetadata
	s.useBL = true
	s.useUserContext = useUserContext
	s.client = &http.Client{}
	return s
}

// RoleStore ...
type RoleStore struct {
	baseStore
	baseRoute string
}

// Role ...
type Role struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func (s RoleStore) buildRoleRequest() (*http.Request, error) {
	return s.buildKinveyRequest(s.baseRoute, "custom", false, s.useUserContext)
}

func (s RoleStore) makeRoleRequest(req *http.Request) ([]byte, error) {
	if s.taskMetadata.ObjectName == s.baseRoute && (s.useBL || s.useUserContext) {
		return nil, errors.New("Not Allowed")
	}
	return s.makeRequest(req)
}

// Create ...
func (s RoleStore) Create(role Role) ([]byte, error) {
	requestOptions, err := s.buildRoleRequest()
	if err != nil {
		return nil, err
	}
	requestOptions.Method = "POST"

	json, err := json.Marshal(role)
	requestOptions.Body = ioutil.NopCloser(bytes.NewBuffer(json))

	return s.makeRoleRequest(requestOptions)
}

// Update ...
func (s RoleStore) Update(role Role) ([]byte, error) {
	requestOptions, err := s.buildRoleRequest()
	if err != nil {
		return nil, err
	}

	if role.ID == "" {
		return nil, errors.New("ID required")
	}

	requestOptions.Method = "PUT"

	u, err := url.Parse(requestOptions.URL.String() + role.ID)
	if err != nil {
		return nil, err
	}

	requestOptions.URL = u

	json, err := json.Marshal(role)

	requestOptions.Body = ioutil.NopCloser(bytes.NewReader(json))

	return s.makeRoleRequest(requestOptions)
}

// FindByID ...
func (s RoleStore) FindByID(id string) ([]byte, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	requestOptions, err := s.buildRoleRequest()
	if err != nil {
		return nil, err
	}

	requestOptions.Method = "GET"

	u, err := url.Parse(requestOptions.URL.String() + id)
	if err != nil {
		return nil, err
	}

	requestOptions.URL = u

	return s.makeRoleRequest(requestOptions)
}

// Remove ...
func (s RoleStore) Remove(id string) ([]byte, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	requestOptions, err := s.buildRoleRequest()
	if err != nil {
		return nil, err
	}

	requestOptions.Method = "DELETE"

	u, err := url.Parse(requestOptions.URL.String() + id)
	if err != nil {
		return nil, err
	}

	requestOptions.URL = u

	return s.makeRoleRequest(requestOptions)
}

// ListMembers ...
func (s RoleStore) ListMembers(id string) ([]byte, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	requestOptions, err := s.buildRoleRequest()
	if err != nil {
		return nil, err
	}

	requestOptions.Method = "GET"

	u, err := url.Parse(requestOptions.URL.String() + id + "/membership")
	if err != nil {
		return nil, err
	}

	requestOptions.URL = u

	return s.makeRoleRequest(requestOptions)
}
