package runtime

import (
	"fmt"
	"luag/parser"
)

const (
	START_INDEX = float64(1)
) // Lua tables are 1-indexed
type Runtime struct {
	Variables map[string]interface{}
}
type LuaTable struct {
	Fields map[interface{}]interface{}
}
type LuaFunction struct {
	Parameters []string
	Body       []parser.Statement
}
type executionresult struct {
	returned bool
	value    interface{}
}

func NewLuaTable() *LuaTable {
	return &LuaTable{
		Fields: make(map[interface{}]interface{}),
	}
}
func (t *LuaTable) String() string {
	return fmt.Sprintf("LuaTable{%v}", t.Fields)
}
func NewRuntime() *Runtime {
	return &Runtime{
		Variables: make(map[string]interface{}),
	}
}

func (r *Runtime) ExecuteChunk(chunk *parser.Chunk) {
	for _, stmt := range chunk.Statements {
		r.ExecuteStatement(stmt)
	}
}

// executechunk+catch print statements to a buffer for testing purposes
func (r *Runtime) ExecuteChunkWithOutput(chunk *parser.Chunk) string {
	var output string
	for _, stmt := range chunk.Statements {
		switch s := stmt.(type) {
		case *parser.FunctionCallStatement:
			if s.Name == "print" {
				for _, arg := range s.Args {
					value := r.EvaluateExpression(arg)
					output += fmt.Sprintf("%v\n", value)
				}
			} else {
				r.ExecuteFunctionCall(s)
			}
		default:
			r.ExecuteStatement(stmt)
		}
	}
	return output
}

func (r *Runtime) ExecuteStatement(stmt parser.Statement) executionresult {
	switch s := stmt.(type) {
	case *parser.LocalStatement:
		value := r.EvaluateExpression(s.Value)
		r.Variables[s.Name] = value
		return executionresult{returned: false, value: nil}
	case *parser.ReturnStatement:
		return executionresult{returned: true, value: r.EvaluateExpression(s.ReturnValue)}
	case *parser.IfStatement:
		condition := r.EvaluateExpression(s.Condition)
		body := s.ThenBody
		if cond, ok := condition.(bool); ok && cond {
			for _, stmt := range body {
				result := r.ExecuteStatement(stmt)
				if result.returned {
					return result
				}
			}
		} else {
			for _, stmt := range s.ElseBody {
				result := r.ExecuteStatement(stmt)
				if result.returned {
					return result
				}
			}
		}
		return executionresult{returned: false, value: nil}
	case *parser.FunctionCallStatement:
		r.ExecuteFunctionCall(s)
		return executionresult{returned: false, value: nil}
	case *parser.FunctionDefStatement:
		function := &LuaFunction{
			Parameters: s.Parameters,
			Body:       s.Body,
		}
		r.Variables[s.Name] = function
		return executionresult{returned: false, value: nil}
	default:
		fmt.Printf("Unknown statement type: %T\n", stmt)
		return executionresult{returned: false, value: nil}
	}
}
func (r *Runtime) callLuaFunction(function *LuaFunction, args []interface{}) interface{} {
	globalScope := r.Variables
	localScope := make(map[string]interface{})

	for name, value := range globalScope {
		localScope[name] = value
	}
	for i, param := range function.Parameters {
		if i < len(args) {
			localScope[param] = args[i]
		} else {
			localScope[param] = nil
		}
	}
	r.Variables = localScope
	defer func() {
		r.Variables = globalScope
	}()
	for _, stmt := range function.Body {
		result := r.ExecuteStatement(stmt)
		if result.returned {
			return result.value
		}
	}
	return nil
}
func (r *Runtime) EvaluateExpression(expr parser.Expression) any {
	switch e := expr.(type) {
	case *parser.NumberLiteral:
		return e.Value
	case *parser.StringLiteral:
		return e.Value
	case *parser.Identifier:
		if val, exists := r.Variables[e.Value]; exists {
			return val
		}
		fmt.Printf("Undefined variable: %s\n", e.Value)
		return nil
	case *parser.BinaryExpression:
		left := r.EvaluateExpression(e.Left)
		right := r.EvaluateExpression(e.Right)
		return r.EvaluateBinaryExpression(left, e.Operator, right)
	case *parser.TableConstructorExpression:
		return r.EvaluateTableConstructor(e)
	case *parser.IndexedAccessExpression:
		table := r.EvaluateExpression(e.Table)
		key := r.EvaluateExpression(e.Key)
		if tbl, ok := table.(*LuaTable); ok {
			return tbl.Fields[key]
		}
		fmt.Printf("TypeError: expected LuaTable for indexed access, got %T\n", table)
		return nil
	case *parser.FunctionCallExpression:
		return r.EvaluateFunctionCall(e)
	case *parser.FunctionExpression:
		return &LuaFunction{
			Parameters: e.Parameters,
			Body:       e.Body,
		}
	default:
		fmt.Printf("Unknown expression type: %T\n", expr)
		return nil
	}
}

