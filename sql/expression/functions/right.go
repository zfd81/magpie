package functions

func Right(str string, length int) string {
	if str == "" || length < 0 {
		return ""
	}
	strRune := []rune(str)
	strLen := len(strRune)
	if length < strLen {
		return string(strRune[strLen-length:])
	} else {
		return str
	}
}
