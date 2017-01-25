package utils

// FilterStringArrayFunc function accept a strung value and return boolean to
// indicate whether the value is selected or not
type FilterStringArrayFunc func(string) bool

// FilterStringArray returns a subset of given string array that match condition
// defined in FilterStringArrayFunc
func FilterStringArray(arr []string, f FilterStringArrayFunc) []string {
	result := []string{}

	for _, value := range arr {
		if f(value) {
			result = append(result, value)
		}
	}

	return result
}

// SubtractStringArray returns new array of string which is the result of the first
// array subtracts the second
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
