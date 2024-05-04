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

package nitric

import (
	"github.com/nitrictech/go-sdk/handler"
)

type Schedule interface {
	Cron(cron string, middleware ...handler.IntervalMiddleware)
	Every(rate string, middleware ...handler.IntervalMiddleware)
}

type schedule struct {
	Schedule

	name    string
	manager Manager
}

// NewSchedule provides a new schedule, which can be configured with a rate/cron and a callback to run on the schedule.
func NewSchedule(name string) Schedule {
	return defaultManager.newSchedule(name)
}

func (m *manager) newSchedule(name string) Schedule {
	return &schedule{
		name:    name,
		manager: m,
	}
}

func (s *schedule) Cron(cron string, middleware ...handler.IntervalMiddleware) {
}

// The rate is e.g. '7 days'. All rates accept a number and a frequency. Valid frequencies are 'days', 'hours' or 'minutes'.
func (s *schedule) Every(rate string, middleware ...handler.IntervalMiddleware) {
	// TODO: create schedule worker
}
