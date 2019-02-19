package flex

// Auth ...
type Auth interface {
	clearAll()
	getHandlers() []string
	process(task *Task, modules Modules) (*Task, *Task)
	resolve(taskName string) func(context *Request, complete AuthCompletionHandler, modules Modules) (*Task, *Task)
	Register(taskName string, functionToExecute func(context *Request, complete AuthCompletionHandler, modules Modules) (*Task, *Task))
}

type auth struct {
	authFunctions map[string]func(context *Request, complete AuthCompletionHandler, modules Modules) (*Task, *Task)
}

func (fa *auth) clearAll() {
	fa.authFunctions = make(map[string]func(req *Request, complete AuthCompletionHandler, modules Modules) (*Task, *Task))
}

func (fa *auth) getHandlers() []string {
	keys := make([]string, 0)
	for key := range fa.authFunctions {
		keys = append(keys, key)
	}
	return keys
}

func (fa *auth) process(task *Task, modules Modules) (*Task, *Task) {
	authCompletionHandler := NewAuthCompletionHandler(task)
	requestCompletionHandler := authCompletionHandler
	authHandler := fa.resolve(task.TaskName)
	return authHandler(&task.Request, requestCompletionHandler, modules)
}

func (fa *auth) resolve(taskName string) func(context *Request, complete AuthCompletionHandler, modules Modules) (*Task, *Task) {
	if i, ok := fa.authFunctions[taskName]; ok {
		return i
	}
	return AuthNotImplementedHandler()
}

// Register ...
func (fa *auth) Register(taskName string, functionToExecute func(context *Request, complete AuthCompletionHandler, modules Modules) (*Task, *Task)) {
	fa.authFunctions[taskName] = functionToExecute
}

func newAuth() Auth {
	ff := &auth{
		authFunctions: make(map[string]func(req *Request, complete AuthCompletionHandler, modules Modules) (*Task, *Task)),
	}
	return ff
}
