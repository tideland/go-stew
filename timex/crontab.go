// Tideland Go Stew - Time Extensions
//
// Copyright (C) 2009-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package timex // import "tideland.dev/go/stew/timex"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"fmt"
	"time"

	"tideland.dev/go/stew/loop"
)

//--------------------
// CONSTANTS
//--------------------

// Status defines the status of a crontab.
type Status = loop.Status

const (
	StatusStopped    Status = loop.StatusStopped
	StatusStarting          = loop.StatusStarting
	StatusWorking           = loop.StatusWorking
	StatusRepairing         = loop.StatusRepairing
	StatusStopping          = loop.StatusStopping
	StatusFinalizing        = loop.StatusFinalizing
	StatusError             = loop.StatusError
)

//--------------------
// CRONTAB
//--------------------

// Job is executed by the crontab.
type Job func() (bool, error)

// crontab is the internal type for the cron server.
type cronjob struct {
	id        string
	frequency time.Duration
	last      time.Time
	job       Job
}

// Crontab is one cron server. A system can run multiple ones
// in parallel.
type Crontab struct {
	ctx        context.Context
	interval   time.Duration
	jobs       map[string]*cronjob
	terminated map[string]error
	addCh      chan *cronjob
	removeCh   chan string
	loop       *loop.Loop
}

// NewCrontab creates a cron server.
func NewCrontab(ctx context.Context, interval time.Duration) (*Crontab, error) {
	c := &Crontab{
		ctx:        ctx,
		interval:   interval,
		jobs:       make(map[string]*cronjob),
		terminated: make(map[string]error),
		addCh:      make(chan *cronjob, 1),
		removeCh:   make(chan string, 1),
	}
	l, err := loop.Go(c.worker, loop.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("error starting backend worker: %v", err)
	}
	c.loop = l
	return c, nil
}

// Stop terminates the cron server.
func (c *Crontab) Stop() error {
	c.loop.Stop()
	return c.loop.Err()
}

// Status returns the current status of the cron server.
func (c *Crontab) Status() loop.Status {
	return c.loop.Status()
}

// Add adds a new job to the server.
func (c *Crontab) Add(id string, frequency time.Duration, job Job) {
	if frequency < c.interval {
		frequency = c.interval
	}
	cj := &cronjob{
		id:        id,
		frequency: frequency,
		last:      time.Now(),
		job:       job,
	}
	c.addCh <- cj
}

// Remove removes a job from the server.
func (c *Crontab) Remove(id string) {
	c.removeCh <- id
}

// JobStatus returns if a job is still active or if it possibly
// terminated with an error.
func (c *Crontab) JobStatus(id string) (bool, error) {
	_, ok := c.jobs[id]
	if ok {
		// Still active.
		return true, nil
	}
	err, ok := c.terminated[id]
	if ok {
		// Terminated with error.
		return false, err
	}
	// Not found.
	return false, nil
}

// worker runs the server backend.
func (c *Crontab) worker(ctx context.Context) error {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case addJob := <-c.addCh:
			c.jobs[addJob.id] = addJob
		case id := <-c.removeCh:
			delete(c.jobs, id)
		case now := <-ticker.C:
			for id, job := range c.jobs {
				c.handleJob(id, job, now)
			}
		}
	}
}

// handleJob checks if a job shall be executed and starts it as goroutine
// if yes.
func (c *Crontab) handleJob(id string, cj *cronjob, now time.Time) {
	if cj.last.Add(cj.frequency).Before(now) {
		cj.last = now
		go func() {
			cont, err := cj.job()
			if err != nil {
				c.terminated[id] = err
				cont = false
			}
			if !cont {
				c.Remove(id)
			}
		}()
	}
}

// EOF
