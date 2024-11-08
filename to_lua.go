package main

import (
	"fmt"
	"strings"
)

const INDENT_STR = "  "

type LuaRenderState struct {
	Indent int
}

func (state *LuaRenderState) WithIndent(str string) string {
	indentation := strings.Repeat(INDENT_STR, state.Indent)
	return indentation + str
}

func (l Lines) ToLua(state *LuaRenderState) (string, error) {
	var buf strings.Builder

	for _, line := range l.Lines {
		if line == nil {
			continue
		}

		lineNode, success := line.(Node)
		if !success {
			return "", fmt.Errorf("unknown node type: %T", line)
		}
		marshaled, err := lineNode.ToLua(state)
		if err != nil {
			return "", err
		}
		buf.WriteString(state.WithIndent(marshaled))
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

func (n CommentNode) ToLua(state *LuaRenderState) (string, error) {
	return fmt.Sprintf("--%s", n.Comment), nil
}

func (n PrimitiveNode) ToLua(state *LuaRenderState) (string, error) {
	return n.Primitive, nil
}

func (n ExpressionList) ToLua(state *LuaRenderState) (string, error) {
	var buf strings.Builder

	if n.IsEmpty() {
		return "", fmt.Errorf("ExpressionList: is empty, can't convert to lua")
	}

	for i, expr := range n.Expressions {
		if i > 0 {
			buf.WriteString(", ")
		}

		exprNode, success := expr.(Node)
		if !success {
			return "", fmt.Errorf("ExpressionList: unknown expression type: %T", expr)
		}
		marshaled, err := exprNode.ToLua(state)
		if err != nil {
			return "", err
		}
		buf.WriteString(marshaled)
	}

	return buf.String(), nil
}

func (n AssignmentNode) ToLua(state *LuaRenderState) (string, error) {
	var buf strings.Builder

	for i, name := range n.Names {
		if i > 0 {
			buf.WriteString(", ")
		}

		var marshaled string
		var err error

		switch n := name.(type) {
		case Node:
			marshaled, err = n.ToLua(state)
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

	marshaled, err := n.ExpressionList.ToLua(state)
	if err != nil {
		return "", err
	}
	buf.WriteString(marshaled)

	return buf.String(), nil
}

func (n StringNode) ToLua(state *LuaRenderState) (string, error) {
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

func (n NumberNode) ToLua(state *LuaRenderState) (string, error) {
	return n.Text, nil
}

func (n RefNode) ToLua(state *LuaRenderState) (string, error) {
	return n.Ref, nil
}

func (n SelfRefNode) ToLua(state *LuaRenderState) (string, error) {
	return "self." + n.Ref, nil
}

func (n ExpressionNode) ToLua(state *LuaRenderState) (string, error) {
	var buf strings.Builder

	headNode, success := n.Head.(Node)
	if !success {
		return "", fmt.Errorf("ExpressionNode: unknown head type: %T", n.Head)
	}

	head, err := headNode.ToLua(state)
	if err != nil {
		return "", err
	}

	buf.WriteString(head)

	for _, op := range n.Rest {
		val, success := op.Value.(Node)
		if !success {
			return "", fmt.Errorf("ExpressionNode: unknown value type: %T", op.Value)
		}

		valStr, err := val.ToLua(state)
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

func (n ParensNode) ToLua(state *LuaRenderState) (string, error) {
	exprNode, success := n.Expression.(Node)
	if !success {
		return "", fmt.Errorf("ParensNode: unknown expression type: %T", n.Expression)
	}

	expr, err := exprNode.ToLua(state)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(%s)", expr), nil
}

func (n ChainNode) ToLua(state *LuaRenderState) (string, error) {
	var buf strings.Builder

	targetNode, ok := n.Target.(Node)
	if !ok {
		return "", fmt.Errorf("ChainNode: unknown target type: %T", n.Target)
	}

	target, err := targetNode.ToLua(state)
	if err != nil {
		return "", err
	}

	buf.WriteString(target)

	for _, op := range n.Ops {
		opNode, ok := op.(Node)
		if !ok {
			return "", fmt.Errorf("ChainNode: unknown op type: %T", op)
		}

		opStr, err := opNode.ToLua(state)
		if err != nil {
			return "", err
		}

		buf.WriteString(opStr)
	}

	return buf.String(), nil
}

func (n ChainCallNode) ToLua(state *LuaRenderState) (string, error) {
	args := ""

	if !n.IsEmpty() {
		var err error
		args, err = n.ExpressionList.ToLua(state)
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("(%s)", args), nil
}

func (n ChainDotNode) ToLua(state *LuaRenderState) (string, error) {
	return "." + n.Field, nil
}

func (n ChainMethodNode) ToLua(state *LuaRenderState) (string, error) {
	return ":" + n.Field, nil
}
func (n ChainIndexNode) ToLua(state *LuaRenderState) (string, error) {
	indexNode, ok := n.Index.(Node)
	if !ok {
		return "", fmt.Errorf("ChainIndexNode: unknown index type: %T", n.Index)
	}

	index, err := indexNode.ToLua(state)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("[%s]", index), nil
}

func (n TableNode) ToLua(state *LuaRenderState) (string, error) {
	var buf strings.Builder
	buf.WriteString("{")

	for i, tuple := range n.Tuples {
		if i > 0 {
			buf.WriteString(", ")
		}

		if tuple.Key == nil {
			// Array-style entry
			valueNode, ok := tuple.Value.(Node)
			if !ok {
				return "", fmt.Errorf("TableNode: unknown value type: %T", tuple.Value)
			}
			value, err := valueNode.ToLua(state)
			if err != nil {
				return "", err
			}
			buf.WriteString(value)
		} else {
			// Key-value pair
			switch tuple.Key.(type) {
			case string:
				key := tuple.Key.(string)
				buf.WriteString(key)
				buf.WriteString(" = ")
			case Node:
				keyNode := tuple.Key.(Node)

				keyLua, err := keyNode.ToLua(state)
				if err != nil {
					return "", err
				}

				buf.WriteString("[")
				buf.WriteString(keyLua)
				buf.WriteString("] = ")
			default:
				return "", fmt.Errorf("TableNode: unknown key type: %T", tuple.Key)
			}

			valueNode, ok := tuple.Value.(Node)
			if !ok {
				return "", fmt.Errorf("TableNode: unknown value type: %T", tuple.Value)
			}

			value, err := valueNode.ToLua(state)
			if err != nil {
				return "", err
			}

			buf.WriteString(value)
		}
	}

	buf.WriteString("}")
	return buf.String(), nil
}

func (n IfStatementNode) ToLua(state *LuaRenderState) (string, error) {
	conditionStr, err := n.Condition.(Node).ToLua(state)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	buf.WriteString("if ")
	buf.WriteString(conditionStr)
	buf.WriteString(" then\n")

	state.Indent += 1

	linesStr, err := n.Lines.ToLua(state)
	if err != nil {
		return "", err
	}

	buf.WriteString(linesStr)

	state.Indent -= 1

	for _, elseif := range n.ElseIfs {
		elseifConditionStr, err := elseif.Condition.(Node).ToLua(state)
		if err != nil {
			return "", err
		}

		buf.WriteString(state.WithIndent("elseif "))
		buf.WriteString(elseifConditionStr)
		buf.WriteString(" then\n")

		state.Indent += 1

		linesStr, err := elseif.Lines.ToLua(state)
		if err != nil {
			return "", err
		}

		buf.WriteString(linesStr)

		state.Indent -= 1
	}

	if !n.ElseLines.IsEmpty() {
		buf.WriteString(state.WithIndent("else\n"))

		state.Indent += 1

		linesStr, err := n.ElseLines.ToLua(state)
		if err != nil {
			return "", err
		}

		buf.WriteString(linesStr)

		state.Indent -= 1
	}

	buf.WriteString(state.WithIndent("end"))
	return buf.String(), nil
}

func (n FunctionExpressionNode) ToLua(state *LuaRenderState) (string, error) {
	var buf strings.Builder

	buf.WriteString("function(")
	k := 0

	if n.IsMethod {
		buf.WriteString("self")
		k++
	}

	for _, arg := range n.Arguments {
		if k > 0 {
			buf.WriteString(", ")
		}
		argName, ok := arg.Name.(string)
		if !ok {
			return "", fmt.Errorf("FunctionExpressionNode: argument name is not a string: %T", arg.Name)
		}
		buf.WriteString(argName)
		k++
	}

	buf.WriteString(")\n")

	state.Indent += 1

	linesStr, err := n.Lines.ToLua(state)
	if err != nil {
		return "", err
	}

	buf.WriteString(linesStr)

	state.Indent -= 1

	buf.WriteString(state.WithIndent("end"))
	return buf.String(), nil
}

func (n ReturnNode) ToLua(state *LuaRenderState) (string, error) {
	var buf strings.Builder

	buf.WriteString("return")

	for i, expr := range n.Expressions {
		if i == 0 {
			buf.WriteString(" ")
		} else if i > 0 {
			buf.WriteString(", ")
		}
		exprNode, ok := expr.(Node)
		if !ok {
			return "", fmt.Errorf("ReturnNode: unknown expression type: %T", expr)
		}
		exprStr, err := exprNode.ToLua(state)
		if err != nil {
			return "", err
		}
		buf.WriteString(exprStr)
	}

	return buf.String(), nil
}

func (n FlowControlNode) ToLua(state *LuaRenderState) (string, error) {
	var buf strings.Builder

	switch n.Type {
	case Break:
		buf.WriteString("break")
	case Continue:
		buf.WriteString("continue")
	default:
		return "", fmt.Errorf("FlowControlNode: unknown flow control type: %v", n.Type)
	}

	return buf.String(), nil
}

func (n WhileStatementNode) ToLua(state *LuaRenderState) (string, error) {
	var buf strings.Builder

	buf.WriteString(state.WithIndent("while "))
	conditionStr, err := n.Condition.(Node).ToLua(state)
	if err != nil {
		return "", err
	}
	buf.WriteString(conditionStr)
	buf.WriteString(" do\n")

	state.Indent += 1

	linesStr, err := n.Lines.ToLua(state)
	if err != nil {
		return "", err
	}

	buf.WriteString(linesStr)

	state.Indent -= 1

	buf.WriteString(state.WithIndent("end"))
	return buf.String(), nil
}

func (n ForLoopNode) ToLua(state *LuaRenderState) (string, error) {
	var buf strings.Builder

	buf.WriteString(state.WithIndent("for "))
	indexName, ok := n.IndexName.(string)
	if !ok {
		return "", fmt.Errorf("ForLoopNode: index name is not a string: %T", n.IndexName)
	}
	buf.WriteString(indexName)
	buf.WriteString(" = ")

	startStr, err := n.Start.(Node).ToLua(state)
	if err != nil {
		return "", err
	}
	buf.WriteString(startStr)
	buf.WriteString(", ")

	endStr, err := n.End.(Node).ToLua(state)
	if err != nil {
		return "", err
	}
	buf.WriteString(endStr)

	if n.Step != nil {
		buf.WriteString(", ")
		stepStr, err := n.Step.(Node).ToLua(state)
		if err != nil {
			return "", err
		}
		buf.WriteString(stepStr)
	}

	buf.WriteString(" do\n")

	state.Indent += 1

	linesStr, err := n.Lines.ToLua(state)
	if err != nil {
		return "", err
	}

	buf.WriteString(linesStr)

	state.Indent -= 1

	buf.WriteString(state.WithIndent("end"))
	return buf.String(), nil
}

func (n ForEachLoopNode) ToLua(state *LuaRenderState) (string, error) {
	var buf strings.Builder

	buf.WriteString(state.WithIndent("for "))
	for i, name := range n.Names {
		if i > 0 {
			buf.WriteString(", ")
		}
		nameStr, ok := name.(string)
		if !ok {
			return "", fmt.Errorf("ForEachLoopNode: name is not a string: %T", name)
		}
		buf.WriteString(nameStr)
	}
	buf.WriteString(" in ")

	exprStr, err := n.ExpressionList.ToLua(state)
	if err != nil {
		return "", err
	}
	buf.WriteString(exprStr)
	buf.WriteString(" do\n")

	state.Indent += 1

	linesStr, err := n.Lines.ToLua(state)
	if err != nil {
		return "", err
	}

	buf.WriteString(linesStr)

	state.Indent -= 1

	buf.WriteString(state.WithIndent("end"))
	return buf.String(), nil
}
