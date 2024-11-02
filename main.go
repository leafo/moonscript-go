package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: moonscript-go <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	result, err := ParseFile(filename)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		os.Exit(1)
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error serializing result to JSON:", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonResult))
}
