// Copyright 2024, Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lifecycle

import (
	"fmt"
	"os"
)

// LifecycleStage represents the different stages of Nitric execution
type LifecycleStage string

const (
	// LocalRun represents local development run (using nitric run/start)
	LocalRun LifecycleStage = "run"
	// Build represents local development requirements building/collection (using nitric up)
	Build LifecycleStage = "build"
	// Cloud represents when the code is running in a deployed environment
	Cloud LifecycleStage = "cloud"
)

const (
	// NITRIC_ENVIRONMENT is the environment variable key used to determine the current Nitric lifecycle
	NITRIC_ENVIRONMENT = "NITRIC_ENVIRONMENT"
)

// GetCurrentLifecycle returns the current lifecycle stage
func GetCurrentLifecycle() (LifecycleStage, error) {
	lifecycle := os.Getenv(NITRIC_ENVIRONMENT)
	if lifecycle == "" {
		return "", fmt.Errorf("unable to determine the current Nitric lifecycle, please ensure the %s environment variable is set", NITRIC_ENVIRONMENT)
	}

	stage := LifecycleStage(lifecycle)
	switch stage {
	case LocalRun, Build, Cloud:
		return stage, nil
	default:
		return "", fmt.Errorf("invalid lifecycle stage: %s", lifecycle)
	}
}

// IsInLifecycle checks if the current environment is one of the provided stages
func IsInLifecycle(stages ...LifecycleStage) bool {
	currentStage, err := GetCurrentLifecycle()
	if err != nil {
		return false
	}

	for _, stage := range stages {
		if currentStage == stage {
			return true
		}
	}
	return false
}

// WhenInLifecycles executes the provided callback if the current environment is one of the provided stages
func WhenInLifecycles[T any](stages []LifecycleStage, callback func() T) T {
	if IsInLifecycle(stages...) {
		return callback()
	}
	var zero T
	return zero
}

// WhenRunning executes the provided callback if the current environment is running (LocalRun or Cloud)
func WhenRunning[T any](callback func() T) T {
	return WhenInLifecycles([]LifecycleStage{LocalRun, Cloud}, callback)
}

// WhenCollecting executes the provided callback if the current environment is collecting requirements (Build)
func WhenCollecting[T any](callback func() T) T {
	return WhenInLifecycles([]LifecycleStage{Build}, callback)
}

// IsRunning checks if the current lifecycle is running the app (LocalRun or Cloud)
func IsRunning() bool {
	return IsInLifecycle(LocalRun, Cloud)
}

// IsCollecting checks if the current lifecycle is collecting application requirements (Build)
func IsCollecting() bool {
	return IsInLifecycle(Build)
}
