package utils

func SliceEqual[T comparable](as, bs []T) bool {
	if len(as) != len(bs) {
		return false
	}

	for i := 0; i < len(as); i++ {
		if as[i] != bs[i] {
			return false
		}
	}

	return true
}

// EX:) []string{"hello", "world"} -> hello,world
func StrArrayToStr(targets []string) string {
	var val_str string
	for i, target := range targets {
		val_str += target
		if i != len(targets)-1 {
			val_str += ","
		}
	}
	return val_str
}
