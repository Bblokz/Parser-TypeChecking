/*
 * Parser and Lexical Analyser Parser.go
 * Copyright (C) 2021-2023 Bas Blokzijl Leiden, The Netherlands.
 */

package parser

import (
	"Parser-TypeChecking/Globals"
	"Parser-TypeChecking/LexicalAnalyser"
	"Parser-TypeChecking/Tokens"
	"fmt"
	"os"
)

// Judgement Initiates Recursive Descent Parsing.
// Expects a non-empty expression followed by a ':' and a TypeExpression.
// context: contains the expression.
func Judgement(context *Globals.Vars) {
	Expr(context)
	JudgementFunction(context)
	TypeExpr(context)
} // Judgement

// Expr Part of parsing <expr> in Recursive Descent Parsing.
// Expects a non-empty <expr> using the LExpr function.
// after that allows for an end-of-line or ':' using MsExpr.
//  context: contains the expression.
func Expr(context *Globals.Vars) {
	if context.DebugMode {
		fmt.Println("   Expr called")
	}
	// We expect a <lvar>, left-bracket or lambda-expression.
	LExpr(context)
	// Allows for end-of-line or ':'.
	MsExpr(context)
}

// JudgementFunction
// "Peeks" in the string to be lexically analysed and checks whether the next token
// is a JudgementFunction. If so, the token is parsed and Recursive Descent continues accordingly.
// Else, there is a syntax error.
// context: contains the expression.
func JudgementFunction(context *Globals.Vars) bool {
    if context.DebugMode {
		fmt.Println(" JudgementFunction called")
	}
	if context.Tree.IndexDoubleDot != -1 {
		fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: double Judgement Function (:).")
		os.Exit(1)
	} else if context.Token == Tokens.TokenDoubleDot {
		context.Tree.AddToken(Tokens.TokenDoubleDot, "Judge", 0)
            
        context.Tree.IndexDoubleDot = len(context.Tree.Nodes)

		context.Token = LexicalAnalyser.LexicalAnalyser(context)
		return true
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: Expected Judgement Function (:).")
		os.Exit(1)
	}
	return false
}

// MsExpr
// Expr-prime, is allowed to be empty; meaning an end-of-line or ':' token is allowed
// If end-of-line or ':' is found, the parsing will stop or continue with <type> respectively.
// Else, the Recursive Descent is continued with an <expr>.
// context: Contains the whole expression.
func MsExpr(context *Globals.Vars) {
	if context.DebugMode {
		fmt.Println("   MsExpr called")
	}
	if context.Token == Tokens.LexicalEndOfLine || context.Token == Tokens.TokenDoubleDot {
		return
	} else {
		if context.DebugMode {
			fmt.Println("Line not finished, found: " + string(context.ReadChar))
		}
	}
	// No end of line; expected non-empty expression!
	LExpr(context)
	// Possible empty expression allowed.
	MsExpr(context)
}

// LExpr
// Determines the next step in recursive descent.
// Contains all possible continuations for <expr> namely, a <lvar>, ( <expr>) or lambda <lvar>'^'<type> <expr>.
// The <expr> <expr> continuations is handled using MsExpr
// context: Contains the whole expression.
func LExpr(context *Globals.Vars) {
	if context.DebugMode {
		fmt.Println("   LExpr called")
	}
	if context.ExpectVariable && context.Token != Tokens.TokenVariable {
		fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: Expected variable")
		os.Exit(1)
	}
	if context.ExpectUp && context.Token != Tokens.TypeSymbol {
		fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: Expected ^ but got: "+string(context.CurrentLine[context.Index]))
		os.Exit(1)
	}
	switch context.Token {
	case Tokens.TokenLambda:
		context.Tree.AddToken(Tokens.TokenLambda, string('λ'), 0)
		context.Token = LexicalAnalyser.LexicalAnalyser(context)
		// Force the next LExpr to contain a variable otherwise we have a syntax error.
		context.ExpectVariable = true
		LExpr(context)
		// Force the next LExpr to contain a '^' otherwise we have a syntax error.
		context.ExpectUp = true
		LExpr(context)
		// TypeExpression required.
		TypeExpr(context)
		// Non-empty expression required.
		LExpr(context)
		break
	case Tokens.TokenVariable:
		VarExpr(context)
		break
	case Tokens.TokenLeftBracket:
		// Increment shows that there is an extra open bracket.
		context.CountBrackets++
		context.Token = LexicalAnalyser.LexicalAnalyser(context)
		// Non-empty expression required.
		LExpr(context)
		// Expected non-empty expression.
		LExpr(context)
		break
	case Tokens.TokenRightBracket:
		// Decrement shows that there is an extra closed bracket.
		context.CountBrackets--
		if context.CountBrackets < 0 {
			fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: No brackets to close.")
			os.Exit(1)
		}
		context.Token = LexicalAnalyser.LexicalAnalyser(context)
		break
	case Tokens.TypeSymbol:
		context.ExpectUp = false
		context.Token = LexicalAnalyser.LexicalAnalyser(context)
		break
	case Tokens.TokenUVar:
		if context.PrevToken == Tokens.TokenUVar {
			fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: Expected TermFunction (->).")
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: Cannot Parse Uvar in Expression.")
			os.Exit(1)
		}
	case Tokens.LexicalEndOfLine:
		if context.CountBrackets > 0 {
			fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: Expected Closing Bracket.")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: Expected Non-empty Expression.")
		os.Exit(1)
	case Tokens.SyntaxError:
		fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: unknown character.")
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: unknown token.")
		os.Exit(1)
	} // switch --- Token
} // LExpr

