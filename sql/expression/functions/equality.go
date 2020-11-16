package functions

import "github.com/spf13/cast"

func Equality(a interface{}, b interface{}) bool {
	switch a.(type) {
	case string:
		val, ok := b.(string)
		if ok {
			return a == val
		}
		return false
	case int:
		val, ok := b.(int)
		if ok {
			return a == val
		}
		return false
	case int64:
		val, ok := b.(int64)
		if ok {
			return a == val
		}
		return false
	case bool:
		val, ok := b.(bool)
		if ok {
			return a == val
		}
		return false
	case float64:
		val, ok := b.(float64)
		if ok {
			return a == val
		}
		return false
	default:
		return cast.ToString(a) == cast.ToString(b)
	}
}
