package main

import (
	"fmt"
	"strings"
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
	arr := []string{"a", "b", "c", "d", "e", "f"}
	f1 := key([]int{2, 1})
	f2 := key([]int{2, 4})
	fmt.Println(f1(arr), f2(arr))
}
