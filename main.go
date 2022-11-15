package main

import (
	Arguments "alperb/dirgod/arguments"
	Parser "alperb/dirgod/parser"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	arguments := Arguments.NewArguments(args)
	parser := Parser.NewDirParser(*arguments)
	parser.Parse()
	fmt.Println("Created!")
}
