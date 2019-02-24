// +build !js,!wasm

package flex

import (
	gocontext "context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type healthCheckResponse struct {
	Healthy bool `json:"healthy"`
}

type functionsResponse struct {
	Request  *Request  `json:"request"`
	Response *Response `json:"response"`
}

type locals struct {
	Body            []byte                 `json:"body"`
	HookType        string                 `json:"hookType"`
	Method          string                 `json:"method"`
	Query           string                 `json:"query"`
	ObjectName      string                 `json:"objectName"`
	EntityID        string                 `json:"entityId"`
	TempObjectStore map[string]interface{} `json:"tempObjectStore"`
}

type httpReceiver struct {
	server *http.Server
}

func (rec *httpReceiver) healthCheck() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		healthCheckResponse := healthCheckResponse{
			Healthy: true,
		}

		json, err := json.Marshal(healthCheckResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	})
}

func (rec *httpReceiver) mapPostToElements(ctx *context) error {
	err := json.Unmarshal(ctx.Request.Body, &ctx.Locals)
	if err != nil {
		return err
	}

	return nil
}

func (rec *httpReceiver) generateBaseTask(ctx *context) error {
	var method string

	if ctx.Locals.Method != "" {
		method = ctx.Locals.Method
	} else {
		method = *ctx.Request.Method
	}

	var appMetadata kinveyAppMetadata
	if kam := ctx.Request.Header.Get("X-Kinvey-App-Metadata"); kam != "" {
		k := kinveyAppMetadata{}
		if err := json.Unmarshal([]byte(kam), &k); err == nil {
			appMetadata = k
		}
	}

	requestHeaders := make(map[string]string)
	if orh := ctx.Request.Header.Get("X-Kinvey-Original-Request-Headers"); orh != "" {
		k := kinveyOriginalRequestHeaders{}
		if err := json.Unmarshal([]byte(orh), &k); err == nil {
			requestHeaders["x-kinvey-api-version"] = k.KinveyAPIVersion
			requestHeaders["authorization"] = k.Authorization
			requestHeaders["x-kinvey-client-app-version"] = k.KinveyClientAppVersion
		}
	}

	ctx.Task.AppID = ctx.Request.Header.Get("X-Kinvey-Environment-Id")
	ctx.Task.AppMetadata = appMetadata
	ctx.Task.AuthKey = ctx.Request.Header.Get("X-Auth-Key")
	ctx.Task.RequestID = ctx.Request.Header.Get("X-Kinvey-Request-Id")
	ctx.Task.Method = strings.ToUpper(method)
	ctx.Task.Request = Request{
		Method:   method,
		Headers:  requestHeaders,
		Username: ctx.Request.Header.Get("X-Kinvey-Username"),
		UserID:   ctx.Request.Header.Get("X-Kinvey-User-Id"),
	}
	ctx.Task.Response = ctx.Response

	if krs, err := strconv.Atoi(ctx.Request.Header.Get("X-Kinvey-Response-Status")); err == nil {
		ctx.Task.Response.Status = krs
	} else if ctx.Response.Status != 0 {
		ctx.Task.Response.Status = ctx.Response.Status
	} else {
		ctx.Task.Response.Status = 0
	}

	if krh := ctx.Request.Header.Get("X-Kinvey-Response-Headers"); krh != "" {
		responseHeaders := make(map[string]string)
		if err := json.Unmarshal([]byte(krh), &responseHeaders); err == nil {
			ctx.Task.Response.Headers = responseHeaders
		}
	}

	if krb := ctx.Request.Header.Get("X-Kinvey-Response-Body"); krb != "" {
		ctx.Task.Response.Body = []byte(krb)
	}

	return nil
}

func (rec *httpReceiver) addFunctionsTaskAttributes(ctx *context) error {
	ctx.Task.TaskType = "functions"

	if ctx.Locals.ObjectName != "" {
		ctx.Task.Request.ObjectName = ctx.Locals.ObjectName
	}

	ctx.Task.TaskName = path.Base(ctx.Request.URL.Path)

	if ctx.Locals.HookType != "" {
		ctx.Task.HookType = ctx.Locals.HookType
	} else {
		ctx.Task.HookType = "customEndpoint"
	}

	ctx.Task.Request.TempObjectStore = ctx.Locals.TempObjectStore

	return nil
}

func (rec *httpReceiver) addDataTaskAttributes(ctx *context) error {
	ctx.Task.TaskType = "data"
	ctx.Task.Request.ServiceObjectName = strings.Split(ctx.Request.URL.Path, "/")[1]

	return nil
}

