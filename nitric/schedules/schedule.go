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

package schedules

import (
	"strings"

	"github.com/nitrictech/go-sdk/internal/handlers"
	"github.com/nitrictech/go-sdk/nitric/workers"
	schedulespb "github.com/nitrictech/nitric/core/pkg/proto/schedules/v1"
)

type Schedule interface {
	// Run a function at a certain interval defined by the cronExpression.
	// Valid function signatures for handler are:
	//
	//	func()
	//	func() error
	//	func(*schedules.Ctx)
	//	func(*schedules.Ctx) error
	//	Handler[schedules.Ctx]
	Cron(cron string, handler interface{})

	// Run a function at a certain interval defined by the rate. The rate is e.g. '7 days'. All rates accept a number and a frequency. Valid frequencies are 'days', 'hours' or 'minutes'.
	// Valid function signatures for handler are:
	//
	//	func()
	//	func() error
	//	func(*schedules.Ctx)
	//	func(*schedules.Ctx) error
	//	Handler[schedules.Ctx]
	Every(rate string, handler interface{})
}

type schedule struct {
	name    string
	manager *workers.Manager
}

var _ Schedule = (*schedule)(nil)

// NewSchedule - Create a new Schedule resource
func NewSchedule(name string) Schedule {
	return &schedule{
		name:    name,
		manager: workers.GetDefaultManager(),
	}
}

func (s *schedule) Cron(cron string, handler interface{}) {
	scheduleCron := &schedulespb.ScheduleCron{
		Expression: cron,
	}

	registrationRequest := &schedulespb.RegistrationRequest{
		ScheduleName: s.name,
		Cadence: &schedulespb.RegistrationRequest_Cron{
			Cron: scheduleCron,
		},
	}

	typedHandler, err := handlers.HandlerFromInterface[Ctx](handler)
	if err != nil {
		panic(err)
	}

	opts := &scheduleWorkerOpts{
		RegistrationRequest: registrationRequest,
		Handler:             typedHandler,
	}

	worker := newScheduleWorker(opts)
	s.manager.AddWorker("IntervalWorkerCron:"+strings.Join([]string{
		s.name,
		cron,
	}, "-"), worker)
}

func (s *schedule) Every(rate string, handler interface{}) {
	scheduleEvery := &schedulespb.ScheduleEvery{
		Rate: rate,
	}

	registrationRequest := &schedulespb.RegistrationRequest{
		ScheduleName: s.name,
		Cadence: &schedulespb.RegistrationRequest_Every{
			Every: scheduleEvery,
		},
	}

	typedHandler, err := handlers.HandlerFromInterface[Ctx](handler)
	if err != nil {
		panic(err)
	}

	opts := &scheduleWorkerOpts{
		RegistrationRequest: registrationRequest,
		Handler:             typedHandler,
	}

	worker := newScheduleWorker(opts)
	s.manager.AddWorker("IntervalWorkerEvery:"+strings.Join([]string{
		s.name,
		rate,
	}, "-"), worker)
}
