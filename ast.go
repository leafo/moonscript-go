package main

import "encoding/json"

type Lines struct {
	Nodes []any
}

type CommentNode struct {
	Text string
}

type PrimitiveNode struct {
	Value string
}

type AssignmentNode struct {
	Names []any
	Expr  any
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
