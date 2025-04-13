package parsex

import (
	"errors"
	"sync"
)

// Batches executable functions together.
// Calls them sequentially and returns on the first encountered error.
func BatchSeq(execs ...Executable) Executable {
	var err error
	return func(args []string) error {
		for _, function := range execs {
			err = function(args)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// Batches executable functions together.
// Calls them in parallel and returns the combined error.
func Batch(execs ...Executable) Executable {
	return func(args []string) error {
		var wg sync.WaitGroup
		var mutex sync.Mutex
		wg.Add(len(execs))

		var errs []error

		for _, function := range execs {
			go func(wg *sync.WaitGroup, args []string) error {
				defer wg.Done()
				err := function(args)

				mutex.Lock()
				errs = append(errs, err)
				mutex.Unlock()
				return nil
			}(&wg, args)
		}

		wg.Wait()

		return errors.Join(errs...)
	}
}
