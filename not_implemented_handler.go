package flex

// AuthNotImplementedHandler ...
func AuthNotImplementedHandler() func(context *Request, complete AuthCompletionHandler, modules Modules) (*Task, *Task) {
	return func(context *Request, complete AuthCompletionHandler, modules Modules) (*Task, *Task) {
		return complete.NotImplemented().Done()
	}
}

// KinveyNotImplementedHandler ...
func KinveyNotImplementedHandler() func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task) {
	return func(context *Request, complete KinveyCompletionHandler, modules Modules) (*Task, *Task) {
		body := make(map[string]interface{})
		body["message"] = "These methods are not implemented"

		resBytes, _ := json.Marshal(body)

		return complete.SetBody(resBytes).NotImplemented().Done()
	}
}
