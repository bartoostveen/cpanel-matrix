package util

import (
	"os"
	"os/signal"
)

func OnSignal(handler func(), sig ...os.Signal) chan interface{} {
	c := make(chan os.Signal)
	after := make(chan interface{})
	signal.Notify(c, sig...)
	go func() {
		received := <-c
		handler()
		after <- received
	}()
	return after
}
