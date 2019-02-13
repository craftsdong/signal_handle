package signal_handle

import (
	"fmt"
	"syscall"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var msg string

func mockCB() {
	msg = "end"
}
func TestSignalHandle(t *testing.T) {
	Convey("SignalHandle", t, func() {
		Convey("new", func() {
			So(signalHandle, ShouldNotBeNil)
		})
		Convey("callback", func() {
			RegisterCallback(mockCB)
			go signalHandle.Listen()
			//mock kill signal
			signalHandle.sigChan <- syscall.SIGINT

			signalHandle.Wait()
			So(msg, ShouldEqual, "end")
		})
	})
}

func ExampleSignalHandle() {
	s := NewSignalHandle()
	//register call back
	s.RegisterCallback(func() {
		//I am a callback function
		fmt.Println("callback")
	})

	//listen kill signal
	go s.Listen()
	s.sigChan <- syscall.SIGINT

	//waiting for deal with all callbacks or after timeout
	s.Wait()

	// Output:
	// callback
}
