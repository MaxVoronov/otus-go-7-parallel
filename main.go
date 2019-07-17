package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const TotalJobs = 500
const MaxWorkingTime = 10
const MaxParallelJobs = 50
const JobsErrorLimit = 25

func main() {
	// Generate slice of random jobs
	jobs := make([]func() error, 0, TotalJobs)
	for i := 0; i < TotalJobs; i++ {
		jobs = append(jobs, someUsefulWork)
	}

	// Start processing of jobs
	fmt.Printf("Prepare processing of %d jobs...\n", len(jobs))
	Run(jobs)
}

func Run(fnList []func() error) {
	var wg sync.WaitGroup
	workers := make(chan struct{}, MaxParallelJobs)
	jobErrors := make(chan error, JobsErrorLimit)
	errorCounter := &ErrorCounter{}

	// Read errors from channel and increase counter
	go func(counter *ErrorCounter) {
		for {
			select {
			case <-jobErrors:
				counter.Increase()
			}
		}
	}(errorCounter)

	for i, f := range fnList {
		if !errorCounter.Less(JobsErrorLimit) {
			break
		}

		workers <- struct{}{}
		wg.Add(1)

		fmt.Printf("Starting new job #%d\n", i)
		go func(fn func() error, idx int) {
			defer wg.Done()
			if err := fn(); err != nil {
				fmt.Printf("Failed to finish job #%d\n", idx)
				jobErrors <- err
			} else {
				fmt.Printf("Job #%d was succesufully done\n", idx)
			}
			<-workers
		}(f, i)
	}

	wg.Wait()
}

func someUsefulWork() error {
	rand.Seed(time.Now().UnixNano())
	duration := rand.Intn(MaxWorkingTime * 1000)
	time.Sleep(time.Duration(duration) * time.Millisecond)

	// Every 3rd job will fail
	if rand.Intn(100)%3 == 0 {
		return errors.New("Ops! Something went wrong")
	}

	return nil
}
