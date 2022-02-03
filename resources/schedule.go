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

package resources

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nitrictech/go-sdk/faas"
)

var singularRates = []string{"minute", "hour", "day"}

func rateSplit(rate string) (int, faas.Frequency, error) {
	rateParts := strings.Split(rate, " ")
	if len(rateParts) == 1 {
		for _, r := range singularRates {
			if r == rateParts[0] {
				return 1, faas.Frequency(r + "s"), nil
			}
		}
	}
	if len(rateParts) != 2 {
		return 0, "", fmt.Errorf("not enough parts to rate expression %s", rate)
	}
	rateNum := rateParts[0]
	rateType := rateParts[1]

	num, err := strconv.Atoi(rateNum)
	if err != nil {
		return 0, "", fmt.Errorf("invalid rate expression %s; %w", rate, err)
	}

	for _, r := range singularRates {
		if r+"s" == rateType {
			return num, faas.Frequency(rateType), nil
		}
	}
	return 0, "", fmt.Errorf("invalid rate expression %s; %s must be one of [minutes, hours, days]", rate, rateType)
}

func NewSchedule(name, rate string, handlers ...faas.EventMiddleware) error {
	return run.NewSchedule(name, rate, handlers...)
}

func (m *manager) NewSchedule(name, rate string, handlers ...faas.EventMiddleware) error {
	f, ok := m.builders[name]
	if !ok {
		f = faas.New()
	}

	r, freq, err := rateSplit(rate)
	if err != nil {
		return err
	}

	f.Event(handlers...).
		WithRateWorkerOpts(faas.RateWorkerOptions{
			Description: name,
			Rate:        r,
			Frequency:   freq,
		})

	m.addStarter(fmt.Sprintf("schedule:%s/%s", name, rate), f)
	m.builders[name] = f

	return nil
}
