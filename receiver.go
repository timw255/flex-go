package flex

type receiver interface {
	Start(flex Flex, taskReceivedCallback func(task *Task) (*Task, *Task), options string) error
	Stop() error
}

func newReceiver(options *Options) receiver {
	if options.receiverType == receiverTypeHTTP {
		return &httpReceiver{}
	}
	return &tcpReceiver{}
}
