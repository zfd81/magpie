package main

import (
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
)

func key(keys []int) func(strs []string) string {
	return func(strs []string) string {
		var builder strings.Builder
		for _, v := range keys {
			builder.WriteString(strs[v])
		}
		return builder.String()
	}
}

func main() {
	env := map[string]interface{}{
		"greet":   "Hello, %v!",
		"names":   []string{"world", "you"},
		"sprintf": fmt.Sprintf,
	}

	code := `sprintf(greet, names[0])`

	program, err := expr.Compile(code, expr.Env())
	if err != nil {
		panic(err)
	}
	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}

	fmt.Println(output)
}
