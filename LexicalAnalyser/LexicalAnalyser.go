/*
 * Parser and Lexical Analyser LexicalAnalyser.go
 * Copyright (C) 2021-2023 Bas Blokzijl Leiden, The Netherlands.
 */

package LexicalAnalyser

import (
	"Parser-TypeChecking/Constants"
	"Parser-TypeChecking/Globals"
	"Parser-TypeChecking/Tokens"
	"fmt"
	"os"
	"unicode"
)

// AddChar
// Adds the provided character to the lexeme.
func AddChar(context *Globals.Vars) {
	context.Lexeme = append(context.Lexeme, context.ReadChar)
}

// GetChar
// Gets the next character in our currently analysed line and determines its character class.
// context: Contains the whole expression.
func GetChar(context *Globals.Vars) {
	context.Index++
	if context.Index >= len(context.CurrentLine) {
		context.CharClass = CharClass.ENDOFLINE
	} else {
		context.ReadChar = context.CurrentLine[context.Index]
		if unicode.IsDigit(context.ReadChar) {
			context.CharClass = CharClass.DIGIT
		} else if unicode.IsLetter(context.ReadChar) && context.ReadChar != 'λ' {
			if unicode.IsLower(context.ReadChar) {
				context.CharClass = CharClass.LOWLETTER
			} else {
				context.CharClass = CharClass.UPLETTER
			}
		} else if (context.ReadChar == '\\') || (context.ReadChar == 'λ') {
			context.CharClass = CharClass.LAMBDA
		} else if context.ReadChar == '(' {
			context.CharClass = CharClass.LBRACKET
		} else if context.ReadChar == ')' {
			context.CharClass = CharClass.RBRACKET
		} else if context.ReadChar == '^' {
			context.CharClass = CharClass.TYPESYMBOL
		} else if context.ReadChar == '-' {
			context.CharClass = CharClass.FUNCTION1
		} else if context.ReadChar == '>' {
			context.CharClass = CharClass.FUNCTION2
		} else if context.ReadChar == ' ' {
			context.CharClass = CharClass.SPACE
		} else if context.ReadChar == ':' {
			context.CharClass = CharClass.DoubleDot
		} else {
			context.CharClass = CharClass.UNDEFINED
		}
	}
}

// LexicalAnalyser
// Analyses the variables.CurrentLine until the next token is found.
// returns the token that is found.
// now also works for typetree.
// context: Contains the whole expression.
func LexicalAnalyser(context *Globals.Vars) int {
	if context.DebugMode {
		fmt.Println("   Lexical called")
	}
	GetChar(context)
	// Save the previous token.
	context.PrevToken = context.Token
	if context.CharClass != CharClass.UNDEFINED {
		switch context.CharClass {
		case CharClass.LOWLETTER:
			// loop through until the variable is complete.
			for context.CharClass == CharClass.LOWLETTER ||
				context.CharClass == CharClass.UPLETTER || context.CharClass == CharClass.DIGIT {
				AddChar(context)
				GetChar(context)
			}
			// Do not lose the last read character, this will be handled on the next call.
			context.Index--
			return Tokens.TokenVariable
		case CharClass.UPLETTER:
			// loop through until the variable is complete
			for context.CharClass == CharClass.LOWLETTER ||
				context.CharClass == CharClass.UPLETTER || context.CharClass == CharClass.DIGIT {
				AddChar(context)
				GetChar(context)
			}
			// do not lose the last read character, this will be handled on the next call.
			context.Index--
			return Tokens.TokenUVar
		case CharClass.TYPESYMBOL: // ^
			// todo implement this right
			return Tokens.TypeSymbol
		case CharClass.FUNCTION1: // -
			// todo implement this right
			GetChar(context) //get the function 2
			return Tokens.TokenFunction
		case CharClass.FUNCTION2: // >
			fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: Function type had invalid declaration")
			os.Exit(1)
		case CharClass.DIGIT:
			fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: Cannot start a variable with digit")
			os.Exit(1)
		case CharClass.LBRACKET:
			return Tokens.TokenLeftBracket
		case CharClass.RBRACKET:
			return Tokens.TokenRightBracket
		case CharClass.LAMBDA:
			return Tokens.TokenLambda
		case CharClass.SPACE:
			return LexicalAnalyser(context)
		case CharClass.DoubleDot:
			return Tokens.TokenDoubleDot
		case CharClass.ENDOFLINE:
			return Tokens.LexicalEndOfLine
		}
	}
	return Tokens.SyntaxError
}

// GoBackToLastToken
// Reverts the last read token, and resets the LexicalAnalyser to read the same token on
// the next function call.
// context: Contains the whole expression.
func GoBackToLastToken(context *Globals.Vars) {
	context.Index--
	// Go back to previous token.
	for context.CurrentLine[context.Index] == ' ' {
		context.Index--
	}
	context.Lexeme = nil
}