func (rec *httpReceiver) appendID(ctx *context) error {
	if id := path.Base(ctx.Request.URL.Path); id != "" {
		ctx.Task.Request.EntityID = id
	} else if ctx.Locals.EntityID != "" {
		ctx.Task.Request.EntityID = ctx.Locals.EntityID
	}

	return nil
}

func (rec *httpReceiver) appendQuery(ctx *context) error {
	if ctx.Request.URL.RawQuery != "" {
		ctx.Task.Request.Query = ctx.Request.URL.Query()
	} else if ctx.Locals.Query != "" {
		qry, err := url.ParseQuery(ctx.Locals.Query)
		if err != nil {
			return err
		}
		ctx.Task.Request.Query = qry
	}

	return nil
}

func (rec *httpReceiver) appendCount(ctx *context) error {
	ctx.Task.Endpoint = "_count"

	return nil
}

func (rec *httpReceiver) appendBody(ctx *context) error {
	if ctx.Locals.Body != nil {
		ctx.Task.Request.Body = ctx.Locals.Body
	} else {
		ctx.Task.Request.Body = ctx.Request.Body
	}
	return nil
}

func (rec *httpReceiver) getDataBody(task *Task) string {
	return task.Response.GetBody()
}

func (rec *httpReceiver) getFunctionsBody(task *Task) string {
	fr := functionsResponse{
		Request:  &task.Request,
		Response: &task.Response,
	}

	json.Unmarshal(fr.Request.Body, &fr.Request.JSONBody)

	json.Unmarshal(fr.Response.Body, &fr.Response.JSONBody)

	json, _ := json.Marshal(fr)

	return string(json)
}

func (rec *httpReceiver) getAuthBody(task *Task) string {
	return task.Response.GetBody()
}

func (rec *httpReceiver) getDiscoveryBody(task *Task) string {
	resBytes, _ := json.Marshal(task.DiscoveryObjects)

	return string(resBytes)
}

func (rec *httpReceiver) sendTask(task *Task, taskReceivedCallback func(task *Task) (*Task, *Task)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err, result := taskReceivedCallback(task)

		if err != nil {
			c := true

			w.Header().Set("X-Kinvey-Request-Continue", strconv.FormatBool(c))
		} else if result != nil {
			var body string

			switch task.TaskType {
			case "data":
				body = rec.getDataBody(result)
				break
			case "functions":
				body = rec.getFunctionsBody(result)
				break
			case "serviceDiscovery":
				body = rec.getDiscoveryBody(result)
				break
			case "auth":
				body = rec.getAuthBody(result)
				break
			default:
				body = ""
			}

			w.Header().Set("Connection", "close")
			w.Header().Set("Content-Type", "application/json")

			w.Write([]byte(body))
		}

	})
}

func (rec *httpReceiver) buildDiscoverTask(ctx *context) error {
	ctx.Task.TaskType = "serviceDiscovery"
	ctx.Task.Request = Request{}
	ctx.Task.Response = Response{}

	return nil
}

func (rec *httpReceiver) addAuthTaskAttributes(ctx *context) error {
	ctx.Task.TaskType = "auth"
	ctx.Task.TaskName = path.Base(ctx.Request.URL.Path)

	return nil
}

func (rec *httpReceiver) dataGroupHandler(taskReceivedCallback func(task *Task) (*Task, *Task)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		segments := strings.Split(r.URL.Path, "/")

		if len(segments) == 3 {
			if segments[2] == "_count" {
				newChain(rec.generateBaseTask, rec.addDataTaskAttributes, rec.appendQuery, rec.appendCount).then(rec.sendTask, taskReceivedCallback).ServeHTTP(w, r)
				return
			}

			newChain(rec.generateBaseTask, rec.addDataTaskAttributes, rec.appendQuery, rec.appendID).then(rec.sendTask, taskReceivedCallback).ServeHTTP(w, r)
			return
		}

		newChain(rec.generateBaseTask, rec.addDataTaskAttributes, rec.appendQuery).then(rec.sendTask, taskReceivedCallback).ServeHTTP(w, r)
	})
}

