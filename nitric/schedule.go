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
	"strings"

	"github.com/nitrictech/go-sdk/handler"
	"github.com/nitrictech/go-sdk/workers"
	schedulespb "github.com/nitrictech/nitric/core/pkg/proto/schedules/v1"
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

// Run middleware at a certain interval defined by the cronExpression.
func (s *schedule) Cron(cron string, middleware ...handler.IntervalMiddleware) {
	scheduleCron := &schedulespb.ScheduleCron{
		Expression: cron,
	}

	registrationRequest := &schedulespb.RegistrationRequest{
		ScheduleName: s.name,
		Cadence: &schedulespb.RegistrationRequest_Cron{
			Cron: scheduleCron,
		},	
	}

	composeHandler := handler.ComposeIntervalMiddleware(middleware...)

	opts := &workers.IntervalWorkerOpts{
		RegistrationRequest: registrationRequest,
		Middleware: composeHandler,	
	}

	worker := workers.NewIntervalWorker(opts)
	s.manager.addWorker("IntervalWorkerCron:" + strings.Join([]string{
		s.name,
		cron,
	}, "-"), worker)
}

// Run middleware at a certain interval defined by the rate. The rate is e.g. '7 days'. All rates accept a number and a frequency. Valid frequencies are 'days', 'hours' or 'minutes'.
func (s *schedule) Every(rate string, middleware ...handler.IntervalMiddleware) {
	scheduleEvery := &schedulespb.ScheduleEvery{
		Rate: rate,
	}

	registrationRequest := &schedulespb.RegistrationRequest{
		ScheduleName: s.name,
		Cadence: &schedulespb.RegistrationRequest_Every{
			Every: scheduleEvery,
		},
	}

	composeHandler := handler.ComposeIntervalMiddleware(middleware...)

	opts := &workers.IntervalWorkerOpts{
		RegistrationRequest: registrationRequest,
		Middleware: composeHandler,	
	}

	worker := workers.NewIntervalWorker(opts)
	s.manager.addWorker("IntervalWorkerEvery:" + strings.Join([]string{
		s.name,
		rate,
	}, "-"), worker)
}
