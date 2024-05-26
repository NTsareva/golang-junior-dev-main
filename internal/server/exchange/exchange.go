package exchange

import (
	"context"
	"errors"
	"sync"
	"time"

	"gitlab.com/llcmediatel/recruiting/golang-junior-dev/internal/server/processor"
)

type Input struct {
	Amount    int   `json:"amount"`
	Banknotes []int `json:"banknotes"`
}

type Output struct {
	Exchange [][]int `json:"exchange"`
}

// PresentResults(input Input) Output converts results type [][]int tahen from Input values to Output
func PresentResults(input Input) (Output, error) {
	resultCh := make(chan []int)
	var wg sync.WaitGroup

	//cancel context is made for situation where our amount is too large so we can load our machine with processing
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	wg.Add(1)

	go processor.Process(ctx, input.Amount, input.Banknotes, 0, []int{}, resultCh, &wg)

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var results [][]int

	for combination := range resultCh {
		results = append(results, combination)
	}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return Output{}, ctx.Err()
	}
	return Output{Exchange: results}, nil
}
