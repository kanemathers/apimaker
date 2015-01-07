package main

import (
	"code.google.com/p/go-uuid/uuid"
)

type Scheduler struct {
	Jobs map[string]*Job
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		Jobs: make(map[string]*Job),
	}
}

func (self *Scheduler) AddJob(job *Job) string {
	id := uuid.NewUUID().String()
	self.Jobs[id] = job

	return id
}
