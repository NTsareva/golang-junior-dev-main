package twodsliceutils

import "sort"

func Sort2D(slices [][]int) {
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
