package runtime

import (
	"fmt"
	"github.com/shreerangdixit/lox/ast"
	"math"
	"time"
)

// ------------------------------------
// User function
// ------------------------------------

type UserFunction struct {
	node ast.FunctionNode
}

func NewUserFunction(node ast.FunctionNode) *UserFunction {
	return &UserFunction{
		node: node,
	}
}

func (f *UserFunction) Type() ObjectType { return TypeFunc }
func (f *UserFunction) String() string   { return f.node.Identifier.Token.Literal }
func (f *UserFunction) Arity() int       { return len(f.node.Parameters) }
func (f *UserFunction) Variadic() bool   { return false }
func (f *UserFunction) Call(e *Evaluator, args []Object) (Object, error) {
	// Reset environment at the end of function call
	prev := e.env
	defer func() {
		e.env = prev
	}()

	// New environment at the beginning of function call
	e.env = NewEnvWithEnclosing(e.env)
	for i := range args {
		// Bind function arguments to values
		err := e.env.Declare(f.node.Parameters[i].Token.Literal, args[i])
		if err != nil {
			return NIL, err
		}
	}
	return e.eval(f.node.Body)
}

// ------------------------------------
// Native (in-built) functions
// ------------------------------------

type NativeFunctionHandler func(e *Evaluator, args []Object) (Object, error)

type NativeFunction struct {
	name     string
	arity    int
	variadic bool
	handler  NativeFunctionHandler
}

func NewNativeFunction(name string, arity int, variadic bool, handler NativeFunctionHandler) *NativeFunction {
	return &NativeFunction{
		name:     name,
		arity:    arity,
		variadic: variadic,
		handler:  handler,
	}
}

func (f *NativeFunction) Type() ObjectType                                 { return TypeFunc }
func (f *NativeFunction) String() string                                   { return f.name }
func (f *NativeFunction) Arity() int                                       { return f.arity }
func (f *NativeFunction) Variadic() bool                                   { return f.variadic }
func (f *NativeFunction) Call(e *Evaluator, args []Object) (Object, error) { return f.handler(e, args) }

var NativeFunctions = []*NativeFunction{
	NewNativeFunction("sleep", 1, false, sleepHandler),
	NewNativeFunction("time", 0, false, timeHandler),
	NewNativeFunction("abs", 1, false, absHandler),
	NewNativeFunction("max", 2, false, maxHandler),
	NewNativeFunction("min", 2, false, minHandler),
	NewNativeFunction("avg", 1, false, avgHandler),
	NewNativeFunction("sqrt", 1, false, sqrtHandler),
	NewNativeFunction("type", 1, false, typeHandler),
	NewNativeFunction("len", 1, false, lenHandler),
	NewNativeFunction("print", 0, true, printHandler),
	NewNativeFunction("println", 0, true, printlnHandler),
}

func sleepHandler(e *Evaluator, args []Object) (Object, error) {
	arg, ok := args[0].(Number)
	if !ok {
		return NIL, fmt.Errorf("sleep() expects a number")
	}

	time.Sleep(time.Duration(arg.Value) * time.Millisecond)
	return NIL, nil
}

func timeHandler(e *Evaluator, args []Object) (Object, error) {
	ms := time.Now().UnixNano() / int64(time.Millisecond)
	return NewNumber(float64(ms)), nil
}

func absHandler(e *Evaluator, args []Object) (Object, error) {
	arg, ok := args[0].(Number)
	if !ok {
		return NIL, fmt.Errorf("abs() expects a number")
	}
	return NewNumber(math.Abs(arg.Value)), nil
}

func maxHandler(e *Evaluator, args []Object) (Object, error) {
	arg1, ok := args[0].(Number)
	if !ok {
		return NIL, fmt.Errorf("max() expects a number")
	}

	arg2, ok := args[1].(Number)
	if !ok {
		return NIL, fmt.Errorf("max() expects a number")
	}

	return NewNumber(math.Max(arg1.Value, arg2.Value)), nil
}

func minHandler(e *Evaluator, args []Object) (Object, error) {
	arg1, ok := args[0].(Number)
	if !ok {
		return NIL, fmt.Errorf("min() expects a number")
	}

	arg2, ok := args[1].(Number)
	if !ok {
		return NIL, fmt.Errorf("min() expects a number")
	}

	return NewNumber(math.Min(arg1.Value, arg2.Value)), nil
}

func avgHandler(e *Evaluator, args []Object) (Object, error) {
	seq, ok := args[0].(Sequence)
	if !ok {
		return NIL, fmt.Errorf("avg() expects a sequence")
	}

	sum := NewNumber(0)
	for _, arg := range seq.Elements() {
		if num, ok := arg.(Number); ok {
			s, _ := sum.Add(num)
			sum = s.(Number)
		} else {
			return NIL, fmt.Errorf("avg() expects numbers")
		}
	}

	return sum.Divide(seq.Size())
}

func sqrtHandler(e *Evaluator, args []Object) (Object, error) {
	if num, ok := args[0].(Number); ok {
		return NewNumber(math.Sqrt(num.Value)), nil
	}
	return NIL, fmt.Errorf("sqrt() expects a number")
}

func typeHandler(e *Evaluator, args []Object) (Object, error) {
	arg := args[0]
	return NewType(arg.Type()), nil
}

func lenHandler(e *Evaluator, args []Object) (Object, error) {
	arg, ok := args[0].(Sequence)
	if !ok {
		return NIL, fmt.Errorf("len() expects a sequence")
	}
	return arg.Size(), nil
}

func printHandler(e *Evaluator, args []Object) (Object, error) {
	for _, obj := range args {
		fmt.Printf("%v", obj)
		fmt.Printf(" ")
	}
	return NIL, nil
}

func printlnHandler(e *Evaluator, args []Object) (Object, error) {
	_, _ = printHandler(e, args)
	fmt.Println()
	return NIL, nil
}
