package processor

import (
	"context"
	"sync"
)

// Function Process(amount int, banknotes []int, init int, curentCombiation []int, result *[][]int) processes current
// amount of money (amount) with given banknotes (banknotes []int) and out all combination of exchange (result *[][]int)
// init int shows from what index of banknotes shoul we start change
// curentCombiation []int show current combinetions of banknotes we already have for exchange
func Process(ctx context.Context, amount int, banknotes []int, init int, curentCombiation []int, resultCh chan<- []int, wg *sync.WaitGroup) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		return
	default:
	}

	if amount == 0 {
		combination := make([]int, len(curentCombiation))
		copy(combination, curentCombiation)
		resultCh <- combination
		return
	}

	for i := init; i < len(banknotes); i++ {
		if amount < banknotes[i] {
			continue
		}

		newCombination := append(curentCombiation, banknotes[i])

		wg.Add(1)

		countedAmount := amount - banknotes[i]
		go Process(ctx, countedAmount, banknotes, i, newCombination, resultCh, wg)
	}
}
