package flex

import (
	"net/url"
)

type kinveyOriginalRequestHeaders struct {
	Authorization          string `json:"authorization"`
	KinveyAPIVersion       string `json:"x-kinvey-api-version"`
	KinveyClientAppVersion string `json:"x-kinvey-client-app-version"`
}

type kinveyAppMetadata struct {
	ID            string                 `json:"_id"`
	ApplicationID string                 `json:"applicationId"`
	AppSecret     string                 `json:"appsecret"`
	BaaSURL       string                 `json:"bassUrl"`
	BLFlags       map[string]interface{} `json:"blFlags"`
	Maintenance   maintenance            `json:"maintenance"`
	MasterSecret  string                 `json:"mastersecret"`
	Name          string                 `json:"name"`
}

type maintenance struct {
	ObjectIDMigration objectIDMigration `json:"objectid_migration"`
}

type objectIDMigration struct {
	Status string `json:"status"`
}

// Task ...
type Task struct {
	AppID            string            `json:"appId"`
	AppMetadata      kinveyAppMetadata `json:"appMetadata"`
	AuthKey          string            `json:"authKey"`
	BaaSURL          string            `json:"baasUrl"`
	ContainerID      string            `json:"containerId"`
	DiscoveryObjects discoveryObjects  `json:"discoveryObjects"`
	Endpoint         string
	HookType         string   `json:"hookType"`
	Method           string   `json:"method"`
	ObjectName       string   `json:"objectName"`
	Request          Request  `json:"request"`
	RequestID        string   `json:"requestId"`
	Response         Response `json:"response"`
	SDKVersion       string
	Target           string `json:"target"`
	TaskID           string `json:"taskId"`
	TaskName         string `json:"taskName"`
	TaskType         string `json:"taskType"`
}

type discoveryObjects struct {
	Auth          authDiscovery `json:"auth"`
	BusinessLogic businessLogic `json:"businessLogic"`
	DataLink      dataLink      `json:"dataLink"`
}

type dataLink struct {
	ServiceObjects []string `json:"serviceObjects"`
}

type businessLogic struct {
	Handlers []string `json:"handlers"`
}

type authDiscovery struct {
	Handlers []string `json:"handlers"`
}

type netType interface {
	GetHeaders() map[string]string
	GetBody() string
}

// RequestMetadata ...
type RequestMetadata struct {
	AuthenticatedUsername   string                 `json:"authenticatedUsername"`
	AuthenticatedUserID     string                 `json:"authenticatedUserId"`
	APIVersion              string                 `json:"apiVersion"`
	Authorization           string                 `json:"authorization"`
	ClientAppVersion        string                 `json:"clientAppVersion"`
	CustomRequestProperties map[string]interface{} `json:"customRequestProperties"`
	RequestID               string                 `json:"requestId"`
	SecurityContext         string                 `json:"securityContext"`
}

// TaskMetadata ...
type TaskMetadata struct {
	ContainerID string `json:"containerId"`
	HookType    string `json:"hookType"`
	ObjectName  string `json:"objectName"`
	Target      string `json:"target"`
	TaskID      string `json:"taskId"`
	TaskType    string `json:"taskType"`
}

// Request ...
type Request struct {
	Body              []byte                 `json:"-"`
	JSONBody          map[string]interface{} `json:"body"`
	CollectionName    string                 `json:"collectionName"`
	EntityID          string
	Headers           map[string]string `json:"headers"`
	HookType          string
	Method            string `json:"method"`
	ObjectName        string
	ServiceObjectName string
	TempObjectStore   map[string]interface{} `json:"tempObjectStore"`
	UserID            string
	Username          string `json:"username"`
	Query             url.Values
}

// GetHeaders ...
func (r *Request) GetHeaders() map[string]string {
	return r.Headers
}

// GetBody ...
func (r *Request) GetBody() string {
	json, _ := json.Marshal(r.Body)

	return string(json)
}

// Response ...
type Response struct {
	Body     []byte                 `json:"-"`
	JSONBody map[string]interface{} `json:"body"`
	Headers  map[string]string      `json:"headers"`
	HookType string
	Status   int `json:"status"`
}

// GetHeaders ...
func (r *Response) GetHeaders() map[string]string {
	return r.Headers
}

// GetBody ...
func (r *Response) GetBody() string {
	return string(r.Body)
}
