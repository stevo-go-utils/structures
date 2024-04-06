package structures

func AnyArr[T any](arr []T) (res []any) {
	for _, v := range arr {
		res = append(res, v)
	}
	return
}

func ParseAnyArr[T any](arr []any) (res []T) {
	for _, v := range arr {
		switch v := v.(type) {
		case T:
			res = append(res, v)
		}
	}
	return
}
