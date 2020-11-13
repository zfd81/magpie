package functions

func Left(str string, length int) string {
	if str == "" || length < 0 {
		return ""
	}
	strRune := []rune(str)
	if length < len(strRune) {
		return string(strRune[:length])
	} else {
		return str
	}
}
