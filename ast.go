package main

import "encoding/json"

type Node interface{}

type Lines struct {
	Nodes []Node
}

type CommentNode struct {
	Text string
}

type StringNode struct {
	Delimiter   string
	StringParts interface{}
}

type InterpolationNode struct {
	Expression interface{}
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
