// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package constants

import (
	"os"
	"strconv"
	"time"
)

const (
	nitricServiceHostDefault        = "127.0.0.1"
	nitricServicePortDefault        = "50051"
	nitricServiceDialTimeoutDefault = "5000"
	nitricServiceAddress            = "SERVICE_ADDRESS"
)

// getEnvWithFallback - Returns an envirable variable's value from its name or a default value if the variable isn't set
func GetEnvWithFallback(varName string, fallback string) string {
	if v := os.Getenv(varName); v == "" {
		return fallback
	} else {
		return v
	}
}

// NitricDialTimeout - Retrieves default service dial timeout in milliseconds
func NitricDialTimeout() time.Duration {
	tInt, _ := strconv.ParseInt(nitricServiceDialTimeoutDefault, 10, 64)

	return time.Duration(tInt) * time.Millisecond
}

// nitricAddress - constructs the full address i.e. host:port, of the nitric service based on config or defaults
func NitricAddress() string {
	return GetEnvWithFallback(nitricServiceAddress, nitricServiceHostDefault)
}
