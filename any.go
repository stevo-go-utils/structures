package structures

func AnyArr[T any](arr []T) (res []any) {
	for _, v := range arr {
		res = append(res, v)
	}
	return
}
