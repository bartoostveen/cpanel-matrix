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

func KeysOf[V any](m map[string]V) []string {
	keys := make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}
