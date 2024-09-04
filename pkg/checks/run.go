package checks

import (
	"sync"
)

// RunChecksInParallel for the given input and handle errors with the callback provided
// Output is provided in a unbuffered channel that is closed once all checks finished
func RunChecksInParallel(checks []Check, checkInput CheckInput, errHandler func(error)) chan []CheckFinding {
	wg := sync.WaitGroup{}
	checkResults := make(chan []CheckFinding)
	for _, check := range checks {
		wg.Add(1)
		go func() {
			defer wg.Done()

			results, err := check.Run(&checkInput)
			if err != nil {
				errHandler(err)
			}
			checkResults <- results
		}()
	}

	go func() {
		wg.Wait()
		close(checkResults)
	}()
	return checkResults
}
