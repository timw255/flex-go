package flex

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// DataStoreModule ...
type DataStoreModule struct {
	appMetadata    kinveyAppMetadata
	requestContext RequestMetadata
	taskMetadata   TaskMetadata
}

func newDataStoreModule(appMetadata kinveyAppMetadata, requestMetadata RequestMetadata, taskMetadata TaskMetadata) DataStoreModule {
	return DataStoreModule{
		appMetadata:    appMetadata,
		requestContext: requestMetadata,
		taskMetadata:   taskMetadata,
	}
}

// NewDataStore ...
func (m DataStoreModule) NewDataStore(useBL bool, useUserContext bool) DataStore {
	s := DataStore{}

	s.appMetadata = m.appMetadata
	s.requestContext = m.requestContext
	s.taskMetadata = m.taskMetadata
	s.useBL = useBL
	s.useUserContext = useUserContext

	s.client = &http.Client{}

	return s
}

type baseStore struct {
	appMetadata    kinveyAppMetadata
	requestContext RequestMetadata
	taskMetadata   TaskMetadata
	useBL          bool
	useUserContext bool

	client *http.Client
}

func (bs baseStore) makeRequest(req *http.Request) ([]byte, error) {
	resp, err := bs.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (bs baseStore) buildKinveyRequest(baseRoute string, collection string, useAppSecret bool, useUserContext bool) (*http.Request, error) {
	if baseRoute == "" {
		return nil, errors.New("Missing Base Route")
	}

	url := bs.appMetadata.BaaSURL + "/" + baseRoute + "/" + bs.appMetadata.ID + "/"

	if collection != "" {
		url += collection + "/"
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Kinvey-API-Version", bs.requestContext.APIVersion)

	if bs.useBL != true {
		req.Header.Set("x-kinvey-skip-business-logic", strconv.FormatBool(true))
	}

	if useAppSecret {
		req.SetBasicAuth(bs.appMetadata.ID, bs.appMetadata.AppSecret)
	} else if (bs.useUserContext && useUserContext) || useUserContext {
		req.Header.Add("authorization", bs.requestContext.Authorization)
	} else {
		req.SetBasicAuth(bs.appMetadata.ID, bs.appMetadata.MasterSecret)
	}

	return req, nil
}

// DataStore ...
type DataStore struct {
	baseStore
}

// NewCollection ...
func (s DataStore) NewCollection(collectionName string) Collection {
	return collection{
		dataStore:      s,
		collectionName: collectionName,
	}
}

func (s DataStore) buildAppDataRequest(collectionName string) (*http.Request, error) {
	baseRoute := "appdata"
	return s.buildKinveyRequest(baseRoute, collectionName, false, s.useUserContext)
}

func (s DataStore) makeAppDataRequest(req *http.Request, collectionName string) ([]byte, error) {
	if s.taskMetadata.ObjectName == collectionName && (s.useBL || s.useUserContext) {
		return nil, errors.New("Not Allowed")
	}
	return s.makeRequest(req)
}

// Entity ...
type Entity interface {
	GetID() *string
}

// Collection ...
type Collection interface {
	Find(query string) ([]byte, error)
	FindByID(id string) ([]byte, error)
	Save(entity Entity) ([]byte, error)
	Remove(query string) (int, error)
	RemoveByID(id string) (int, error)
	Count(query string) (int, error)
}

type collection struct {
	dataStore      DataStore
	collectionName string
}

type countResponse struct {
	Count int `json:"count"`
}

func (c collection) Find(query string) ([]byte, error) {
	requestOptions, err := c.dataStore.buildAppDataRequest(c.collectionName)
	if err != nil {
		return nil, err
	}

	requestOptions.Method = "GET"

	if query != "" {
		q := requestOptions.URL.Query()
		q.Add("query", query)
		requestOptions.URL.RawQuery = q.Encode()
	}

	return c.dataStore.makeAppDataRequest(requestOptions, c.collectionName)
}

func (c collection) FindByID(id string) ([]byte, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	requestOptions, err := c.dataStore.buildAppDataRequest(c.collectionName)
	if err != nil {
		return nil, err
	}

	requestOptions.Method = "GET"

	u, err := url.Parse(requestOptions.URL.String() + id)
	if err != nil {
		return nil, err
	}

	requestOptions.URL = u

	return c.dataStore.makeAppDataRequest(requestOptions, c.collectionName)
}

func (c collection) Save(entity Entity) ([]byte, error) {
	requestOptions, err := c.dataStore.buildAppDataRequest(c.collectionName)
	if err != nil {
		return nil, err
	}

	if entity.GetID() != nil {
		requestOptions.Method = "PUT"

		u, err := url.Parse(requestOptions.URL.String() + *entity.GetID())
		if err != nil {
			return nil, err
		}

		requestOptions.URL = u
	} else {
		requestOptions.Method = "POST"
	}

	json, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}

	requestOptions.Body = ioutil.NopCloser(bytes.NewReader(json))

	return c.dataStore.makeAppDataRequest(requestOptions, c.collectionName)
}

func (c collection) Remove(query string) (int, error) {
	requestOptions, err := c.dataStore.buildAppDataRequest(c.collectionName)
	if err != nil {
		return 0, err
	}

	requestOptions.Method = "DELETE"

	if query != "" {
		q := requestOptions.URL.Query()
		q.Add("query", query)
		requestOptions.URL.RawQuery = q.Encode()
	}

	kinveyResponse, err := c.dataStore.makeAppDataRequest(requestOptions, c.collectionName)
	if err != nil {
		return 0, err
	}

	countResponse := countResponse{}
	err = json.Unmarshal(kinveyResponse, &countResponse)
	if err != nil {
		return 0, err
	}

	return countResponse.Count, nil
}

func (c collection) RemoveByID(id string) (int, error) {
	if id == "" {
		return 0, errors.New("id is required")
	}

	requestOptions, err := c.dataStore.buildAppDataRequest(c.collectionName)
	if err != nil {
		return 0, err
	}

	requestOptions.Method = "DELETE"

	u, err := url.Parse(requestOptions.URL.String() + id)
	if err != nil {
		return 0, err
	}

	requestOptions.URL = u

	kinveyResponse, err := c.dataStore.makeAppDataRequest(requestOptions, c.collectionName)
	if err != nil {
		return 0, err
	}

	countResponse := countResponse{}
	err = json.Unmarshal(kinveyResponse, &countResponse)
	if err != nil {
		return 0, err
	}

	return countResponse.Count, nil
}

func (c collection) Count(query string) (int, error) {
	requestOptions, err := c.dataStore.buildAppDataRequest(c.collectionName)
	if err != nil {
		return 0, err
	}

	requestOptions.Method = "GET"

	u, err := url.Parse(requestOptions.URL.String() + "_count")
	if err != nil {
		return 0, err
	}

	requestOptions.URL = u

	if query != "" {
		q := requestOptions.URL.Query()
		q.Add("query", query)
		requestOptions.URL.RawQuery = q.Encode()
	}

	kinveyResponse, err := c.dataStore.makeAppDataRequest(requestOptions, c.collectionName)
	if err != nil {
		return 0, err
	}

	countResponse := countResponse{}
	err = json.Unmarshal(kinveyResponse, &countResponse)
	if err != nil {
		return 0, err
	}

	return countResponse.Count, nil
}
