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

func main() {
	// Generate slice of random jobs
	jobs := make([]func() error, 0, TotalJobs)
	for i := 0; i < TotalJobs; i++ {
		jobs = append(jobs, someUsefulWork)
	}

	// Start processing of jobs
	fmt.Printf("Prepare processing of %d jobs...\n", len(jobs))
	Worker(jobs)
	fmt.Println("Done!")
}

func Worker(fnList []func() error) {
	var wg sync.WaitGroup

	for i, fn := range fnList {
		fmt.Printf("Starting job #%d\n", i+1)
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if err := fn(); err != nil {
				fmt.Printf("Failed to finish job #%d\n", idx)
				return
			}
			fmt.Printf("Job #%d done\n", idx)
		}(i + 1)
	}

	wg.Wait()
}

func someUsefulWork() error {
	rand.Seed(time.Now().UnixNano())
	duration := rand.Intn(MaxWorkingTime * 1000)
	time.Sleep(time.Duration(duration) * time.Millisecond)

	// Every 3rd job will fail
	if rand.Intn(100)%3 == 0 {
		return errors.New("Ops! Some error")
	}

	return nil
}
