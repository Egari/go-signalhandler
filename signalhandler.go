package SignalHandler

import (
	"fmt"
	"log"
	"os"
	"os/signal"
)

type SignalHandler struct {
	signalChannel chan os.Signal
	callbacks     map[os.Signal][]func() error
	outputFunc    func(string)
}

func ConstructSignalHandler() *SignalHandler {
	return &SignalHandler{
		signalChannel: make(chan os.Signal, 1),
		callbacks:     make(map[os.Signal][]func() error, 10),
		outputFunc     log.Errorln
	}
}

func (handler *SignalHandler) setOutput(outputFunction func(string)) {
	handler.outputFunc = outputFunction
}

func (handler *SignalHandler) Start() {
	go handler.Listen()
}

func (handler *SignalHandler) Listen() {
	for {
		receivedSignal := <-handler.signalChannel

		for _, signalFunction := range handler.callbacks[receivedSignal] {
			err := signalFunction()

			if err != nil {
				outputFunc(fmt.Sprintf("SignalHandler: error during signal %d func: %v", receivedSignal, err))
			}
		}
	}
}

func (handler *SignalHandler) RegisterSignalFunction(sig os.Signal, callback func() error) {
	if len(handler.callbacks[sig]) == 0 {
		signal.Ignore(sig)
		signal.Notify(handler.signalChannel, sig)
	}

	handler.callbacks[sig] = append(handler.callbacks[sig], callback)
}
