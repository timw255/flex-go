package flex

import (
	"encoding/base64"
	"strings"
)

// Modules ...
type Modules struct {
	BackendContext  BackendContextModule
	DataStore       DataStoreModule
	Email           EmailModule
	EndpointRunner  EndpointRunnerModule
	KinveyEntity    KinveyEntityModule
	RoleStore       RoleStoreModule
	TempObjectStore TempObjectStoreModule
	UserStore       UserStoreModule
	GroupStore      GroupStoreModule
	Logger          Logger
	//kinveyDate,
	//push: push(appMetadata),
	//Query,
	//requestContext: requestContext(requestMetadata),
}

func getSecurityContextString(authorizationHeader string, appMetadata kinveyAppMetadata) string {
	encodedCredentials := strings.Split(authorizationHeader, " ")

	if len(encodedCredentials) == 2 {
		decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials[1])
		if err != nil {
			return "unknown"
		}

		credentials := strings.Split(string(decodedCredentials), ":")

		if len(credentials) == 2 {
			if credentials[1] != appMetadata.ID {
				return "user"
			}
			if credentials[1] == appMetadata.ID && credentials[2] != "" {
				return "app"
			}
			if credentials[2] == appMetadata.MasterSecret {
				return "master"
			}
		}
	}

	return "unknown"
}

func generateModules(task *Task) Modules {
	var clientAppVersion string
	var customRequestProperties map[string]interface{}

	var baasURL string

	if task.BaaSURL != "" {
		baasURL = task.BaaSURL
	} else if task.AppMetadata.BaaSURL != "" {
		baasURL = task.AppMetadata.BaaSURL
	} else {
		forwardedProto := task.Request.Headers["x-forwarded-proto"]

		if forwardedProto == "" {
			forwardedProto = "https"
		}

		host := task.Request.Headers["host"]
		baasURL = forwardedProto + "://" + host
	}

	for _, key := range task.Request.Headers {
		if strings.ToLower(task.Request.Headers[key]) == "x-kinvey-client-app-version" {
			clientAppVersion = task.Request.Headers[key]
		}
		if strings.ToLower(task.Request.Headers[key]) == "x-kinvey-custom-request-properties" {
			err := json.UnmarshalFromString(task.Request.Headers[key], &customRequestProperties)
			if err != nil {
				customRequestProperties = make(map[string]interface{})
			}
		}
	}

	appMetadata := kinveyAppMetadata{
		ID:            task.AppMetadata.ID,
		ApplicationID: task.AppMetadata.ApplicationID,
		BLFlags:       task.AppMetadata.BLFlags,
		AppSecret:     task.AppMetadata.AppSecret,
		MasterSecret:  task.AppMetadata.MasterSecret,
		BaaSURL:       baasURL,
	}

	requestMetadata := RequestMetadata{
		AuthenticatedUsername:   task.Request.Username,
		AuthenticatedUserID:     task.Request.UserID,
		Authorization:           task.Request.Headers["authorization"],
		ClientAppVersion:        clientAppVersion,
		CustomRequestProperties: customRequestProperties,
		RequestID:               task.RequestID,
	}

	if task.Response.Headers["x-kinvey-api-version"] != "" {
		requestMetadata.APIVersion = task.Response.Headers["x-kinvey-api-version"]
	} else {
		requestMetadata.APIVersion = "3"
	}

	taskMetadata := TaskMetadata{
		TaskType:    task.TaskType,
		HookType:    task.HookType,
		Target:      task.Target,
		TaskID:      task.TaskID,
		ContainerID: task.ContainerID,
	}

	if task.Request.ServiceObjectName != "" {
		task.ObjectName = task.Request.ServiceObjectName
	} else if task.Request.ObjectName != "" {
		task.ObjectName = task.Request.ObjectName
	} else {
		task.ObjectName = task.Request.CollectionName
	}

	requestMetadata.SecurityContext = getSecurityContextString(task.Request.Headers["authorization"], appMetadata)

	useBSONObjectID := task.AppMetadata.Maintenance.ObjectIDMigration.Status != "done"

	return Modules{
		DataStore:       newDataStoreModule(appMetadata, requestMetadata, taskMetadata),
		EndpointRunner:  newEndpointRunnerModule(appMetadata, requestMetadata, taskMetadata),
		TempObjectStore: newTempObjectStoreModule(),
		KinveyEntity:    newKinveyEntityModule(appMetadata.ID, useBSONObjectID),
		BackendContext:  newBackendContextModule(appMetadata),
		RoleStore:       newRoleStoreModule(appMetadata, requestMetadata, taskMetadata),
		Email:           newEmailModule(appMetadata),
		UserStore:       newUserStoreModule(appMetadata, requestMetadata, taskMetadata),
		GroupStore:      newGroupStoreModule(appMetadata, requestMetadata, taskMetadata),
		Logger:          newLogger(),
	}
}
