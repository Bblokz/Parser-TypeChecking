/*
 * Parser and Lexical Analyser Global.go
 * Copyright (C) 2021-2023 Bas Blokzijl, Leiden, The Netherlands.
 */

/*
 * Contains all variables needed for context when parsing and type checking.
 * Variables are passed on by reference.
 */

package Globals

import( 
    ParseTree "Parser-TypeChecking/Parsetree"
)

// VariableType
// Denotes the Type and Var name both in string form.
type VariableType struct {
	// The token value of this node.
	VarName string
	TypeName string
}

type VarTypeList struct{
    list []VariableType
}

// FindInList
// Returns index of variable in the type list exactly if it has a type in this context, else returns -1
// variableName: The variable of which the type will be searched for.
func (context *VarTypeList) FindInList (variableName string) int{
    if (len(context.list) == 0 || len(variableName) == 0){
        return -1
    }
    for i := len(context.list) - 1; i >= 0; i-- {   
        //if a var name comes up in the list see if entire name is the same.
        if context.list[i].VarName[0] == variableName[0]{
            for ind := 0; ind <= len(context.list[i].VarName); ind++ {
                if ind == len(context.list[i].VarName) { //-1????? todo
                    return i
                }else if context.list[i].VarName[ind] != variableName[ind] {
                    break
                }
            }
        }
	}
    return -1
}

// AddVarType
// Function to add the provided variable with typeName to context.
// variableName: Denotes the variable.
// typeName: The type of the variable.
func (context *VarTypeList) AddVarType (variableName string, typeName string) {
    var newvartype VariableType
    newvartype.VarName = variableName
    newvartype.TypeName = typeName
    context.list = append(context.list, newvartype)
}

// GetLast
// Function to find rightmost type for variable in context.
// variableName: Name of the variable.
func (context *VarTypeList) GetLast (variableName string) string{
    var foundindex = context.FindInList(variableName)
    if (foundindex == -1){
        return "err"
    }
    return context.list[foundindex].TypeName
}


type Vars struct {
	DebugMode bool
    
    // List of known variables, order is important.
    Context VarTypeList

	// Abstract syntax tree containing all tokens after parsing.
	Tree ParseTree.ParseTree

	// Used to store the line that is currently analysed.
	CurrentLine []rune

	// Index on CurrentLine.
	Index int

	// Character currently read from line.
	ReadChar rune

	// Character class of the currently read character.
	CharClass int

	// Token read in by lexical analyser.
	Token int

	// PrevToken is the previous token obtained by the LexicalAnalyser.
	PrevToken int

	// Lexeme that is currently analysed.
	Lexeme []rune

	// The last token was a variable and we just encountered a whitespace.
	ExpectVariable bool

	// The last token was a lambda variable we need ^.
	ExpectUp bool

	// counter to keep track of the amount of closing and open brackets.
	CountBrackets int
}