// Start ...
func (rec *httpReceiver) Start(flex Flex, taskReceivedCallback func(task *Task) (*Task, *Task), options string) error {
	router := gin.New()

	router.POST("/healthcheck", gin.WrapH(rec.healthCheck()))

	// FlexFunctions
	ff := router.Group("/_flexFunctions/")
	{
		for _, h := range flex.Functions.getHandlers() {
			ff.POST(h, gin.WrapH(newChain(rec.mapPostToElements, rec.generateBaseTask, rec.addFunctionsTaskAttributes, rec.appendQuery, rec.appendID, rec.appendBody).then(rec.sendTask, taskReceivedCallback)))
		}
	}

	// FlexAuth
	fa := router.Group("/_auth/")
	{
		for _, h := range flex.Auth.getHandlers() {
			fa.POST(h, gin.WrapH(newChain(rec.mapPostToElements, rec.generateBaseTask, rec.addFunctionsTaskAttributes, rec.appendQuery, rec.appendID, rec.appendBody).then(rec.sendTask, taskReceivedCallback)))
		}
	}

	// Command
	router.POST("/_command/discover", gin.WrapH(newChain(rec.buildDiscoverTask).then(rec.sendTask, taskReceivedCallback)))

	// FlexData
	for _, so := range flex.Data.getServiceObjects() {
		g := router.Group("/" + so)
		{
			g.POST("", gin.WrapH(newChain(rec.generateBaseTask, rec.addDataTaskAttributes, rec.appendBody).then(rec.sendTask, taskReceivedCallback)))
			g.DELETE("", gin.WrapH(newChain(rec.generateBaseTask, rec.addDataTaskAttributes, rec.appendQuery).then(rec.sendTask, taskReceivedCallback)))
			g.GET("", gin.WrapH(rec.dataGroupHandler(taskReceivedCallback)))

			j := g.Group("/:param")
			{
				j.GET("", gin.WrapH(rec.dataGroupHandler(taskReceivedCallback)))
				j.PUT("", gin.WrapH(newChain(rec.generateBaseTask, rec.addDataTaskAttributes, rec.appendID, rec.appendBody).then(rec.sendTask, taskReceivedCallback)))
				j.DELETE("", gin.WrapH(newChain(rec.generateBaseTask, rec.addDataTaskAttributes, rec.appendID, rec.appendQuery).then(rec.sendTask, taskReceivedCallback)))
			}
		}
	}

	/*d := router.Group("/debug/pprof/")
	{
		d.GET("", gin.WrapH(http.HandlerFunc(pprof.Index)))
		d.GET("cmdline", gin.WrapH(http.HandlerFunc(pprof.Cmdline)))
		d.GET("profile", gin.WrapH(http.HandlerFunc(pprof.Profile)))
		d.GET("symbol", gin.WrapH(http.HandlerFunc(pprof.Symbol)))
		d.GET("trace", gin.WrapH(http.HandlerFunc(pprof.Trace)))
	}*/

	address := ":10001"
	rec.server = &http.Server{
		Addr:    address,
		Handler: router,
		//ReadTimeout:  5 * time.Second,
		//WriteTimeout: 10 * time.Second,
	}

	rec.server.ListenAndServe()

	fmt.Println(fmt.Sprintf("Service listening on %s", address))

	return nil
}

// Stop ...
func (rec *httpReceiver) Stop() error {
	ctx := gocontext.Background()
	err := rec.server.Shutdown(ctx)

	if err != nil {
		return err
	}

	return nil
}

type context struct {
	Response Response
	Request  request
	Locals   locals
	Task     *Task
}

type request struct {
	URL    *url.URL
	Method *string
	Header *http.Header
	Body   []byte
}

type constructor func(r *context) error

type chain struct {
	constructors []constructor
}

func newChain(constructors ...constructor) chain {
	return chain{append(([]constructor)(nil), constructors...)}
}

func (c chain) then(handler func(task *Task, taskReceivedCallback func(task *Task) (*Task, *Task)) http.Handler, taskReceivedCallback func(task *Task) (*Task, *Task)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := context{
			Request: request{},
			Locals:  locals{},
			Task:    &Task{},
		}

		context.Request.Method = &r.Method
		context.Request.Header = &r.Header
		context.Request.URL = r.URL

		body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 26214400))

		context.Request.Body = body

		j := make(map[string]interface{})
		json.Unmarshal(body, &j)

		if rb := j["response"]; rb != nil {
			json.Unmarshal([]byte(fmt.Sprint(rb)), &context.Response)
		}

		for i := range c.constructors {
			c.constructors[i](&context)
		}

		handler(context.Task, taskReceivedCallback).ServeHTTP(w, r)
	})
}
