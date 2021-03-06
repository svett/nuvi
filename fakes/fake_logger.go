// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/svett/nuvi"
)

type FakeLogger struct {
	PrintlnStub        func(v ...interface{})
	printlnMutex       sync.RWMutex
	printlnArgsForCall []struct {
		v []interface{}
	}
	PrintfStub        func(format string, v ...interface{})
	printfMutex       sync.RWMutex
	printfArgsForCall []struct {
		format string
		v      []interface{}
	}
}

func (fake *FakeLogger) Println(v ...interface{}) {
	fake.printlnMutex.Lock()
	fake.printlnArgsForCall = append(fake.printlnArgsForCall, struct {
		v []interface{}
	}{v})
	fake.printlnMutex.Unlock()
	if fake.PrintlnStub != nil {
		fake.PrintlnStub(v...)
	}
}

func (fake *FakeLogger) PrintlnCallCount() int {
	fake.printlnMutex.RLock()
	defer fake.printlnMutex.RUnlock()
	return len(fake.printlnArgsForCall)
}

func (fake *FakeLogger) PrintlnArgsForCall(i int) []interface{} {
	fake.printlnMutex.RLock()
	defer fake.printlnMutex.RUnlock()
	return fake.printlnArgsForCall[i].v
}

func (fake *FakeLogger) Printf(format string, v ...interface{}) {
	fake.printfMutex.Lock()
	fake.printfArgsForCall = append(fake.printfArgsForCall, struct {
		format string
		v      []interface{}
	}{format, v})
	fake.printfMutex.Unlock()
	if fake.PrintfStub != nil {
		fake.PrintfStub(format, v...)
	}
}

func (fake *FakeLogger) PrintfCallCount() int {
	fake.printfMutex.RLock()
	defer fake.printfMutex.RUnlock()
	return len(fake.printfArgsForCall)
}

func (fake *FakeLogger) PrintfArgsForCall(i int) (string, []interface{}) {
	fake.printfMutex.RLock()
	defer fake.printfMutex.RUnlock()
	return fake.printfArgsForCall[i].format, fake.printfArgsForCall[i].v
}

var _ nuvi.Logger = new(FakeLogger)