// VarExpr
// Prints the read in variable and resets the lexeme and ExpectedVariable members.
// Since a <lvar> is a terminal, the next token is obtained at the end of this function.
// context: Contains the whole expression.
func VarExpr(context *Globals.Vars) {
	if context.DebugMode {
		fmt.Println("   VarExpr called")
	}
	context.Tree.AddToken(Tokens.TokenVariable, string(context.Lexeme), CalcBrack(context))
	context.Lexeme = nil
	context.ExpectVariable = false
	context.Token = LexicalAnalyser.LexicalAnalyser(context)
} // VarExpr

// TypeExpr
// Obtains the next token from LexicalAnalyser and determines the next step in recursive descent accordingly.
// Contains all possible continuations for <type> namely <uvar>, ( <type> ) and <type> "->" <type>.
// The <type> "->" <type> continuation is handled using TypeFunction which "Peeks" to see if the next token
// corresponds with a "->", if so an extra call to TypeExpr is made after parsing the "->".
// context: Contains the whole expression.
func TypeExpr(variables *Globals.Vars) {
	if variables.DebugMode {
		fmt.Println("   TypeExpr called")
	}
	switch variables.Token {
	case Tokens.TokenUVar:
		variables.Tree.AddToken(Tokens.TokenUVar, string(variables.Lexeme), CalcBrack(variables))
		// Reset.
		variables.Lexeme = nil
		variables.Token = LexicalAnalyser.LexicalAnalyser(variables)

		// Check for possible <type> "->" <type> continuation.
		if TypeFunction(variables) {
			TypeExpr(variables)
		}
		break
	case Tokens.TokenLeftBracket:
		// Increment shows that there is an extra open bracket.
		variables.CountBrackets++

		variables.Token = LexicalAnalyser.LexicalAnalyser(variables)

		// Expected nested TypeExpr.
		TypeExpr(variables)
		if variables.Token == Tokens.TokenRightBracket {
			variables.CountBrackets--
			variables.Token = LexicalAnalyser.LexicalAnalyser(variables)
			if TypeFunction(variables) {
				TypeExpr(variables)
			}
		}
		break
	case Tokens.TokenRightBracket:
		fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: Expected UVar.")
		os.Exit(1)
	case Tokens.TokenVariable:
		fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: LVar cannot be parsed in Type Expression.")
		os.Exit(1)
	case Tokens.TokenLambda:
		fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: Lambda cannot be parsed in Type Expression.")
		os.Exit(1)
	case Tokens.LexicalEndOfLine:
		fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: Type Expression cannot be empty.")
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "%s\n", "Syntax error: Expected Type Expression")
		os.Exit(1)
	}
}

// TypeFunction
// "Peeks" in the string to be lexically analysed and checks whether the next token
// is a TypeFunction. If so, the token is parsed and Recursive Descent continues accordingly.
// Else, the "Peek" is reverted since a token other than a TypeFunction was found,
// which allows the LexicalAnalyser to read this token again on the next call.
// context: Contains the whole expression.
func TypeFunction(context *Globals.Vars) bool {
	if context.Token == Tokens.TokenFunction {
		context.Tree.AddToken(Tokens.TokenFunction, "->", 0)
		//get the next token
		context.Token = LexicalAnalyser.LexicalAnalyser(context)
		return true
	}
	return false
}

// CalcBrack
// Given the current index in the global variables we determine the amount of brackets among this token
// an opening bracket on the left increments the counter
// a closing bracket on the right decrements the counter
// returns an integer, that counts the difference between opening and closing brackets around the token
// context: Contains the whole expression.
func CalcBrack(context *Globals.Vars) int {
	var counter = 0
	var index = context.Index
	for i := index - 1; i >= 0; i-- {
		if context.CurrentLine[i] == '(' {
			counter++
			// we ignore the first lambda to the left since the counter for a lambda token is stored in the lambda variable
		} else if context.CurrentLine[i] == ' ' || (context.CurrentLine[i] == '\\' || context.CurrentLine[i] == 'λ') {
			continue
		} else {
			break
		}
	} // for --- walk back
	for i := index + len(context.Lexeme); i < len(context.CurrentLine); i++ {
		if context.CurrentLine[i] == ')' {
			counter--
		} else if context.CurrentLine[i] == ' ' {
			continue
		} else {
			break
		}
	} // for --- walk forward
	return counter
} // CountBrack
