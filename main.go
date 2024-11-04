package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	jsonOutput := flag.Bool("json", false, "Generate JSON from the AST")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: moonscript-go [-json] <filename>")
		os.Exit(1)
	}

	filename := flag.Arg(0)
	result, err := ParseFile(filename)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		os.Exit(1)
	}

	if *jsonOutput {
		jsonResult, err := json.Marshal(result)
		if err != nil {
			fmt.Println("Error serializing result to JSON:", err)
			os.Exit(1)
		}
		fmt.Println(string(jsonResult))
	} else {
		luaCode, err := result.(Node).ToLua()
		if err != nil {
			fmt.Println("Error converting to Lua:", err)
			os.Exit(1)
		}
		fmt.Println(luaCode)
	}
}