func (r *Runtime) EvaluateBinaryExpression(left interface{}, operator string, right interface{}) interface{} {
	switch operator {
	case ">", "<", ">=", "<=":
		switch l := left.(type) {
		case float64:
			if r, ok := right.(float64); ok {
				switch operator {
				case ">":
					return l > r
				case "<":
					return l < r
				case ">=":
					return l >= r
				case "<=":
					return l <= r
				}
			}
		case string:
			if r, ok := right.(string); ok {
				switch operator {
				case ">":
					return l > r
				case "<":
					return l < r
				}
			}
		}
		fmt.Printf("TypeError: invalid operands for %s: %T and %T\n", operator, left, right)
		return nil
	case "==":
		return left == right
	case "~=":
		return left != right
	case "+":
		switch l := left.(type) {
		case float64:
			if r, ok := right.(float64); ok {
				return l + r
			}
		case string:
			if r, ok := right.(string); ok {
				return l + r
			}
		}
		fmt.Printf("TypeError: invalid operands for %s: %T and %T\n", operator, left, right)
		return nil
	case "-":
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l - r
			}
		}
		fmt.Printf("TypeError: invalid operands for %s: %T and %T\n", operator, left, right)
		return nil
	case "*":
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l * r
			}
		}
		fmt.Printf("TypeError: invalid operands for %s: %T and %T\n", operator, left, right)
		return nil
	case "/":
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				if r == 0 {
					fmt.Println("Error: Division by zero")
					return nil
				}
				return l / r
			}
		}
		fmt.Printf("TypeError: invalid operands for %s: %T and %T\n", operator, left, right)
		return nil
	default:
		fmt.Printf("Unknown operator: %s\n", operator)
		return nil
	}
}
func (r *Runtime) EvaluateFunctionCall(
	call *parser.FunctionCallExpression,
) interface{} {
	functionValue := r.EvaluateExpression(call.Function)

	function, ok := functionValue.(*LuaFunction)
	if !ok {
		fmt.Printf(
			"TypeError: attempted to call value of type %T\n",
			functionValue,
		)
		return nil
	}

	args := make([]interface{}, len(call.Args))
	for i, arg := range call.Args {
		args[i] = r.EvaluateExpression(arg)
	}

	return r.callLuaFunction(function, args)
}

func (r *Runtime) EvaluateTableConstructor(expr *parser.TableConstructorExpression) *LuaTable {
	table := NewLuaTable()
	nextIndex := START_INDEX
	for _, field := range expr.Fields {
		value := r.EvaluateExpression(field.Value)
		if field.Key == nil { // Implicit integer key
			table.Fields[nextIndex] = value
			nextIndex++
			continue
		}
		key := r.EvaluateExpression(field.Key)
		if key == nil {
			fmt.Println("Error: Table field key evaluated to nil")
			continue
		}
		table.Fields[key] = value
	}
	return table
}

func (r *Runtime) ExecuteFunctionCall(call *parser.FunctionCallStatement) {
	switch call.Name {
	case "print":
		for _, arg := range call.Args {
			value := r.EvaluateExpression(arg)
			fmt.Println(value)
		}
	default:
		fmt.Printf("Unknown function: %s\n", call.Name)
	}
}

func (r *Runtime) PrintVariables() {
	fmt.Println("Current Variables:")
	for name, value := range r.Variables {
		fmt.Printf("%s = %v\n", name, value)
	}
}
