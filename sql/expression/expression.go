package expression

import (
	"github.com/antonmedv/expr"
	"github.com/zfd81/magpie/sql/expression/functions"
)

func Eval(script string, env map[string]interface{}) (interface{}, error) {
	env["if"] = functions.If
	env["decode"] = functions.Decode
	env["eq"] = functions.Equality
	env["left"] = functions.Left
	env["right"] = functions.Right
	env["substr"] = functions.Substr
	program, err := expr.Compile(script, expr.Env(env))
	if err != nil {
		return "", err
	}
	return expr.Run(program, env)
}
