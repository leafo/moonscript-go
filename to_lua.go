package main

import (
	"fmt"
	"strings"
)

func (l Lines) ToLua() (string, error) {
	var out []string

	for _, line := range l.Lines {
		if line == nil {
			continue
		}

		lineNode, success := line.(Node)
		if !success {
			return "", fmt.Errorf("unknown node type: %T", line)
		}
		marshaled, err := lineNode.ToLua()
		if err != nil {
			return "", err
		}
		out = append(out, marshaled)
	}

	return strings.Join(out, "\n"), nil
}

func (n CommentNode) ToLua() (string, error) {
	return fmt.Sprintf("--%s", n.Comment), nil
}

func (n PrimitiveNode) ToLua() (string, error) {
	return n.Primitive, nil
}

func (n AssignmentNode) ToLua() (string, error) {
	var buf strings.Builder

	for i, name := range n.Names {
		if i > 0 {
			buf.WriteString(", ")
		}

		var marshaled string
		var err error

		switch n := name.(type) {
		case Node:
			marshaled, err = n.ToLua()
			if err != nil {
				return "", err
			}
		case string:
			marshaled = n
		default:
			return "", fmt.Errorf("AssignmentNode: unknown name type: %T", name)
		}

		buf.WriteString(marshaled)
	}

	buf.WriteString(" = ")

	for i, expr := range n.Exprs {
		if i > 0 {
			buf.WriteString(", ")
		}

		exprNode, success := expr.(Node)
		if !success {
			return "", fmt.Errorf("AssignmentNode: unknown value type: %T", expr)
		}
		marshaled, err := exprNode.ToLua()
		if err != nil {
			return "", err
		}
		buf.WriteString(marshaled)
	}

	return buf.String(), nil
}

func (n StringNode) ToLua() (string, error) {
	var out []string
	for _, part := range n.StringParts {
		val, isString := part.(string)
		if isString {
			out = append(out, val)
			continue
		}

		return "", fmt.Errorf("interpolation not implemented yet")
	}

	return fmt.Sprintf("%s%s%s", n.Delimiter, strings.Join(out, ""), n.Delimiter), nil
}

func (n NumberNode) ToLua() (string, error) {
	return n.Text, nil
}

func (n RefNode) ToLua() (string, error) {
	return n.Ref, nil
}

func (n ExpressionNode) ToLua() (string, error) {
	var buf strings.Builder

	headNode, success := n.Head.(Node)
	if !success {
		return "", fmt.Errorf("ExpressionNode: unknown head type: %T", n.Head)
	}

	head, err := headNode.ToLua()
	if err != nil {
		return "", err
	}

	buf.WriteString(head)

	for _, op := range n.Rest {
		val, success := op.Value.(Node)
		if !success {
			return "", fmt.Errorf("ExpressionNode: unknown value type: %T", op.Value)
		}

		valStr, err := val.ToLua()
		if err != nil {
			return "", err
		}

		buf.WriteString(" ")
		buf.WriteString(op.Operator)
		buf.WriteString(" ")
		buf.WriteString(valStr)
	}

	return buf.String(), nil
}

func (n ParensNode) ToLua() (string, error) {
	exprNode, success := n.Expression.(Node)
	if !success {
		return "", fmt.Errorf("ParensNode: unknown expression type: %T", n.Expression)
	}

	expr, err := exprNode.ToLua()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(%s)", expr), nil
}

func (n ChainNode) ToLua() (string, error) {
	var buf strings.Builder

	targetNode, ok := n.Target.(Node)
	if !ok {
		return "", fmt.Errorf("ChainNode: unknown target type: %T", n.Target)
	}

	target, err := targetNode.ToLua()
	if err != nil {
		return "", err
	}

	buf.WriteString(target)

	for _, op := range n.Ops {
		opNode, ok := op.(Node)
		if !ok {
			return "", fmt.Errorf("ChainNode: unknown op type: %T", op)
		}

		opStr, err := opNode.ToLua()
		if err != nil {
			return "", err
		}

		buf.WriteString(opStr)
	}

	return buf.String(), nil
}

func (n ChainCallNode) ToLua() (string, error) {
	var args []string
	for _, arg := range n.Arguments {
		argNode, ok := arg.(Node)
		if !ok {
			return "", fmt.Errorf("ChainCallNode: unknown argument type: %T", arg)
		}

		argStr, err := argNode.ToLua()
		if err != nil {
			return "", err
		}
		args = append(args, argStr)
	}

	return fmt.Sprintf("(%s)", strings.Join(args, ", ")), nil
}

func (n ChainDotNode) ToLua() (string, error) {
	return "." + n.Field, nil
}

func (n ChainIndexNode) ToLua() (string, error) {
	indexNode, ok := n.Index.(Node)
	if !ok {
		return "", fmt.Errorf("ChainIndexNode: unknown index type: %T", n.Index)
	}

	index, err := indexNode.ToLua()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("[%s]", index), nil
}
