package main

import (
	"fmt"
	"log"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

type Worker interface {
	Run() error
	GetInterval() time.Duration
}

type Task struct {
	Id       string
	Job      Worker
	Interval time.Duration
	stopChan chan struct{}
}

type Scheduler struct {
	Tasks       map[string]Task
	newTaskChan chan Task
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		Tasks:       make(map[string]Task),
		newTaskChan: make(chan Task),
	}
}

func (self *Scheduler) NewTask(job Worker) Task {
	task := Task{
		Id:       uuid.NewUUID().String(),
		Job:      job,
		Interval: job.GetInterval(),
		stopChan: make(chan struct{}),
	}

	self.Tasks[task.Id] = task

	self.newTaskChan <- task

	return task
}

func (self *Scheduler) RemoveTask(id string) error {
	task, ok := self.Tasks[id]

	if !ok {
		return fmt.Errorf("unknown worker: %s\n", id)
	}

	task.stopChan <- struct{}{}

	return nil
}

func (self *Scheduler) Start() {
	go func() {
		for task := range self.newTaskChan {
			go func() {
				for {
					select {
					case <-task.stopChan:
						return

					default:
						log.Printf("running worker: %s\n", task.Id)

						if err := task.Job.Run(); err != nil {
							log.Printf("error running worker: %s: %s\n", task.Id, err)
						}

						time.Sleep(task.Interval)
					}
				}
			}()
		}
	}()
}
