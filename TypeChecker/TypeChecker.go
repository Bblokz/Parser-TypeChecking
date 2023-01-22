/*
 * Parser and Lexical Analyser TypeChecker.go
 * Copyright (C) 2021-2023 Bas Blokzijl Leiden, The Netherlands.
 */

package TypeChecker

import (
    "Parser-TypeChecking/Globals"
    "fmt"
    "os"
    "unicode"
)

/* calcParts
 * Uses the standard-output on ParseTree to get the final slices used for typechecking
 * variables: passed on arguments containing the TreeOutput.
 * expressionSlice: by-reference, contains the expression.
 * typeSLice: by-reference, contains the types.
 */
func calcParts(variables *Globals.Vars, expressionSlice *string, typeSlice *string) {
    *expressionSlice = variables.Tree.SubTreeToStandardOutput(1)
    *typeSlice = variables.Tree.SubTreeToStandardOutput(variables.Tree.IndexDoubleDot)
}

const (
	VariableRule = iota
	ApplicationRule 
	LamdbaRule 
)

/* findE
 * Finds the part where the brackets close and returns encapsulated Expression.
 * expressionSlice: The analysed expression.
 * index: Index in the expression at which the opening bracket occurred.
 * firstE: Whether this the first expression.
 */
func findE(expressionSlice string, index int, firstE bool) string {
    word := ""
    var currLetter byte
    countbrack := 1
    for i := index; i < len(expressionSlice); i++ {
        currLetter = expressionSlice[i]
        if currLetter == '(' {
            countbrack++
        } else if currLetter == ')' || (firstE && currLetter == ' '){
            countbrack--
        }
        if countbrack == 0 {
            break
        }
        word +=  string(currLetter)
    }
    return word
}

/* findWord
 * uses the first provided index of the word to return entire word as string, which stopped at first non-letter
 * expressionSlice: Part of the total expression, contains the word we are looking for.
 * index: Index at which the word starts in the provided string.
 */
func findWord(expressionSlice string, index int) string{
    word := ""
    // walk until end of the expression.
    for _, r := range expressionSlice[index:] {
        if !unicode.IsLetter(r) {
            return word
        }
         word += string(r)
     }
     return word
}


/* findRule
 * Decide which rule needs to be applied next to find the type for the expression.
 * expressionSlice: The total expression.
 */
func findRule(expressionSlice string) int{
    if len(expressionSlice) >=1 && expressionSlice[0] == '(' {
        if len(expressionSlice) >=2 && expressionSlice[1] == '\\' { 
            return LamdbaRule
        } else {
            return ApplicationRule
        }
    } else {
        return VariableRule
    }
}

/* findType
 * Finds a unique type for given expression, using context in variables.
 * Recursive function that slices up the expression slice until it is a variable.
 * variables: Contains expression context.
 * expressionSlice: The entire expression.
 */
func findType(variables *Globals.Vars, expressionSlice string) string{
    T1 := ""
    T2 := ""
    E1 := ""
    E2 := ""
    foundVar := ""
    word := ""
    switch findRule(expressionSlice) {
    case VariableRule:

        word = variables.Context.GetLast(expressionSlice) 
        if variables.Context.FindInList (expressionSlice) == -1 {
            fmt.Fprintf(os.Stderr, "%s\n", "Typecheck error: Variable not in context")
            os.Exit(1)
        }
        word = variables.Context.GetLast(expressionSlice)
        break
    case ApplicationRule:
        // Break expression apart and find type of the parts.
        E1 = findE(expressionSlice, 1, true)
        E2 = findE(expressionSlice, 2 + len(E1), false)
        // Recursive.
        T1 = findType(variables, E1)
        T2 = findType(variables, E2)
        if T1[:(len(T2) + 3)] != (T2 + "->"){
            fmt.Fprintf(os.Stderr, "%s\n", "Typecheck error: E1 should find a function type")
        }
        if T2 != T1[:len(T2)] { 
            fmt.Fprintf(os.Stderr, "%s\n", "Typecheck error: codomain E1 not the same as type of E2")
        }
        word = T1[(len(T2) + 3):(len(T1)-1)]
        break
    case LamdbaRule:
        // Add lambda type to the type that is applied and find type of applied term.
        foundVar = findWord(expressionSlice, 2)
        T1 = findE(expressionSlice, 4+len(foundVar), false)
        // Add variable to context.
        variables.Context.AddVarType(foundVar, T1)
        E1 = findE(expressionSlice, 5+len(foundVar)+len(T1), false)
        T2 = findType(variables, E1)
            word = "(" + T1 + "->"+ T2 + ")"
    default:
        break
    }
    return word
}

/* TypeCheker
 * Uses findType() recursively on the expression part of the judgement and finds whether
 * the found type is the same as given type.
 * variables: Provides context.
 */
func TypeChecker(variables *Globals.Vars) {
    var finalExpression, finalType string
    calcParts(variables, &finalExpression, &finalType)
    
    //does not account for redundant brackets
    if finalType == findType(variables, finalExpression){
        fmt.Println("Type checks out")
    } else {
        fmt.Println("Does not type check")
    }
}
