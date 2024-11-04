package main

func copyIndentStack(state map[string]any) []int {
	if state["current_indent"] == nil {
		return make([]int, 0)
	}
	existing := state["current_indent"].([]int)
	return append([]int{}, existing...)
}

func currentIndent(state map[string]any) int {
	if state["current_indent"] == nil {
		return 0
	}

	stack := state["current_indent"].([]int)

	if len(stack) == 0 {
		return 0
	}

	return stack[len(stack)-1]
}

func pushIndent(state map[string]any, indent int) {
	// fmt.Print("-> pushing ", state["current_indent"])
	stack := append(copyIndentStack(state), indent)
	state["current_indent"] = stack
	// fmt.Println(" ->", state["current_indent"])
}

func popIndent(state map[string]any) {
	// fmt.Print("<- popping ", state["current_indent"])
	stack := copyIndentStack(state)
	state["current_indent"] = stack[:len(stack)-1]
	// fmt.Println(" ->", state["current_indent"])
}
