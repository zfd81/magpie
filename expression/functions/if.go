package functions

func If(c bool, t interface{}, f interface{}) interface{} {
	if c {
		return t
	} else {
		return f
	}
}
