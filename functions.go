package flex

// Functions ...
type Functions interface {
	getHandlers() []string
	clearAll()
	resolve(taskName string) func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)
	process(task *Task, modules Modules) (*Task, *Task)
	Register(taskName string, functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task))
}

type functions struct {
	registeredFunctions map[string]func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)
}

func (ff *functions) getHandlers() []string {
	keys := make([]string, 0)
	for key := range ff.registeredFunctions {
		keys = append(keys, key)
	}
	return keys
}

func (ff *functions) resolve(taskName string) func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task) {
	if i, ok := ff.registeredFunctions[taskName]; ok {
		return i
	}
	return KinveyNotImplementedHandler()
}

func (ff *functions) process(task *Task, modules Modules) (*Task, *Task) {
	context := &Request{}
	var currentContext netType

	if task.HookType == "post" {
		currentContext = &task.Response
	} else {
		currentContext = &task.Request
	}

	context.Method = task.Request.Method
	context.Headers = currentContext.GetHeaders()
	context.Username = task.Request.Username
	context.UserID = task.Request.UserID

	if task.Request.ObjectName != "" {
		context.ObjectName = task.Request.ObjectName
	} else {
		context.ObjectName = task.Request.CollectionName
	}

	context.HookType = task.HookType

	if task.Request.EntityID != "" {
		context.EntityID = task.Request.EntityID
	}

	functionCompletionHandler := NewKinveyCompletionHandler(task)
	functionHandler := ff.resolve(task.TaskName)

	return functionHandler(context, functionCompletionHandler, modules)
}

func (ff *functions) clearAll() {
	ff.registeredFunctions = make(map[string]func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task))
}

// Register ...
func (ff *functions) Register(taskName string, functionToExecute func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)) {
	ff.registeredFunctions[taskName] = functionToExecute
}

func newFunctions() Functions {
	ff := &functions{
		registeredFunctions: make(map[string]func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task)),
	}
	return ff
}
