package entities

func ReadStr(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func PtrStr(s string) *string {
	return &s
}
