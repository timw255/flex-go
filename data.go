package flex

import (
	"errors"
)

// ServiceObject ...
type ServiceObject interface {
	OnDeleteAll(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error
	OnDeleteByID(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error
	OnDeleteByQuery(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error
	OnGetAll(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error
	OnGetByID(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error
	OnGetByQuery(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error
	OnGetCount(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error
	OnGetCountByQuery(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error
	OnInsert(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error
	OnUpdate(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error
}

type serviceObject struct {
	name     string
	eventMap map[string]func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)
}

func (so *serviceObject) register(dataOp string, functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error {
	if dataOp == "" {
		return errors.New("Operation not permitted")
	}
	so.eventMap[dataOp] = functionToExecute
	return nil
}

func (so *serviceObject) unregister(dataOp string) error {
	if dataOp == "" {
		return errors.New("Operation not permitted")
	}
	delete(so.eventMap, dataOp)
	return nil
}

func (so *serviceObject) OnDeleteAll(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error {
	so.register("onDeleteAll", functionToExecute)
	return nil
}

func (so *serviceObject) OnDeleteByID(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error {
	so.register("onDeleteByID", functionToExecute)
	return nil
}

func (so *serviceObject) OnDeleteByQuery(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error {
	so.register("onDeleteByQuery", functionToExecute)
	return nil
}

func (so *serviceObject) OnGetAll(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error {
	so.register("onGetAll", functionToExecute)
	return nil
}

func (so *serviceObject) OnGetByID(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error {
	so.register("onGetByID", functionToExecute)
	return nil
}

func (so *serviceObject) OnGetByQuery(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error {
	so.register("onGetByQuery", functionToExecute)
	return nil
}

func (so *serviceObject) OnGetCount(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error {
	so.register("onGetCount", functionToExecute)
	return nil
}

func (so *serviceObject) OnGetCountByQuery(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error {
	so.register("onGetCountByQuery", functionToExecute)
	return nil
}

func (so *serviceObject) OnInsert(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error {
	so.register("onInsert", functionToExecute)
	return nil
}

func (so *serviceObject) OnUpdate(functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) error {
	so.register("onUpdate", functionToExecute)
	return nil
}

func (so *serviceObject) RemoveHandler(dataOp string) error {
	so.unregister(dataOp)
	return nil
}

func (so *serviceObject) resolve(dataOp string) func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task) {
	if i, ok := so.eventMap[dataOp]; ok {
		return i
	}
	return KinveyNotImplementedHandler()
}

// Data ...
type Data interface {
	NewServiceObject(name string) ServiceObject
	getServiceObjects() []string
	RemoveServiceObject(serviceObjectToRemove string) error
	clearAll()
	serviceObject(serviceObjectName string) serviceObject
	process(task *Task, modules Modules) (*Task, *Task)
}

type data struct {
	registeredServiceObjects map[string]serviceObject
}

func newData() Data {
	fd := &data{
		registeredServiceObjects: make(map[string]serviceObject),
	}
	return fd
}

// NewServiceObject ...
func (fd *data) NewServiceObject(name string) ServiceObject {
	so := fd.newServiceObject(name)
	return &so
}

func (fd *data) newServiceObject(name string) serviceObject {
	so := serviceObject{
		name: name,
	}
	so.eventMap = make(map[string]func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task))
	fd.registeredServiceObjects[name] = so
	return so
}

func (fd *data) getServiceObjects() []string {
	keys := make([]string, 0)
	for key := range fd.registeredServiceObjects {
		keys = append(keys, key)
	}
	return keys
}

func (fd *data) serviceObject(serviceObjectName string) serviceObject {
	if _, ok := fd.registeredServiceObjects[serviceObjectName]; !ok {
		fd.registeredServiceObjects[serviceObjectName] = fd.newServiceObject(serviceObjectName)
	}
	return fd.registeredServiceObjects[serviceObjectName]
}

func (fd *data) process(task *Task, modules Modules) (*Task, *Task) {
	serviceObjectToProcess := fd.serviceObject(task.Request.ServiceObjectName)

	var dataOp string
	if task.Method == "POST" {
		dataOp = "onInsert"
	} else if task.Method == "PUT" {
		dataOp = "onUpdate"
	} else if task.Method == "GET" && task.Endpoint != "_count" {
		taskRequest := task.Request
		if taskRequest.EntityID != "" {
			dataOp = "onGetByID"
		} else if len(taskRequest.Query) > 0 {
			dataOp = "onGetByQuery"
		} else {
			dataOp = "onGetAll"
		}
	} else if task.Method == "GET" && task.Endpoint == "_count" {
		taskRequest := task.Request
		if len(taskRequest.Query) > 0 {
			dataOp = "onGetCountByQuery"
		} else {
			dataOp = "onGetCount"
		}
	} else if task.Method == "DELETE" {
		taskRequest := task.Request
		if taskRequest.EntityID != "" {
			dataOp = "onDeleteByID"
		} else if len(taskRequest.Query) > 0 {
			dataOp = "onDeleteByQuery"
		} else {
			dataOp = "onDeleteAll"
		}
	} else {
		// 'BadRequest', 'Cannot determine data operation'
	}

	operationHandler := serviceObjectToProcess.resolve(dataOp)
	dataCompletionHandler := NewKinveyCompletionHandler(task)
	return operationHandler(&task.Request, dataCompletionHandler, modules)
}

// RemoveServiceObject ...
func (fd *data) RemoveServiceObject(serviceObjectToRemove string) error {
	if serviceObjectToRemove == "" {
		return errors.New("Must list ServiceObject name")
	}
	delete(fd.registeredServiceObjects, serviceObjectToRemove)
	return nil
}

func (fd *data) clearAll() {
	fd.registeredServiceObjects = make(map[string]serviceObject)
}
