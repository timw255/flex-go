package flex

type authCompletionResponse struct {
	Token string `json:"token"`
}

// AuthCompletionHandler ...
type AuthCompletionHandler struct {
	task *Task
}

// NewAuthCompletionHandler ...
func NewAuthCompletionHandler(task *Task) AuthCompletionHandler {
	a := AuthCompletionHandler{
		task: task,
	}

	return a
}

// SetToken ...
func (a *AuthCompletionHandler) SetToken(token string) *AuthCompletionHandler {
	acr := authCompletionResponse{
		Token: token,
	}
	bytes, _ := json.Marshal(acr)
	a.task.Response.Body = bytes
	return a
}

// AddAttribute ...
func (a *AuthCompletionHandler) AddAttribute(key string, value interface{}) {

}

// RemoveAttribute ...
func (a *AuthCompletionHandler) RemoveAttribute(key string) {

}

// OK ...
func (a *AuthCompletionHandler) OK() *AuthCompletionHandler {
	return a
}

// ServerError ...
func (a *AuthCompletionHandler) ServerError() {

}

// AccessDenied ...
func (a *AuthCompletionHandler) AccessDenied() {

}

// TemporarilyUnavailable ...
func (a *AuthCompletionHandler) TemporarilyUnavailable() {

}

// NotImplemented ...
func (a *AuthCompletionHandler) NotImplemented() *AuthCompletionHandler {
	return a
}

// Next ...
func (a *AuthCompletionHandler) Next() *AuthCompletionHandler {
	return a
}

// Done ...
func (a *AuthCompletionHandler) Done() (*Task, *Task) {
	return nil, a.task
}
