package utils

type Integer interface {
	int | int64 | int32
}

func CountDigits[I Integer](x I) I {
	var length I = 0
	for x != 0 {
		x /= 10
		length++
	}

	return length
}
