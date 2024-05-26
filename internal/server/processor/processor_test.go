package processor

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name      string
		amount    int
		banknotes []int
		expected  [][]int
		byTimeout bool
	}{
		{
			name:      "Average amount more than biggest banknote",
			amount:    300,
			banknotes: []int{5000, 2000, 1000, 500, 200, 100, 50},
			expected: [][]int{
				{200, 100},
				{200, 50, 50},
				{100, 100, 100},
				{100, 100, 50, 50},
				{100, 50, 50, 50, 50},
				{50, 50, 50, 50, 50, 50},
			},
			byTimeout: false,
		},
		{
			name:      "Zero amount",
			amount:    0,
			banknotes: []int{500, 200, 100, 50},
			expected:  [][]int{},
			byTimeout: false,
		},
		{
			name:      "Negative amount",
			amount:    -300,
			banknotes: []int{200, 100, 50},
			expected:  [][]int{},
			byTimeout: false,
		}, {
			name:      "Amount of minimum banknote",
			amount:    50,
			banknotes: []int{200, 100, 50},
			expected: [][]int{
				{50},
			},
			byTimeout: false,
		}, {
			name:      "Amount less than minimum banknote",
			amount:    50,
			banknotes: []int{500, 200, 100},
			expected:  [][]int{},
			byTimeout: false,
		}, {
			name:      "No banknotes",
			amount:    300,
			banknotes: []int{},
			expected:  [][]int{},
			byTimeout: false,
		},
		{
			name:      "Finish by timeout",
			amount:    1000000,
			banknotes: []int{5000, 2000, 1000, 500, 200, 100, 50, 5, 1},
			expected:  [][]int{},
			byTimeout: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result [][]int
			resultChan := make(chan []int)
			var wg sync.WaitGroup

			ctx := context.Background()
			var cancel context.CancelFunc

			if test.byTimeout {
				ctx, cancel = context.WithTimeout(ctx, 100*time.Millisecond)
				defer cancel()
			}

			wg.Add(1)
			go Process(ctx, test.amount, test.banknotes, 0, []int{}, resultChan, &wg)

			go func() {
				wg.Wait()
				close(resultChan)
			}()

			for combinetion := range resultChan {
				result = append(result, combinetion)
			}

			sort2D(result)
			sort2D(test.expected)

			if test.byTimeout && !errors.Is(ctx.Err(), context.DeadlineExceeded) {
				t.Errorf("expected deadline but got %v", ctx.Err())
			}

			if is2DSliceEmpty(test.expected) && is2DSliceEmpty(result) {
				return
			}

			if !test.byTimeout && !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected %v but got %v", test.expected, result)
			}
		})
	}
}

func sort2D(slices [][]int) {
	for _, s := range slices {
		sort.Ints(s)
	}

	sort.Slice(slices, func(i, j int) bool {
		for k := 0; k < len(slices[i]) && k < len(slices[j]); k++ {
			if slices[i][k] != slices[j][k] {
				return slices[i][k] < slices[j][k]
			}
		}

		return len(slices[i]) < len(slices[j])
	})
}

func is2DSliceEmpty(slice [][]int) bool {
	if slice == nil || len(slice) == 0 {
		return true
	}

	for _, innerSlice := range slice {
		if len(innerSlice) != 0 {
			return false
		}
	}
	return true
}
