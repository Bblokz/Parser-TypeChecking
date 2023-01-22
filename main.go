/*
 * Parser and Lexical Analyser LexicalAnalyser.go
 * Copyright (C) 2021-2023 Bas Blokzijl Leiden, The Netherlands.
 */

package main

import (
	"Parser-TypeChecking/Globals"
	"Parser-TypeChecking/LexicalAnalyser"
	"Parser-TypeChecking/Parser"
	"Parser-TypeChecking/TypeChecker"
	"bufio"
	"fmt"
	"os"
)

// check: Checks if the file can be opened if not,
// an exception is throw with the corresponding error.
// error: possible obtained error.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// get command arguments provided.
	commandArgs := os.Args

	if len(commandArgs) == 1 {
		fmt.Printf("Please provide a filename in the commandline")
		return
	}
	if len(commandArgs) > 2 {
		fmt.Printf("Too many arguments provided, please only provde the filename used as input!")
		return
	}
	// Open the file as provided by the commandline arguments
	data, err := os.Open(commandArgs[1])
	check(err)
	var variables = new(Globals.Vars)
	// Initiate a bufio scanner to analyse the data line by line.
	var scanner = bufio.NewScanner(data)
	for scanner.Scan() {
		variables.CurrentLine = []rune(scanner.Text())
		variables.CountBrackets = 0
		variables.ExpectVariable = false
		variables.DebugMode = false
		variables.Index = -1
		variables.Tree.IndexDoubleDot = -1
		variables.Token = LexicalAnalyser.LexicalAnalyser(variables)
		parser.Judgement(variables)
		TypeChecker.TypeChecker(variables)
		fmt.Println(variables.Tree.SubTreeToStandardOutput(0))
		variables.Tree.ClearTree()

	}

} // main
