package services

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	jobFun interface{}
}

// NewScheduler creates a scheduler
func NewScheduler(intervalSecs int, job interface{}) *Scheduler {
	s := &Scheduler{
		jobFun: job,
	}

	go s.runScheduler(intervalSecs)

	return s
}

func (s *Scheduler) runScheduler(intervalSecs int) {
	scheduler := gocron.NewScheduler(time.UTC)
	_, err := scheduler.Every(intervalSecs).Second().Do(s.jobFun)
	if err != nil {
		panic(fmt.Errorf("failed to run sscheduled job: %v", err))
	}

	scheduler.StartBlocking()
}
