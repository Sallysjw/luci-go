// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package parallel

import (
	"sync"

	"github.com/luci/luci-go/common/errors"
)

// WorkPool creates a fixed-size pool of worker goroutines. A supplied generator
// method creates task functions and passes them through to the work pool.
// Available workers will consume tasks from the pool and execute them until the
// generator is finished.
//
// WorkPool blocks until all the generator completes and all workers have
// finished their tasks.
func WorkPool(workers int, gen func(chan<- func() error)) error {
	if workers < 0 {
		return errors.New("invalid number of workers")
	}

	sem := make(Semaphore, workers)
	errchan := make(chan error, workers)
	funchan := make(chan func() error, workers)

	go func() {
		defer close(funchan)
		gen(funchan)
	}()

	go func() {
		grp := sync.WaitGroup{}

		for fn := range funchan {
			sem.Lock()

			grp.Add(1)
			fn := fn

			go func() {
				defer func() {
					sem.Unlock()
					grp.Done()
				}()

				errchan <- fn()
			}()
		}

		grp.Wait()
		close(errchan)
	}()

	return errors.MultiErrorFromErrors(errchan)
}
