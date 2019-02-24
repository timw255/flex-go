// +build !js,!wasm

package flex

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type tcpReceiver struct {
	server net.Listener

	handlers sync.WaitGroup
}

func (rec *tcpReceiver) composeErrorReply() {
	return
}

func (rec *tcpReceiver) parseTask(data []byte) (*Task, error) {
	parsedTask := &Task{}

	err := json.Unmarshal(data, &parsedTask)
	if err != nil {
		return nil, errors.New("Error parsing task")
	}

	return parsedTask, nil
}

func (rec *tcpReceiver) Start(flex Flex, taskReceivedCallback func(task *Task) (*Task, *Task), options string) error {
	healthCheckBytes := []byte(`{"healthCheck":1}`)

	processTask := func(c net.Conn) {
		defer c.Close()

		buf := bufio.NewReader(c)
		for {
			data, _, err := buf.ReadLine()
			if err == io.EOF {
				break
			}

			if bytes.Equal(data, healthCheckBytes) {
				c.Write([]byte(`{"status":"ready"}`))
				c.Write([]byte("\n"))
				continue
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

			c.Write([]byte(string(json)))
			c.Write([]byte("\n"))
		}
	}

	server, err := net.Listen("tcp4", ":7000")
	if err != nil {
		return err
	}
	rec.server = server

	rec.serve(processTask)

	return nil
}

func (rec *tcpReceiver) serve(processTaskFunction func(c net.Conn)) {
	for {
		conn, err := rec.server.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			fmt.Println("Failed to accept connection:", err.Error())
		}
		rec.handlers.Add(1)
		go func() {
			processTaskFunction(conn)
			rec.handlers.Done()
		}()
	}
}

func (rec *tcpReceiver) Stop() error {
	err := rec.server.Close()
	if err != nil {
		return err
	}

	rec.handlers.Wait()

	return nil
}
