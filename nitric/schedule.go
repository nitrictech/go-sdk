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
	"fmt"
	"strconv"
	"strings"

	"github.com/nitrictech/go-sdk/faas"
)

type Schedule interface {
	Cron(cron string, middleware ...faas.EventMiddleware)
	Every(rate string, middleware ...faas.EventMiddleware) error
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

func (s *schedule) Cron(cron string, middleware ...faas.EventMiddleware) {
	f := s.manager.getBuilder(s.name)
	if f == nil {
		f = faas.New()
	}

	f.Event(middleware...).
		WithCronWorkerOpts(faas.CronWorkerOptions{
			Description: s.name,
			Cron:        cron,
		})

	s.manager.addWorker(fmt.Sprintf("schedule:%s/%s", s.name, cron), f)
	s.manager.addBuilder(s.name, f)
}

func rateSplit(rate string) (int, faas.Frequency, error) {
	rateParts := strings.Split(rate, " ")

	if len(rateParts) < 1 || len(rateParts) > 2 {
		return -1, "", fmt.Errorf("invalid rate expression %s; rate should be in the form '[rate] [frequency]' e.g. '7 days'", rate)
	}

	// Handle a single rate e.g. 'day'
	if len(rateParts) == 1 {
		for _, r := range []string{"minute", "hour", "day"} {
			if r == rateParts[0] {
				return 1, faas.Frequency(r + "s"), nil
			}
		}
	}

	// Handle a full rate expression e.g. '7 days'
	rateNum := rateParts[0]
	rateType := rateParts[1]

	num, err := strconv.Atoi(rateNum)
	if err != nil {
		return -1, "", fmt.Errorf("invalid rate expression %s; %w", rate, err)
	}

	for _, r := range faas.Frequencies {
		if string(r) == rateType {
			return num, faas.Frequency(rateType), nil
		}
	}

	return -1, "", fmt.Errorf("invalid rate expression %s; %s must be one of [minutes, hours, days]", rate, rateType)
}

// The rate is e.g. '7 days'. All rates accept a number and a frequency. Valid frequencies are 'days', 'hours' or 'minutes'.
func (s *schedule) Every(rate string, middleware ...faas.EventMiddleware) error {
	f := s.manager.getBuilder(s.name)
	if f == nil {
		f = faas.New()
	}

	rateNum, frequency, err := rateSplit(rate)
	if err != nil {
		return err
	}

	f.Event(middleware...).
		WithRateWorkerOpts(faas.RateWorkerOptions{
			Description: s.name,
			Frequency:   frequency,
			Rate:        rateNum,
		})

	s.manager.addBuilder(s.name, f)
	s.manager.addWorker(fmt.Sprintf("schedule:%s/%s", s.name, rate), f)

	return nil
}
