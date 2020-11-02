package expression

import "github.com/antonmedv/expr"

func Eval(script string, env map[string]interface{}) (interface{}, error) {
	program, err := expr.Compile(script, expr.Env(env))
	if err != nil {
		return "", err
	}
	return expr.Run(program, env)
}
