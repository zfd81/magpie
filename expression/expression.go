package expression

import (
	"github.com/antonmedv/expr"
	"github.com/zfd81/magpie/expression/functions"
)

func Eval(script string, env map[string]interface{}) (interface{}, error) {
	env["if"] = functions.If
	env["decode"] = functions.Decode
	program, err := expr.Compile(script, expr.Env(env))
	if err != nil {
		return "", err
	}
	return expr.Run(program, env)
}
