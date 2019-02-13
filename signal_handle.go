package signal_handle

import (
	"os"
	"os/signal"
	"time"
)

var (
	signalHandle *SignalHandle
)

func init() {
	signalHandle = NewSignalHandle()
}

// SignalHandle defines
type SignalHandle struct {
	sigChan    chan os.Signal
	closing    chan struct{}
	cb         []func()
	maxTimeOut time.Duration
}

//NewSignalHandle create a new SignalHandle object
func NewSignalHandle() *SignalHandle {
	return &SignalHandle{
		sigChan: make(chan os.Signal),
		closing: make(chan struct{}),
	}
}

//RegisterCallback add callback function to global signalHandle
func RegisterCallback(f func()) {
	signalHandle.RegisterCallback(f)
}

//RegisterCallback add callback function
func (s *SignalHandle) RegisterCallback(f func()) {
	s.cb = append(s.cb, f)
}

//Listen register and listen kill signal for global signalHandle
func Listen(signals ...os.Signal) {
	signalHandle.Listen(signals...)
}

//Listen register and listen kill signal
func (s *SignalHandle) Listen(signals ...os.Signal) {
	if len(signals) > 0 {
		signal.Notify(
			s.sigChan,
			signals...,
		)
	} else { //all signals
		signal.Notify(s.sigChan)
	}

	//waiting for kill signal
	<-s.sigChan

	//waiting for handle max time
	if s.maxTimeOut > 0 {
		go exitAfter(s.maxTimeOut)
	}

	//shutdown callback
	for _, f := range s.cb {
		f()
	}

	close(s.closing)
}

//Wait global signalHandle run out
func Wait() {
	signalHandle.Wait()
}

//Wait until run out all callbacks
func (s *SignalHandle) Wait() {
	//waiting for deal all callback functions
	<-s.closing
}

//SetMaxTimeOut set max running time for global signalHandle's callback handle
func SetMaxTimeOut(after time.Duration) {
	signalHandle.SetMaxTimeOut(after)
}

//SetMaxTimeOut set max running time for callback handle
func (s *SignalHandle) SetMaxTimeOut(after time.Duration) {
	s.maxTimeOut = after
}

//exitAfter must be exit after some times
func exitAfter(after time.Duration) {
	select {
	case <-time.After(after):
		os.Exit(1)
	}
}
