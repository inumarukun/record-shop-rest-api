package common

// 任意の型のスライスを変換する汎用関数
func MapSlice[T, U any](input []T, transformFunc func(T) U) []U {
	var result []U
	for _, item := range input {
		result = append(result, transformFunc(item))
	}
	return result
}
