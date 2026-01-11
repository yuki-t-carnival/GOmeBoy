package util

// Insert int into []int
func InsertSlice(arr []int, idx, val int) []int {
	arr = append(arr, 0)         // Extend the destination slice by one
	copy(arr[idx+1:], arr[idx:]) // Shift elements from idx onward to idx+1
	arr[idx] = val
	return arr
}
