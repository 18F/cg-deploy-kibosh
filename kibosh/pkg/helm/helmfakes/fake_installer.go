// Code generated by counterfeiter. DO NOT EDIT.
package helmfakes

import (
	"sync"
	"time"

	"github.com/cf-platform-eng/kibosh/pkg/helm"
)

type FakeInstaller struct {
	InstallStub        func() error
	installMutex       sync.RWMutex
	installArgsForCall []struct {
	}
	installReturns struct {
		result1 error
	}
	installReturnsOnCall map[int]struct {
		result1 error
	}
	SetMaxWaitStub        func(time.Duration)
	setMaxWaitMutex       sync.RWMutex
	setMaxWaitArgsForCall []struct {
		arg1 time.Duration
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeInstaller) Install() error {
	fake.installMutex.Lock()
	ret, specificReturn := fake.installReturnsOnCall[len(fake.installArgsForCall)]
	fake.installArgsForCall = append(fake.installArgsForCall, struct {
	}{})
	fake.recordInvocation("Install", []interface{}{})
	fake.installMutex.Unlock()
	if fake.InstallStub != nil {
		return fake.InstallStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.installReturns
	return fakeReturns.result1
}

func (fake *FakeInstaller) InstallCallCount() int {
	fake.installMutex.RLock()
	defer fake.installMutex.RUnlock()
	return len(fake.installArgsForCall)
}

func (fake *FakeInstaller) InstallCalls(stub func() error) {
	fake.installMutex.Lock()
	defer fake.installMutex.Unlock()
	fake.InstallStub = stub
}

func (fake *FakeInstaller) InstallReturns(result1 error) {
	fake.installMutex.Lock()
	defer fake.installMutex.Unlock()
	fake.InstallStub = nil
	fake.installReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeInstaller) InstallReturnsOnCall(i int, result1 error) {
	fake.installMutex.Lock()
	defer fake.installMutex.Unlock()
	fake.InstallStub = nil
	if fake.installReturnsOnCall == nil {
		fake.installReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.installReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeInstaller) SetMaxWait(arg1 time.Duration) {
	fake.setMaxWaitMutex.Lock()
	fake.setMaxWaitArgsForCall = append(fake.setMaxWaitArgsForCall, struct {
		arg1 time.Duration
	}{arg1})
	fake.recordInvocation("SetMaxWait", []interface{}{arg1})
	fake.setMaxWaitMutex.Unlock()
	if fake.SetMaxWaitStub != nil {
		fake.SetMaxWaitStub(arg1)
	}
}

func (fake *FakeInstaller) SetMaxWaitCallCount() int {
	fake.setMaxWaitMutex.RLock()
	defer fake.setMaxWaitMutex.RUnlock()
	return len(fake.setMaxWaitArgsForCall)
}

func (fake *FakeInstaller) SetMaxWaitCalls(stub func(time.Duration)) {
	fake.setMaxWaitMutex.Lock()
	defer fake.setMaxWaitMutex.Unlock()
	fake.SetMaxWaitStub = stub
}

func (fake *FakeInstaller) SetMaxWaitArgsForCall(i int) time.Duration {
	fake.setMaxWaitMutex.RLock()
	defer fake.setMaxWaitMutex.RUnlock()
	argsForCall := fake.setMaxWaitArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeInstaller) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.installMutex.RLock()
	defer fake.installMutex.RUnlock()
	fake.setMaxWaitMutex.RLock()
	defer fake.setMaxWaitMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeInstaller) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ helm.Installer = new(FakeInstaller)
