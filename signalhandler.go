package SignalHandler

import (
	"fmt"
	"log"
	"os"
	"os/signal"
)

// SignalHandler is a basic, abstract signal handler that is implementable in most go applications without much effort
type SignalHandler struct {
	signalChannel chan os.Signal
	callbacks     map[os.Signal][]func() error
	outputFunc    func(v ...interface{})
}

// ConstructSignalHandler will create a new SignalHandler
func ConstructSignalHandler() *SignalHandler {
	return &SignalHandler{
		signalChannel: make(chan os.Signal, 1),
		callbacks:     make(map[os.Signal][]func() error, 10),
		outputFunc:    log.Println,
	}
}

// SetOutput will set the output function used to print errors should they occurr
func (handler *SignalHandler) SetOutput(outputFunction func(v ...interface{})) {
	handler.outputFunc = outputFunction
}

// Start will listen, async, in a seperate go-routine.
func (handler *SignalHandler) Start() {
	go handler.Listen()
}

// Listen will wait for signals to come in, go through registered signal functions and execute them
func (handler *SignalHandler) Listen() {
	for {
		receivedSignal := <-handler.signalChannel

		for _, signalFunction := range handler.callbacks[receivedSignal] {
			err := signalFunction()

			if err != nil {
				handler.outputFunc(fmt.Sprintf("SignalHandler: error during signal %d func: %v", receivedSignal, err))
			}
		}
	}
}

// RegisterSignalFunction allows you to specify a function to be called when a signal is received by the signal handler
func (handler *SignalHandler) RegisterSignalFunction(sig os.Signal, callback func() error) {
	if len(handler.callbacks[sig]) == 0 {
		signal.Ignore(sig)
		signal.Notify(handler.signalChannel, sig)
	}

	handler.callbacks[sig] = append(handler.callbacks[sig], callback)
}
