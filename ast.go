package main

import (
	"encoding/json"
)

type Node interface {
	ToLua(*LuaRenderState) (string, error)
}

type Lines struct {
	Lines []any
}

func (l *Lines) IsEmpty() bool {
	return len(l.Lines) == 0
}

type ExpressionList struct {
	Expressions []any
}

func (e *ExpressionList) IsEmpty() bool {
	return len(e.Expressions) == 0
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
	ExpressionList
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
	Ref string
}

type SelfRefNode struct {
	Ref string
}

type ParensNode struct {
	Expression any
}

// starts a chain with seriers of chain operations like indexing or calling
type ChainNode struct {
	Target any
	Ops    []any
}

// calls a function with arguments
type ChainCallNode struct {
	ExpressionList
}

// indexes field by name using dot notation
type ChainDotNode struct {
	Field string
}

type ChainMethodNode struct {
	Field string
}

// indexes field by expression with brackets
type ChainIndexNode struct {
	Index any
}

type IfStatementNode struct {
	Condition any
	Lines
	ElseIfs   []ElseIfStatementNode
	ElseLines Lines
}

type ElseIfStatementNode struct {
	Condition any
	Lines
}

type FunctionExpressionNode struct {
	IsMethod  bool
	Arguments []ArgumentTuple
	Lines
}

type ArgumentTuple struct {
	Name         any // string or ref to self
	DefaultValue any
}

type ReturnNode struct {
	ExpressionList
}

type FlowControlType int

const (
	Break FlowControlType = iota
	Continue
)

func (fct FlowControlType) String() string {
	switch fct {
	case Break:
		return "Break"
	case Continue:
		return "Continue"
	default:
		return "Unknown"
	}
}

func (fct FlowControlType) MarshalJSON() ([]byte, error) {
	return json.Marshal(fct.String())
}

type FlowControlNode struct {
	Type FlowControlType
}

type WhileStatementNode struct {
	Condition any
	Lines
}

type ForLoopNode struct {
	IndexName any
	Start     any
	End       any
	Step      any
	Lines
}

type ForEachLoopNode struct {
	Names []any
	ExpressionList
	Lines
}
