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

func SubtractStringArray(firstArr, secondArr []string) []string {
	secondArrHash := map[string]bool{}
	for _, v := range secondArr {
		secondArrHash[v] = true
	}

	filterFunc := func(value string) bool {
		return !secondArrHash[value]
	}

	return FilterStringArray(firstArr, filterFunc)
}
