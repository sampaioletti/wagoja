package jswasm

import (
	"github.com/dop251/goja_nodejs/eventloop"
)

//NewJS returns a runtime that includes the necessary polyfils
func NewJS() *eventloop.EventLoop {
	vm := eventloop.NewEventLoop()
	vm.Start()
	return vm
}
