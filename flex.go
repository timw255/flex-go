package flex

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/json-iterator/go"
	"github.com/timw255/flex-go/util"
)

const (
	flexGoVersion    = "0.0.0"
	receiverTypeHTTP = "http"
	receiverTypeTCP  = "tcp"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary

	terminated    bool
	rec           receiver
	flexTaskTypes = []string{
		"data",
		"functions",
		"auth",
		"serviceDiscovery",
	}
)

func init() {
	terminated = false
}

// Options ...
type Options struct {
	host         string
	port         int
	sharedSecret string
	receiverType string
}

// NewOptions ...
func NewOptions(host string, port int, sharedSecret string) *Options {
	o := &Options{
		host,
		port,
		sharedSecret,
		"http",
	}

	return o
}

// Flex ...
type Flex struct {
	Data         Data
	Functions    Functions
	Auth         Auth
	Logger       Logger
	version      string
	sharedSecret string
}

// NewService ...
func NewService(options *Options, initializer func(err error, flex Flex)) {
	sdkReceiver, ok := os.LookupEnv("SDK_RECEIVER")

	if ok && sdkReceiver == receiverTypeTCP {
		options.receiverType = receiverTypeTCP
	} else {
		options.receiverType = receiverTypeHTTP
	}

	//options.receiverType = receiverTypeTCP

	rec = newReceiver(options)

	d := newData()
	f := newFunctions()
	a := newAuth()
	l := newLogger()

	s := Flex{
		Data:         d,
		Functions:    f,
		Auth:         a,
		Logger:       l,
		version:      flexGoVersion,
		sharedSecret: options.sharedSecret,
	}

	taskReceivedCallback := func(task *Task) (*Task, *Task) {
		task.SDKVersion = flexGoVersion

		if options.sharedSecret != "" && task.TaskType != "serviceDiscovery" && task.TaskType != "logger" && task.TaskType != "moduleGenerator" && options.sharedSecret != task.AuthKey {
			return nil, task
		}

		if !util.Contains(flexTaskTypes, task.TaskType) {
			return nil, task
		}

		if task.TaskType == "serviceDiscovery" {
			so := dataLink{
				ServiceObjects: s.Data.getServiceObjects(),
			}
			fh := businessLogic{
				Handlers: s.Functions.getHandlers(),
			}
			ah := authDiscovery{
				Handlers: s.Auth.getHandlers(),
			}

			dco := discoveryObjects{
				DataLink:      so,
				BusinessLogic: fh,
				Auth:          ah,
			}

			task.DiscoveryObjects = dco

			return nil, task
		}

		modules := generateModules(task)

		switch task.TaskType {
		case "data":
			return s.Data.process(task, modules)
		case "functions":
			return s.Functions.process(task, modules)
		case "auth":
			return s.Auth.process(task, modules)
		}

		return nil, nil
	}

	var gracefulShutdown = make(chan os.Signal)
	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)
	go func() {
		sig := <-gracefulShutdown
		if terminated {
			fmt.Println("Forced quit")
			terminate(errors.New("Forced quit"))
			return
		}

		terminated = true
		fmt.Println(fmt.Sprintf("Caught %+v, initiating graceful shutdown. Press ctrl-c or send SIGTERM/SIGINT to force-quit immediately.", sig))

		c := make(chan bool, 1)
		go func() {
			err := rec.Stop()
			if err != nil {
				terminate(err)
			}
			c <- true
		}()
		select {
		case res := <-c:
			_ = res
			return
		case <-time.After(50 * time.Second):
			terminate(nil)
		}

		return
	}()

	initializer(nil, s)

	err := rec.Start(s, taskReceivedCallback, "")
	_ = err
}

func terminate(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	os.Exit(0)
}
