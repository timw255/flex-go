// +build js,wasm

package flex

import (
	"bytes"
	"errors"
	"fmt"
	"syscall/js"
)

type receiver interface {
	Start(flex Flex, taskReceivedCallback func(task *Task) (*Task, *Task), options string) error
	Stop() error
}

func newReceiver(options *Options) receiver {
	return &taskReceiver{}
}

var (
	quit = make(chan bool)
)

type taskReceiver struct {
}

func (rec *taskReceiver) composeErrorReply() {
	return
}

func (rec *taskReceiver) parseTask(data []byte) (*Task, error) {
	parsedTask := &Task{}

	err := json.Unmarshal(data, &parsedTask)
	if err != nil {
		return nil, errors.New("Error parsing task")
	}

	return parsedTask, nil
}

func (rec *taskReceiver) Start(flex Flex, taskReceivedCallback func(task *Task) (*Task, *Task), options string) error {
	healthCheckBytes := []byte(`{"healthCheck":1}`)

	processTaskFunction := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		data := []byte(args[0].String())
		if bytes.Equal(data, healthCheckBytes) {
			return string([]byte(`{"status":"ready"}`))
		}

		parsedTask, err := rec.parseTask(data)
		if err != nil {
			fmt.Println("error parsing task")
		}

		taskErr, result := taskReceivedCallback(parsedTask)
		if taskErr != nil {
			fmt.Println("error taskreceived callback")
		}

		json, err := json.Marshal(result)
		if err != nil {
			fmt.Println("error marshall result")
		}
		return string(json)
	})

	js.Global().Set("processTask", processTaskFunction)

	sd := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Stopping Go!")
		rec.Stop()
		return nil
	})

	defer sd.Release()

	addEventListener := js.Global().Get("addEventListener")
	if addEventListener != js.Undefined() {
		addEventListener.Invoke("beforeunload", sd)
	}

	process := js.Global().Get("process")
	if process != js.Undefined() {
		process.Call("on", "SIGTERM", sd)
		process.Call("on", "SIGINT", sd)

		stdot := process.Get("stdout")
		if stdot != js.Undefined() {
			stdot.Call("write", `{"message":"Ready to Go!","level":"info"}`)
			stdot.Call("write", "\n")
		}
	}

	<-quit

	return nil
}

func (rec *taskReceiver) Stop() error {
	close(quit)
	return nil
}
