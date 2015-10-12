package goprof

import (
	"fmt"
	"os"
	"testing"
)

type mockStopper struct {
	called bool
}

func (m *mockStopper) fxn() stopFxn {
	return func() {
		m.called = true
	}
}

type mockStarter struct {
	profileDir string
}

func (m *mockStarter) fxn(result error) startFxn {
	return func(dir string) error {
		m.profileDir = dir
		return result
	}
}

func TestMain(m *testing.M) {
	noLogging := func(format string, args ...interface{}) {}
	SetLogFunction(noLogging)
	os.Exit(m.Run())
}

func TestStopWhenNotRunning(t *testing.T) {
	if path := StopProfiling(); path != "" {
		t.Fatalf("Expected empty string when stopping not running profiling. Got '%s'", path)
	}
}

func startMockProfiling() (string, error) {
	startTrace, startCPU := &mockStarter{}, &mockStarter{}
	stopTrace, stopCPU := &mockStopper{}, &mockStopper{}
	return startProfiling(startTrace.fxn(nil), stopTrace.fxn(), startCPU.fxn(nil), stopCPU.fxn())
}

func TestStartTraceFailed(t *testing.T) {
	startTrace, startCPU := &mockStarter{}, &mockStarter{}
	stopTrace, stopCPU := &mockStopper{}, &mockStopper{}
	dir, err := startProfiling(startTrace.fxn(fmt.Errorf("test")), stopTrace.fxn(), startCPU.fxn(nil), stopCPU.fxn())
	if dir != "" || err == nil {
		t.Fatalf("Start profiling should return error and no dir. I got '%s' and %v", dir, err)
	}
	if ProfilingInProgress() {
		t.Fatalf("Profiling is running")
	}
	if !stopTrace.called || !stopCPU.called {
		t.Fatalf("Profiles was not stopped")
	}
	if _, statErr := os.Stat(startTrace.profileDir); statErr == nil {
		t.Errorf("Temporary dir for profiling data seem to exist")
	}
}

func TestStartCPUFailed(t *testing.T) {
	startTrace, startCPU := &mockStarter{}, &mockStarter{}
	stopTrace, stopCPU := &mockStopper{}, &mockStopper{}
	dir, err := startProfiling(startTrace.fxn(nil), stopTrace.fxn(), startCPU.fxn(fmt.Errorf("test")), stopCPU.fxn())
	if dir != "" || err == nil {
		t.Fatalf("Start profiling should return error and no dir. I got '%s' and %v", dir, err)
	}
	if ProfilingInProgress() {
		t.Fatalf("Profiling is running")
	}
	if !stopTrace.called || !stopCPU.called {
		t.Fatalf("Profiles was not stopped")
	}
	if _, statErr := os.Stat(startTrace.profileDir); statErr == nil {
		t.Errorf("Temporary dir '%s' for profiling data seem to exist", startTrace.profileDir)
	}
}

func TestStopFailedToDumpHeap(t *testing.T) {
	startDir, err := startMockProfiling()
	if startDir == "" || err != nil {
		t.Fatalf("Profiling should be started successfully. I got '%s' and %v", startDir, err)
	}
	writeHeap := &mockStarter{}
	stopTrace, stopCPU := &mockStopper{}, &mockStopper{}
	stopDir := stopProfiling(writeHeap.fxn(fmt.Errorf("test")), stopTrace.fxn(), stopCPU.fxn())
	if stopDir != startDir {
		t.Fatalf("Different dirs for start and stop: '%s' and '%s'", startDir, stopDir)
	}
	if !stopTrace.called || !stopCPU.called {
		t.Fatalf("Profiles was not stopped")
	}
	if ProfilingInProgress() {
		t.Fatalf("Profiling is running")
	}
	os.RemoveAll(startDir)
}

func TestStopCallsStop(t *testing.T) {
	startDir, err := startMockProfiling()
	if startDir == "" || err != nil {
		t.Fatalf("Profiling should be started successfully. I got '%s' and %v", startDir, err)
	}
	writeHeap := &mockStarter{}
	stopTrace, stopCPU := &mockStopper{}, &mockStopper{}
	stopDir := stopProfiling(writeHeap.fxn(nil), stopTrace.fxn(), stopCPU.fxn())
	if stopDir != startDir {
		t.Fatalf("Different dirs for start and stop: '%s' and '%s'", startDir, stopDir)
	}
	if !stopTrace.called || !stopCPU.called {
		t.Fatalf("Profiles was not stopped")
	}
	if ProfilingInProgress() {
		t.Fatalf("Profiling is running")
	}
	os.RemoveAll(startDir)
}
