package flex

// KinveyCompletionHandler ...
type KinveyCompletionHandler struct {
	Task *Task
}

// NewKinveyCompletionHandler ...
func NewKinveyCompletionHandler(task *Task) KinveyCompletionHandler {
	a := KinveyCompletionHandler{
		Task: task,
	}

	return a
}

// SetBody ...
func (a *KinveyCompletionHandler) SetBody(body []byte) *KinveyCompletionHandler {
	a.Task.Response.Body = body
	return a
}

// NotImplemented ...
func (a *KinveyCompletionHandler) NotImplemented() *KinveyCompletionHandler {
	return a
}

// Done ...
func (a *KinveyCompletionHandler) Done() (*Task, *Task) {
	return nil, a.Task
}
