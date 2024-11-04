package main

import "encoding/json"

type Lines struct {
	Lines []any
}

type TableNode struct {
	Tuples []TableTuple
}

type TableTuple struct {
	Key   any // can be nil for array positional value
	Value any
}

type CommentNode struct {
	Comment string
}

type PrimitiveNode struct {
	Primitive string
}

type AssignmentNode struct {
	Names []any
	Exprs []any
}

type ExpressionNode struct {
	Head any
	Rest []OperatorExpressionNode
}

type OperatorExpressionNode struct {
	Operator string
	Value    any
}

type StringNode struct {
	Delimiter   string
	StringParts []any // string literal or interpolation
}

type InterpolationNode struct {
	Expression any
}

// a parsed number
type NumberType int

const (
	Integer NumberType = iota
	Float
	Exp
)

func (nt NumberType) String() string {
	switch nt {
	case Integer:
		return "Integer"
	case Float:
		return "Float"
	case Exp:
		return "Exp"
	default:
		return "Unknown"
	}
}

func (nt NumberType) MarshalJSON() ([]byte, error) {
	return json.Marshal(nt.String())
}

type NumberNode struct {
	Type NumberType
	Text string
}

// reference to a variable
type RefNode struct {
	Name string
}
