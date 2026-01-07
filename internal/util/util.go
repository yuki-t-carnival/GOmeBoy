package util

// arr[idx]にvalを挿入する（上書きせずずらす）
func InsertSlice(arr []int, idx, val int) []int {
	arr = append(arr, 0)         // 長さを1伸ばす
	copy(arr[idx+1:], arr[idx:]) // idx以降の内容を､idx+1以降にコピー（=右に1ずらす）
	arr[idx] = val
	return arr
}
