package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const TotalJobs = 500
const MaxWorkingTime = 10
const MaxParallelJobs = 10
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
	jobs := make(chan func() error, MaxParallelJobs)
	done := make(chan struct{})
	errorCounter := &ErrorCounter{}

	for i := 0; i < MaxParallelJobs; i++ {
		go func(counter *ErrorCounter, idx int) {
			fmt.Printf("[#%d] Worker started \n", idx)
			for {
				select {
				case fn := <-jobs:
					fmt.Printf("[#%d] Start processing of new job\n", idx)
					if err := fn(); err != nil {
						counter.Increase()
						fmt.Printf("[#%d] Failed to process job\n", idx)
					} else {
						fmt.Printf("[#%d] Successfull done\n", idx)
					}
				case <-done:
					fmt.Printf("[#%d] Got exit signal\n", idx)
					return
				}
			}
		}(errorCounter, i)
	}

	for _, job := range fnList {
		if !errorCounter.Less(JobsErrorLimit) {
			fmt.Printf("--== Max error limit: %d ==--\n", errorCounter.Value())
			done <- struct{}{}
			break
		}

		fmt.Println("<-- Sent job to worker")
		jobs <- job
	}

	fmt.Println("Done!")
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
