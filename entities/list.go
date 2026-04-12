package entities

func List(array []string) string {
	if len(array) == 0 {
		return ""
	}
	if len(array) == 1 {
		return array[0]
	}
	result := array[0]
	for i := 1; i < len(array); i++ {
		result += ", " + array[i]
	}
	return result
}
