package main

import (
	"fmt"
	"log"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

type Scheduler struct {
	Jobs map[string]*Job

	newJobChan chan *Job
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		Jobs:       make(map[string]*Job),
		newJobChan: make(chan *Job),
	}
}

func (self *Scheduler) AddJob(job *Job) string {
	job.Id = uuid.NewUUID().String()
	self.Jobs[job.Id] = job

	self.newJobChan <- job

	return job.Id
}

func (self *Scheduler) Start() {
	for job := range self.newJobChan {
		duration, err := time.ParseDuration(fmt.Sprintf("%ds", job.Interval))

		if err != nil {
			log.Printf("error parsing duration for job: %s\n", job.Id)
		}

		go func() {
			job.Scrape()

			time.Sleep(duration)
		}()
	}
}
