package main

import (
	"sync"
)

type Job interface {
	DoJob()
	GetJobId() int
}

type WorkerPool struct {
	jobQueue   chan Job
	workers    chan Worker
	maxWorkers int
	jobWg      sync.WaitGroup
	workerWg   sync.WaitGroup
}

func NewWorkerPool(maxQueue int, maxWorkers int) *WorkerPool {
	return &WorkerPool{
		jobQueue:   make(chan Job, maxQueue),
		workers:    make(chan Worker, maxWorkers),
		maxWorkers: maxWorkers,
	}
}

func (w *WorkerPool) WorkOnJob(jobList []Job) {
	w.run()

	for i := 0; i < len(jobList); i++ {
		w.jobWg.Add(1)
		w.jobQueue <- jobList[i]
	}
}

func (w *WorkerPool) WaitAndClose() {
	w.jobWg.Wait()
	close(w.jobQueue)
	w.workerWg.Wait()
	close(w.workers)
}

type Worker struct {
	id         int
	jobChannel chan Job
	quit       chan bool
}

func (w *WorkerPool) run() {
	for i := 0; i < w.maxWorkers; i++ {
		worker := newWorker(i)
		w.workerStart(worker)
	}
	w.dispatch()
}

func (w *WorkerPool) dispatch() {
	go func() {
		for {
			select {
			case job, ok := <-w.jobQueue:
				if ok {
					worker := <-w.workers
					worker.jobChannel <- job
				} else { //jobQueue has been closed, close worker goroutines
					select {
					case worker := <-w.workers:
						worker.quit <- true
					default: //no more worker in the workers chan, ready to return
						return
					}
				}
			}
		}
	}()
}

func newWorker(id int) Worker {
	return Worker{
		id,
		make(chan Job),
		make(chan bool),
	}
}

func (w *WorkerPool) workerStart(worker Worker) {
	w.workerWg.Add(1)
	go func() {
		defer w.workerWg.Done()

		for {
			// worker is available, ready to pick up a new job
			w.workers <- worker
			select {
			case job := <-worker.jobChannel:
				// log.Printf("worker %d picks up job %d", worker.id, job.GetJobId())
				job.DoJob()
				w.jobWg.Done()
			case <-worker.quit:
				return
			}
		}
	}()
}

func (w Worker) stop() {
	go func() {
		w.quit <- true
	}()
}
