package functions

func Decode(args ...interface{}) interface{} {
	size := len(args)
	if size < 3 {
		return args[0]
	}
	value := args[0]
	var conv func(interface{}) (interface{}, bool)
	switch value.(type) {
	case string:
		conv = func(i interface{}) (interface{}, bool) {
			val, ok := i.(string)
			return val, ok
		}
	case int:
		conv = func(i interface{}) (interface{}, bool) {
			val, ok := i.(int)
			return val, ok
		}
	case int64:
		conv = func(i interface{}) (interface{}, bool) {
			val, ok := i.(int64)
			return val, ok
		}
	case bool:
		conv = func(i interface{}) (interface{}, bool) {
			val, ok := i.(bool)
			return val, ok
		}
	case float64:
		conv = func(i interface{}) (interface{}, bool) {
			val, ok := i.(float64)
			return val, ok
		}
	default:
		conv = func(i interface{}) (interface{}, bool) {
			val, ok := i.(string)
			return val, ok
		}
	}
	for i := 1; i < size-1; i = i + 2 {
		if val, ok := conv(args[i]); ok && val == value {
			return args[i+1]
		}
	}
	if (size-1)%2 == 1 {
		return args[size-1]
	}
	return value
}
