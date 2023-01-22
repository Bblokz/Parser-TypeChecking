/*
 * Parser and Lexical Analyser Constant.go
 * Copyright (C) 2021-2023 Bas Blokzijl, Leiden, The Netherlands.
 */

package CharClass

// Character classes, used to determine the specific kind of character that is read by GetChar.
const (
	// LOWLETTER To distinguish lower case letters in variables.
	LOWLETTER = iota
	// UPLETTER To distinguish upper case letters in variables.
	UPLETTER
	// DIGIT To distinguish digits in variables.
	DIGIT
	// LBRACKET To distinguish the beginning of a nesting in brackets.
	LBRACKET
	// RBRACKET To distinguish the end of a nesting in brackets.
	RBRACKET
	// LAMBDA To distinguish a lambda expression.
	LAMBDA
	// TYPESYMBOL To distinguish '^' in a lambda expression.
	TYPESYMBOL
	// DoubleDot To distinguish between the <expr> and <type> in a <judgement>.
	DoubleDot
	// FUNCTION1 To distinguish the first part of a type function, the '-' in "->".
	FUNCTION1
	// FUNCTION2 To distinguish the second part of a type function, the '>' in "->".
	FUNCTION2
	// UNDEFINED For unicode characters that have no syntactical meaning to our Lexical Analyser.
	UNDEFINED
	// ENDOFLINE Indicates that the last read character is an end of line.
	ENDOFLINE
	// SPACE To distinguish spaces
	SPACE
)
