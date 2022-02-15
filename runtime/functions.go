package runtime

import (
	"fmt"
	"time"
)

type FunctionHandler func(e *Evaluator, args []Object) (Object, error)

type Function struct {
	name    string
	arity   int
	Handler FunctionHandler
}

func (f Function) Type() ObjectType                                 { return FUNC_OBJ }
func (f Function) String() string                                   { return f.name }
func (f Function) Arity() int                                       { return f.arity }
func (f Function) Call(e *Evaluator, args []Object) (Object, error) { return f.Handler(e, args) }

var NativeFunctions = []Function{
	{
		name:    "sleep",
		arity:   1,
		Handler: sleepHandler,
	},
	{
		name:    "time",
		arity:   0,
		Handler: timeHandler,
	},
}

func sleepHandler(e *Evaluator, args []Object) (Object, error) {
	arg := args[0]

	switch arg := arg.(type) {
	case Float64:
		time.Sleep(time.Duration(arg.Value) * time.Second)
		return NIL, nil
	default:
		return NIL, fmt.Errorf("sleep expects an argument of type float64")
	}
}

func timeHandler(e *Evaluator, args []Object) (Object, error) {
	ms := time.Now().UnixNano() / int64(time.Millisecond)
	return NewFloat64(float64(ms)), nil
}