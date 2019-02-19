# Flex Go

Flex Go is a custom version of the Kinvey Flex SDK, written in Go.

## This is unofficial, unsupported, and does not have all the functionality of the official SDK.

That said, it does actually work... and it's really, really fast. :)

If you're in to this sort of thing, pull requests of all sorts (tests, benchmarks, refactors, etc.) are accepted and absolutely welcome!

# Usage

```go
package main

import (
	"github.com/timw255/flex-go"
	"github.com/timw255/flex-go-example/handler"
)

func main() {
	options := flex.NewOptions("localhost", 10001, "")
	flex.NewService(options, func(err error, f flex.Flex) {
		if err != nil {
			f.Logger.Error("Error initializing the Flex SDK, exiting.")
		}

		// functions
		f.Functions.Register("myFunction", handler.MyFunctionHandler)

		// data
		widgets := f.Data.NewServiceObject("widgets")
		widgets.OnDeleteAll(handler.OnDeleteAll)
		widgets.OnDeleteByID(handler.OnDeleteByID)
		widgets.OnDeleteByQuery(handler.OnDeleteByQuery)
		widgets.OnGetAll(handler.OnGetAll)
		widgets.OnGetByID(handler.OnGetByID)
		widgets.OnGetByQuery(handler.OnGetByQuery)
		widgets.OnGetCount(handler.OnGetCount)
		widgets.OnGetCountByQuery(handler.OnGetCountByQuery)
		widgets.OnInsert(handler.OnInsert)
		widgets.OnUpdate(handler.OnUpdate)

		// auth
		f.Auth.Register("myAuth", handler.MyAuthHandler)
	})
}

```

# Flex Auth

```go
package handler

import (
	"github.com/timw255/flex-go"
)

// MyAuthHandler ...
func MyAuthHandler(context *flex.Request, complete flex.AuthCompletionHandler, modules flex.Modules) (*flex.Task, *flex.Task) {
	return complete.SetToken("41jknt32ntl34h234bthj3b24t").OK().Done()
}

```

# Flex Data

```go
package handler

import (
	"encoding/json"
	"time"

	"github.com/timw255/flex-go"
)

// CustomEntity ...
type CustomEntity struct {
	flex.KinveyEntity
	IsActive *bool   `json:"isActive,omitempty"`
	Key      *string `json:"key,omitempty"`
}

// OnGetAll ...
func OnGetAll(context *flex.Request, complete flex.KinveyCompletionHandler, modules flex.Modules) (*flex.Task, *flex.Task) {
	dataStore := modules.DataStore.NewDataStore(true, true)
	objectsCollection := dataStore.NewCollection("objects")

	data, err := objectsCollection.Find("")
	if err != nil {
		modules.Logger.Error(err.Error())
	}

	result := make([]CustomEntity, 0)
	if err := json.Unmarshal(data, &result); err != nil {
		modules.Logger.Error(err.Error())
	}

	json, err := json.Marshal(result)

	return complete.SetBody(json).Done()
}

// OnInsert ...
func OnInsert(context *flex.Request, complete flex.KinveyCompletionHandler, modules flex.Modules) (*flex.Task, *flex.Task) {
	dataStore := modules.DataStore.NewDataStore(true, true)
	objectsCollection := dataStore.NewCollection("objects")

	entity := CustomEntity{
		IsActive: flex.Bool(true),
		Key:      flex.String("some key"),
	}

	entity.KinveyEntity = modules.KinveyEntity.NewKinveyEntity("")

	entity.ACL.AddReaderRole("0f350bba-1145-e342-cb5a-223f314b650d")

	data, err := objectsCollection.Save(entity)
	if err != nil {
		modules.Logger.Error(err.Error())
	}

	newEntity := CustomEntity{}
	if err := json.Unmarshal(data, &newEntity); err != nil {
		modules.Logger.Error(err.Error())
	}

	json, err := json.Marshal(newEntity)

	return complete.SetBody(json).Done()
}

```

# Flex Functions

```go
package handler

import (
	"encoding/json"

	"github.com/timw255/flex-go"
)

type endpointMessage struct {
	Message string `json:"message"`
}

// MyFunctionHandler ...
func MyFunctionHandler(context *flex.Request, complete flex.KinveyCompletionHandler, modules flex.Modules) (*flex.Task, *flex.Task) {
	endpointRunner := modules.EndpointRunner.NewEndpointRunner(true)
	testEndpoint := endpointRunner.NewEndpoint("test")

	requestMessage := endpointMessage{
		Message: "ping!",
	}

	requestData, err := json.Marshal(requestMessage)
	if err != nil {
		modules.Logger.Error(err.Error())
	}

	responseData, err := testEndpoint.Execute(requestData)
	if err != nil {
		modules.Logger.Error(err.Error())
	}

	return complete.SetBody(responseData).Done()
}

```