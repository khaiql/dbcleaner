package utils

type FilterStringArrayFunc func(string) bool

func FilterStringArray(arr []string, f FilterStringArrayFunc) []string {
	result := []string{}

	for _, value := range arr {
		if f(value) {
			result = append(result, value)
		}
	}

	return result
}
