package runtime

import (
	"fmt"
	"luag/parser"
)

type Runtime struct {
	Variables map[string]interface{}
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

func (r *Runtime) ExecuteStatement(stmt parser.Statement) {
	switch s := stmt.(type) {
	case *parser.LocalStatement:
		value := r.EvaluateExpression(s.Value)
		r.Variables[s.Name] = value
	case *parser.IfStatement:
		condition := r.EvaluateExpression(s.Condition)
		if conditionBool, ok := condition.(bool); ok && conditionBool {
			for _, thenStmt := range s.ThenBody {
				r.ExecuteStatement(thenStmt)
			}
		} else {
			for _, elseStmt := range s.ElseBody {
				r.ExecuteStatement(elseStmt)
			}
		}
	case *parser.FunctionCallStatement:
		r.ExecuteFunctionCall(s)
	default:
		fmt.Printf("Unknown statement type: %T\n", stmt)
	}
}

func (r *Runtime) EvaluateExpression(expr interface{}) interface{} {
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
	default:
		fmt.Printf("Unknown operator: %s\n", operator)
		return nil
	}
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
