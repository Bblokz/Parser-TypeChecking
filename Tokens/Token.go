/*
 * Parser and Lexical Analyser Token.go
 * Copyright (C) 2021-2023 Bas Blokzijl, Leiden, The Netherlands.
 */

package Tokens

// Internal codes for the Lexemes, return values for LexicalAnalyser.
const (
	TokenVariable = iota
	TokenLambda
	TokenLeftBracket
	TokenRightBracket
	// TokenFunction also arrow function.
	TokenFunction
	TypeSymbol
	TokenUVar
	TokenDoubleDot
	// LexicalEndOfLine When the lexical analyser has reached the end of the line.
	LexicalEndOfLine
	SyntaxError
	// Application For the concatenation of two expressions.
	Application
)
