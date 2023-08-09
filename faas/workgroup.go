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

package faas

import (
	"sync"

	multierror "github.com/missionMeteora/toolkit/errors"
)

type workPool struct {
	wg       sync.WaitGroup
	maxGuard chan bool
	ErrorCh  chan error
	err      *multierror.ErrorList
}

func NewWorkPool(maxWorkers int) *workPool {
	return &workPool{
		wg:       sync.WaitGroup{},
		maxGuard: make(chan bool, maxWorkers),
		ErrorCh:  make(chan error, maxWorkers),
		err:      &multierror.ErrorList{},
	}
}

func (w *workPool) Go(f func(a interface{}) error, arg interface{}) {
	w.addOrBlock()

	go func(o interface{}) {
		defer w.done()

		if err := f(o); err != nil {
			w.err.Push(err)
		}
	}(arg)
}

func (w *workPool) AddError(e error) {
	w.err.Push(e)
}

func (w *workPool) Err() error {
	return w.err.Err()
}

func (w *workPool) Wait() {
	w.wg.Wait()
}

func (w *workPool) addOrBlock() {
	w.maxGuard <- true
	w.wg.Add(1)
}

func (w *workPool) done() {
	<-w.maxGuard
	w.wg.Done()
}
