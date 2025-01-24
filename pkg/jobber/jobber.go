package jobber

import (
	"orbital/pkg/stringer"
	"sync"
	"time"
)

const (
	MaxRunInfinte = -1
)

type Task func()

type Runner struct {
	mu    sync.Mutex
	jobs  map[string]*job
	pool  chan struct{}
	close chan struct{}
}

type job struct {
	id       string
	interval time.Duration
	task     Task
	stop     chan struct{}
	maxRuns  int
	runCount int
}

var jobCounter uint64

func New(workerCount int) *Runner {
	return &Runner{
		jobs:  make(map[string]*job),
		pool:  make(chan struct{}, workerCount),
		close: make(chan struct{}),
	}
}

func (r *Runner) AddJob(interval time.Duration, maxRun int, task Task) string {
	r.mu.Lock()
	defer r.mu.Unlock()

	if maxRun == 0 {
		maxRun = MaxRunInfinte
	}

	id, _ := stringer.Random(16, stringer.RandNumber, stringer.RandLowercase)
	j := &job{
		id:       id,
		interval: interval,
		task:     task,
		stop:     make(chan struct{}),
		maxRuns:  maxRun,
		runCount: 0,
	}

	r.jobs[id] = j

	go r.runJob(j)
	return id
}

func (r *Runner) RemoveJob(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if j, exists := r.jobs[id]; exists {
		close(j.stop)
		delete(r.jobs, id)
	}
}

func (r *Runner) Shutdown() {
	close(r.close)
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, j := range r.jobs {
		close(j.stop)
	}

	r.jobs = make(map[string]*job)
}

func (r *Runner) runJob(j *job) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			select {
			case r.pool <- struct{}{}:
				j.task()
				j.runCount++

				if j.maxRuns > 0 && j.runCount >= j.maxRuns {
					r.RemoveJob(j.id)
					return
				}

				<-r.pool
			case <-j.stop:
				return
			}
		case <-j.stop:
			return
		case <-r.close:
			return
		}
	}
}
