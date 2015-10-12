package goprof

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/pprof"
	"runtime/trace"
)

var (
	// if we are writing profiles right now, we keep here the name of the working folder where we write profiles to
	// if we are not writing profiles at the moment, it's just an empty string
	ourProfilesDirectory = ""
)

// these types used for mocking functions which start/stop profiling
type startFxn func(profilesDir string) error
type stopFxn func()

// ProfilingInProgress returns true if we write profiles at the moment
func ProfilingInProgress() bool {
	return ourProfilesDirectory != ""
}

// StartProfiling starts writing profiles and returns path to the directory where they will be placed
// if anything goes wrong corresponding error is returned and no profiling is started
// If writing profiles is in progress it returns an error
func StartProfiling() (profilesDirectory string, err error) {
	dir, err := startProfiling(startWritingTrace, trace.Stop, startCPUProfiling, pprof.StopCPUProfile)
	if err != nil {
		logf("Failed to start writing profiles: %v", err)
	} else {
		logf("Start writing profiles to '%s'", dir)
	}
	return dir, err
}

// StopProfiling stops writing all profiles. Before stopping it tries to write a heap dump
// to the same folder where the other profiles are kept. It returns path to the folder which contains profiling files
// If profiling is not in progress, this method does nothing and returns empty string
func StopProfiling() (profilesDirectory string) {
	dir := stopProfiling(writeHeapProfile, trace.Stop, pprof.StopCPUProfile)
	logf("Stop writing profiles to '%s'", dir)
	return dir
}

// ToggleProfiling changes state of writing profiles to the opposite
func ToggleProfiling() (profilesDirectory string, err error) {
	if ProfilingInProgress() {
		return StopProfiling(), nil
	}
	return StartProfiling()
}

func startProfiling(startTrace startFxn, stopTrace stopFxn, startCPU startFxn, stopCPU stopFxn) (profilesDirectory string, err error) {
	if ProfilingInProgress() {
		return "", fmt.Errorf("Cannot start profiling, since it's already started")
	}
	profiles, err := ioutil.TempDir("", "profiles")
	if err != nil {
		return "", err
	}
	defer func() {
		// if something is wrong, do cleanup
		if err != nil {
			stopCPU()
			stopTrace()
			ourProfilesDirectory = ""
			// TODO: log error if any?
			os.RemoveAll(profiles)
		}
	}()
	if err := startTrace(profiles); err != nil {
		return "", err
	}
	if err := startCPU(profiles); err != nil {
		return "", err
	}
	ourProfilesDirectory = profiles
	return profiles, nil
}

func stopProfiling(writeHeap startFxn, stopTrace, stopCPU stopFxn) (profilesDirectory string) {
	if !ProfilingInProgress() {
		return ""
	}
	defer func() {
		// stop everything when we are finished with writing heap profile
		stopCPU()
		stopTrace()
		ourProfilesDirectory = ""
	}()
	if err := writeHeap(ourProfilesDirectory); err != nil {
		logf("Failed to write heap profile: %v", err)
	}
	return ourProfilesDirectory
}

func startWritingTrace(profilesDir string) error {
	traceFile, err := os.Create(filepath.Join(profilesDir, "trace"))
	if err != nil {
		return err
	}
	return trace.Start(traceFile)
}

func writeHeapProfile(profilesDir string) error {
	heapProfileFile, err := os.Create(filepath.Join(profilesDir, "heap-profile"))
	if err != nil {
		return err
	}
	defer heapProfileFile.Close()
	return pprof.WriteHeapProfile(heapProfileFile)
}

func startCPUProfiling(profilesDir string) error {
	cpuProfileFile, err := os.Create(filepath.Join(profilesDir, "cpu-profile"))
	if err != nil {
		return err
	}
	return pprof.StartCPUProfile(cpuProfileFile)
}
