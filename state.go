package main

import (
	"fmt"
	"os"
)

const DEBUG_INDENT = false
const INVALID_INDENT = -1

func copyIndentStack(state map[string]any) []int {
	if state["current_indent"] == nil {
		return make([]int, 0)
	}
	existing := state["current_indent"].([]int)
	return append([]int{}, existing...)
}

func currentIndent(state map[string]any) int {
	if state["current_indent"] == nil {
		return INVALID_INDENT
	}

	stack := state["current_indent"].([]int)

	if len(stack) == 0 {
		return INVALID_INDENT
	}

	return stack[len(stack)-1]
}

func pushIndent(state map[string]any, indent int) {
	if DEBUG_INDENT {
		fmt.Fprint(os.Stderr, "\033[32m-> pushing ", state["current_indent"], "\033[0m")
	}
	stack := append(copyIndentStack(state), indent)
	state["current_indent"] = stack
	if DEBUG_INDENT {
		fmt.Fprintln(os.Stderr, "\033[32m ->", state["current_indent"], "\033[0m")
	}
}

func popIndent(state map[string]any) {
	if DEBUG_INDENT {
		fmt.Fprint(os.Stderr, "\033[31m<- popping ", state["current_indent"], "\033[0m")
	}
	stack := copyIndentStack(state)
	state["current_indent"] = stack[:len(stack)-1]
	if DEBUG_INDENT {
		fmt.Fprintln(os.Stderr, "\033[31m ->", state["current_indent"], "\033[0m")
	}
}
