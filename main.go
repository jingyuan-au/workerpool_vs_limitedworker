package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

type MyJob struct {
	JobId int
}

func (o MyJob) DoJob() {
	fmt.Printf("Processing job: %d ...\r\n", o.JobId)
	time.Sleep(1 * time.Millisecond)
}

func (o MyJob) GetJobId() int {
	return o.JobId
}

func main() {

	jobList := getJobList(1000000)

	log.Println("Satrt workerpool...")
	workerpoolRuntime := runJobWithWorkerpool(jobList)

	log.Println("Satrt limitedWorker...")
	limitedWorkerRuntime := runJobWithLimitedWorker(jobList)

	log.Println("Wokerpool pattern Total time:", workerpoolRuntime)
	log.Println("LimitedWorker pattern Total time:", limitedWorkerRuntime)

	//check leak
	log.Println("go routines:", runtime.NumGoroutine())
}

func getJobList(length int) []Job {
	jobList := make([]Job, length)
	for i := 0; i < length; i++ {
		jobList[i] = MyJob{i}
	}
	return jobList
}

func runJobWithWorkerpool(jobList []Job) time.Duration {
	start := time.Now()

	workerPool := NewWorkerPool(100, 1000)
	workerPool.WorkOnJob(jobList)
	workerPool.WaitAndClose()

	end := time.Now()
	return end.Sub(start)
}

func runJobWithLimitedWorker(jobList []Job) time.Duration {
	start := time.Now()

	limitedWorker := NewLimitedWorker(1000)
	limitedWorker.WorkOnJob(jobList)
	limitedWorker.WaitAndClose()
	end := time.Now()
	return end.Sub(start)
}
