package functions

func Substr(str string, position int, length int) string {
	strRune := []rune(str)
	if position > len(strRune) {
		return ""
	} else if position > 0 {
		return Left(string(strRune[position-1:]), length)
	} else if (0 - position) <= len(strRune) {
		return Left(string(strRune[position+len(strRune):]), length)
	}
	return ""
}
