package main

import "sync"

type LimitedWorker struct {
	semaphore chan struct{}
	wg        sync.WaitGroup
}

func NewLimitedWorker(maxWorkers int) *LimitedWorker {
	return &LimitedWorker{
		semaphore: make(chan struct{}, maxWorkers),
	}
}

func (w *LimitedWorker) WorkOnJob(jobList []Job) {
	for _, job := range jobList {
		w.wg.Add(1)
		w.semaphore <- struct{}{}
		go func(j Job) {
			j.DoJob()
			w.wg.Done()
			<-w.semaphore
		}(job)
	}
}

func (w *LimitedWorker) WaitAndClose() {
	w.wg.Wait()
	close(w.semaphore)
}
