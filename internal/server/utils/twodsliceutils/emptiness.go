package twodsliceutils

func Is2DSliceEmpty(slice [][]int) bool {
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
