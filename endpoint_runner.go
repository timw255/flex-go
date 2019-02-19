package flex

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// EndpointRunnerModule ...
type EndpointRunnerModule struct {
	appMetadata    kinveyAppMetadata
	requestContext RequestMetadata
	taskMetadata   TaskMetadata
}

func newEndpointRunnerModule(appMetadata kinveyAppMetadata, requestMetadata RequestMetadata, taskMetadata TaskMetadata) EndpointRunnerModule {
	return EndpointRunnerModule{
		appMetadata:    appMetadata,
		requestContext: requestMetadata,
		taskMetadata:   taskMetadata,
	}
}

// NewEndpointRunner ...
func (m EndpointRunnerModule) NewEndpointRunner(useUserContext bool) EndpointRunner {
	er := EndpointRunner{}

	er.appMetadata = m.appMetadata
	er.requestContext = m.requestContext
	er.taskMetadata = m.taskMetadata
	er.useBL = true
	er.useUserContext = useUserContext

	er.client = &http.Client{}

	return er
}

// EndpointRunner ...
type EndpointRunner struct {
	baseStore
}

// NewEndpoint ...
func (er EndpointRunner) NewEndpoint(endpointName string) Endpoint {
	e := endpoint{
		endpointRunner: er,
		endpointName:   endpointName,
	}
	return e
}

func (er EndpointRunner) buildEndpointRequest() (*http.Request, error) {
	baseRoute := "rpc"

	return er.buildKinveyRequest(baseRoute, "custom", false, er.useUserContext)
}

func (er EndpointRunner) makeEndpointRequest(req *http.Request, endpointName string) ([]byte, error) {
	if er.taskMetadata.ObjectName == endpointName && er.taskMetadata.HookType == "customEndpoint" {
		return nil, errors.New("Not Allowed")
	}

	return er.makeRequest(req)
}

// Endpoint ...
type Endpoint interface {
	Execute(body []byte) ([]byte, error)
}

type endpoint struct {
	endpointRunner EndpointRunner
	endpointName   string
}

func (e endpoint) Execute(body []byte) ([]byte, error) {
	requestOptions, err := e.endpointRunner.buildEndpointRequest()
	if err != nil {
		return nil, err
	}

	requestOptions.Method = "POST"

	u, err := url.Parse(requestOptions.URL.String() + e.endpointName)
	if err != nil {
		return nil, err
	}

	requestOptions.URL = u

	requestOptions.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return e.endpointRunner.makeEndpointRequest(requestOptions, e.endpointName)
}
