/*
 * Parser and Lexical Analyser ParseTree.go
 * Copyright (C) 2021-2023 Bas Blokzijl Leiden, The Netherlands.
 */

package parsetree

import (
	"Parser-TypeChecking/Tokens"
	"fmt"
	"strings"
)

type Node struct {
	// Token the token value of this node.
	Token int
	// Lexeme string used to print e.g. name of a variable.
	Lexeme string
	// Depth of the AST.
	Depth int
	// IsInApplication denotes whether this node is a direct child of an application.
	IsInApplication bool
	// IsSecondInApplication denotes which child of the application contains this node.
	IsSecondInApplication bool
	// BracketCounter denotes the amount of brackets around this node, >0 if there are more opening than closing brackets.
	BracketCounter int
}

// TokenToString
// Returns the string associated with this token.
// used to print tree.
func (node Node) TokenToString() string {
	switch node.Token {
	case Tokens.TokenLambda:
		return "Lam:"
	case Tokens.TokenRightBracket:
	case Tokens.TokenLeftBracket:
		return "    "
	case Tokens.TokenVariable:
		return "Var:"
	}
	return "    "
} // TokenToString

// ParseTree
// AST containing the parsed expression.
type ParseTree struct {
	// Indicates the Index in the string slice at which the DoubleDot is stored.
	IndexDoubleDot int
	// All Nodes in the tree.
	Nodes []Node
	// Depth at which the last node was added.
	currentDepth int
}

// AddToken
// Method for ParseTree: creates a node with the token and lexeme information and adds it to the tree.
// token: Type of the token.
// lexeme: the provided lexeme as string.
// bracketsCounter: How many brackets are currently opened.
func (tree *ParseTree) AddToken(token int, lexeme string, bracketsCounter int) {
	if len(tree.Nodes) > 0 {
		// Get the previous node.
		var prevNode = tree.Nodes[len(tree.Nodes)-1]
		// This becomes the new node to add.
		var newNode Node
		// Index that specifies at which position we can find the closest node that has opening brackets around it.
		var indexBracketNode = -1
		// Index that specifies at which position we can find the closest node that is already part of an application.
		var indexAppliedNode = -1

		if token == Tokens.TokenFunction {
			// Functions are added immediately.
			newNode = Node{token, lexeme, tree.currentDepth, false, false, bracketsCounter}
			tree.Nodes = append(tree.Nodes, newNode)
			return
		}
		if token == Tokens.TokenDoubleDot {
			// Added immediately, needs extra information for adding types correctly.
			newNode = Node{token, lexeme, tree.currentDepth, false, false, bracketsCounter}
			tree.Nodes = append(tree.Nodes, newNode)
			return
		}
		switch prevNode.Token {
		case Tokens.TokenDoubleDot:
			// Remove DoubleDot from old location.
			tree.Nodes = tree.Nodes[:len(tree.Nodes)-1]
			// Apply the first part of the Judgement to the Judgment Expression.
			tree.Nodes[0].IsInApplication = true
			tree.Nodes[0].IsSecondInApplication = false
			tree.IncrementDepthFromIndex(0)
			// Place DoubleDot in front
			prevNode.Depth = 0
			tree.InjectNodeAtIndex(prevNode, 0)
			// Apply the second part of the Judgement to the Judgement Expression.
			tree.currentDepth = 1
			newNode = Node{token, lexeme, tree.currentDepth, false, false, bracketsCounter}
			tree.Nodes = append(tree.Nodes, newNode)
			break
		case Tokens.TokenLambda:
			// The new node is a function argument and added on the same level as the lambda.
			newNode = Node{token, lexeme, tree.currentDepth, false, false, bracketsCounter}
			tree.Nodes = append(tree.Nodes, newNode)
			break
		case Tokens.TokenFunction:
			// Find brackets
			if tree.Nodes[len(tree.Nodes)-2].BracketCounter < 0 {
				for i := len(tree.Nodes) - 2; i >= 0; i-- {
					if tree.Nodes[i].BracketCounter > 0 {
						tree.Nodes[len(tree.Nodes)-2].BracketCounter += tree.Nodes[i].BracketCounter
						tree.Nodes[i].BracketCounter = 0
						if tree.Nodes[len(tree.Nodes)-2].BracketCounter > 0 {
							tree.Nodes[i].BracketCounter += tree.Nodes[len(tree.Nodes)-2].BracketCounter
							tree.Nodes[len(tree.Nodes)-2].BracketCounter = 0
						}
						if tree.Nodes[len(tree.Nodes)-2].BracketCounter == 0 {
							indexBracketNode = i
							break
						}
					}
				} // for
			} else {
				// Right associative; find first already applied node.
				if tree.IndexDoubleDot >= 0 {
					// DoubleDot has been added; need restriction on indexAppliedNode.
					for i := len(tree.Nodes) - 2; i >= tree.IndexDoubleDot; i-- {
						if tree.Nodes[i].IsSecondInApplication == true && (tree.Nodes[i].Token == Tokens.TokenUVar ||
							tree.Nodes[i].Token == Tokens.TokenFunction) {
							indexAppliedNode = i
						}
					} // for
				} else {
					for i := len(tree.Nodes) - 2; i >= 0; i-- {
						if tree.Nodes[i].IsSecondInApplication == true && (tree.Nodes[i].Token == Tokens.TokenUVar ||
							tree.Nodes[i].Token == Tokens.TokenFunction) {
							indexAppliedNode = i
						}
					} // for
				}
				if indexAppliedNode >= 0 {
					tree.UVarFunctToNestedFunc(&prevNode, indexAppliedNode, token, lexeme, bracketsCounter)
				} else {
					tree.UVarFuncToPrevUVar(&prevNode, token, lexeme, bracketsCounter)
				}
			}
			// Whether we found a valid indexBracketNode.
			if indexBracketNode >= 0 {
				tree.UVarFuncToBracketNode(&prevNode, indexBracketNode, token, lexeme, bracketsCounter)
			}
			break
		case Tokens.TokenUVar:
			tree.ApplyToClosestLambda(token, lexeme, bracketsCounter)
			break
		case Tokens.TokenVariable:
			if tree.IsPrevLambda(token, lexeme, bracketsCounter) {
				break
			}
			// Find node with opening brackets.
			for i := len(tree.Nodes) - 1; i >= 0; i-- {
				if tree.Nodes[i].BracketCounter > 0 {
					indexBracketNode = i
					break
				}
			}
			// Check if there indeed exists a node with opening brackets.
			if indexBracketNode >= 0 {
				tree.ApplyToBracketNode(indexBracketNode, token, lexeme, bracketsCounter)
			} else {
				if prevNode.IsInApplication {
					tree.ApplyToClosestNotAppliedApplication(token, lexeme, bracketsCounter)
				} else {
					tree.ApplyToPrevNode(token, lexeme, bracketsCounter)
				}
			} // no node with open brackets found.
		} // switch
	} else {
		// Add first node.
		var newNode = Node{token, lexeme, tree.currentDepth, false, false, bracketsCounter}
		tree.Nodes = append(tree.Nodes, newNode)
	}
}

// UVarFunctToNestedFunc
// Uses the previously added -> and applies a recursive arrow function to it.
// prevNode: Last added node.
// inexAppliedNode: Where the applied node lives.
// token: of the new term.
// lexeme: of the new term.
// bracketsCounter: denotes the amount of opened brackets.
func (tree *ParseTree) UVarFunctToNestedFunc(prevNode *Node, indexAppliedNode int, token int, lexeme string, bracketsCounter int) {
	// Remove function from previous location.
	var newNode Node
	tree.Nodes = tree.Nodes[:len(tree.Nodes)-1]
	prevNode.Depth = tree.Nodes[indexAppliedNode].Depth
	prevNode.IsInApplication = true
	prevNode.IsSecondInApplication = true
	// Increment depth since this node moves one level deeper, underneath the previous node.
	tree.Nodes[indexAppliedNode].Depth++
	tree.currentDepth++
	// Becomes the first child in the new level.
	tree.Nodes[indexAppliedNode].IsSecondInApplication = false
	// Inject the function at the applied location.
	tree.InjectNodeAtIndex(*(prevNode), indexAppliedNode)
	// Add the new Node, applied in the injected function as the second child.
	newNode = Node{token, lexeme, tree.currentDepth, true, true, bracketsCounter}
	tree.Nodes = append(tree.Nodes, newNode)
}

// UVarFuncToPrevUVar
// Uses the previously added -> to apply the new node (of type UVar) with the second to last Uvar in the slice.
func (tree *ParseTree) UVarFuncToPrevUVar(prevNode *Node, token int, lexeme string, bracketsCounter int) {
	var newNode Node
	tree.Nodes = tree.Nodes[:len(tree.Nodes)-1]
	// Pass on the brackets
	prevNode.BracketCounter = tree.Nodes[len(tree.Nodes)-1].BracketCounter
	tree.Nodes[len(tree.Nodes)-1].BracketCounter = 0
	// Move node one level deeper
	tree.Nodes[len(tree.Nodes)-1].Depth++
	tree.currentDepth++
	// Move node into application
	tree.Nodes[len(tree.Nodes)-1].IsInApplication = true
	tree.Nodes[len(tree.Nodes)-1].IsSecondInApplication = false
	// Inject function at correct position
	tree.InjectNodeAtIndex(*(prevNode), len(tree.Nodes)-1)
	// Add the new Node, applied in the injected function as the second child.
	newNode = Node{token, lexeme, tree.currentDepth, true, true, bracketsCounter}
	tree.Nodes = append(tree.Nodes, newNode)
	return
}

// UVarFuncToBracketNode
// Handles the correct placement for a UVar when there exists another UVar with Brackets.
// The prevNode has to be of token: Function.
// Moves the Function to "Apply" the new Uvar to the Uvar with Brackets by inserting it in the correct place.
// The brackets are moved on to the next node so no information is lost.
// prevNode: Last added node.
// inexAppliedNode: Where the applied node lives.
// token: of the new term.
// lexeme: of the new term.
// bracketsCounter: denotes the amount of opened brackets.
func (tree *ParseTree) UVarFuncToBracketNode(prevNode *Node, indexBracketNode int, token int, lexeme string, bracketsCounter int) {
	var newNode Node
	var isAlreadyApplied = false
	if tree.Nodes[indexBracketNode].IsInApplication && !tree.Nodes[indexBracketNode].IsSecondInApplication {
		isAlreadyApplied = true
	}
	// Remove function from previous location.
	tree.Nodes = tree.Nodes[:len(tree.Nodes)-1]
	// Prev function is inserted at the original depth
	prevNode.Depth = tree.Nodes[indexBracketNode].Depth
	// Apply the first part of the Function on the prev node.
	tree.Nodes[indexBracketNode].IsInApplication = true
	tree.Nodes[indexBracketNode].IsSecondInApplication = false
	tree.IncrementDepthFromIndex(indexBracketNode)
	// Pass on the bracketCounter.
	prevNode.BracketCounter = tree.Nodes[indexBracketNode].BracketCounter
	tree.Nodes[indexBracketNode].BracketCounter = 0
	tree.InjectNodeAtIndex(*(prevNode), indexBracketNode)
	if isAlreadyApplied {
		tree.currentDepth = prevNode.Depth
	} else {
		tree.currentDepth = prevNode.Depth + 1
	}

	newNode = Node{token, lexeme, tree.currentDepth, true, true, bracketsCounter}
	tree.Nodes = append(tree.Nodes, newNode)
}

// ApplyToClosestLambda
// Applies the new Node to the closest lambda near the end of the Node Slice.
// Passes on the possible remaining brackets on the lambda to the application.
// Removes any redundant brackets that are closed by the brackets provided with the new Node.
// token: of the new term.
// lexeme: of the new term.
// bracketsCounter: denotes the amount of opened brackets.
func (tree *ParseTree) ApplyToClosestLambda(token int, lexeme string, bracketsCounter int) {
	var newNode Node
	newNode = Node{token, lexeme, tree.currentDepth, true, true, bracketsCounter}
	for i := len(tree.Nodes) - 1; newNode.BracketCounter < 0 && i >= 0; i-- {
		if tree.Nodes[i].BracketCounter > 0 {
			tree.RemoveRedundantBrackets(&tree.Nodes[i], &newNode)
		}
	}
	for i := len(tree.Nodes) - 2; i >= 0; i-- {
		if tree.Nodes[i].Token == Tokens.TokenLambda {
			// We create a new application to nest the lambda in, this means we pass on any possible remaining
			// brackets on the lambda.
			newNode = Node{Tokens.Application, "Apply", tree.Nodes[i].Depth, tree.Nodes[i].IsInApplication,
				tree.Nodes[i].IsSecondInApplication, tree.Nodes[i+1].BracketCounter}
			tree.Nodes[i+1].BracketCounter = 0
			tree.currentDepth = tree.Nodes[i].Depth + 1
			tree.IncrementDepthFromIndex(i)
			tree.Nodes[i].IsInApplication = true
			tree.Nodes[i].IsSecondInApplication = false
			tree.InjectNodeAtIndex(newNode, i)
			newNode = Node{token, lexeme, tree.currentDepth, true, true, bracketsCounter}
			tree.Nodes = append(tree.Nodes, newNode)
			break
		}
	}
}

// IsPrevLambda
// Checks whether the previous two tokens have made a lambda expr, if so, then the new token will be nested
// underneath this lambda. The new token is created with the provided token, lexeme and bracketscounter.
// The depth is determined using currentDepth in the tree.
// token: of the new term.
// lexeme: of the new term.
// bracketsCounter: denotes the amount of opened brackets.
func (tree *ParseTree) IsPrevLambda(token int, lexeme string, bracketsCounter int) bool {
	// Check if the previous 2 tokens have made up a lambda expr.
	var newNode Node
	if len(tree.Nodes) > 1 && tree.Nodes[len(tree.Nodes)-2].Token == Tokens.TokenLambda {
		// The new node is the first underneath the lambda expression, hence we increment the depth.
		tree.currentDepth++
		// Add new node to tree.
		newNode = Node{token, lexeme, tree.currentDepth, false, false, bracketsCounter}
		tree.Nodes = append(tree.Nodes, newNode)
		return true
	}
	return false
}

// ApplyToBracketNode
// Applies the new node to the node that contains brackets, provided by the index in the node slice.
// If the new node has closing brackets it will close as many as possible open brackets on the node to which
// the new node is applied. Any possible remaining brackets are passed on to the application node.
// The depth of the tree is adjusted accordingly.
// indexBracketNode: Where the node at which the openbracket start, lives.
// token: of the new term.
// lexeme: of the new term.
// bracketsCounter: denotes the amount of opened brackets.
func (tree *ParseTree) ApplyToBracketNode(indexBracketNode int, token int, lexeme string, bracketsCounter int) {
	var newNode Node
	var oldDepthAtApplication = tree.Nodes[indexBracketNode].Depth
	tree.IncrementDepthFromIndex(indexBracketNode)
	// Create the new node at new depth.
	tree.currentDepth = oldDepthAtApplication + 1
	newNode = Node{token, lexeme, tree.currentDepth, true, true, bracketsCounter}
	// Calculate remaining brackets, if the new node contains closing brackets it may remove brackets
	// on the node to which we apply.
	var remainingBrackets = tree.RemoveRedundantBrackets(&tree.Nodes[indexBracketNode], &newNode)
	// Create application node, pass on possible remaining brackets.
	var application = Node{Tokens.Application, "Apply", oldDepthAtApplication, false, false, remainingBrackets}
	// remove the brackets from the old node since they are passed on
	tree.Nodes[indexBracketNode].BracketCounter = 0
	if indexBracketNode != 0 && tree.Nodes[indexBracketNode-1].Token == Tokens.TokenLambda {
		// We have only incremented the depth from the lambda variable and further, not the lambda itself!
		tree.Nodes[indexBracketNode-1].Depth++
		// The brackets are stored in the variable which is part of the lambda expression but the application
		// needs to be injected before the lambda since the lambda token and its variable make up one line.
		tree.InjectNodeAtIndex(application, indexBracketNode-1)
	} else {
		// Set application or bracketVariable nested in application.
		tree.Nodes[indexBracketNode].IsInApplication = true
		tree.Nodes[indexBracketNode].IsSecondInApplication = false
		tree.InjectNodeAtIndex(application, indexBracketNode)
	}
	// Add new node to the tree.
	tree.Nodes = append(tree.Nodes, newNode)
	return
}

// ApplyToClosestNotAppliedApplication
// Searches backwards in the node slice for the first application in the tree that is not applied itself.
// Then the new node is applied to this application. Since this will always be done only if there are no bracketnodes
// to apply to, we do not need to pass on any brackets.
// token: of the new term.
// lexeme: of the new term.
// bracketsCounter: denotes the amount of opened brackets.
func (tree *ParseTree) ApplyToClosestNotAppliedApplication(token int, lexeme string, bracketsCounter int) {
	var newNode Node
	for i := len(tree.Nodes) - 1; i >= 0; i-- {
		if tree.Nodes[i].Token == Tokens.Application && !(tree.Nodes[i].IsInApplication) {
			// Increment depth for this application and subtree.
			var oldDepthAtApplication = tree.Nodes[i].Depth
			tree.IncrementDepthFromIndex(i)
			// Set application nested in application.
			tree.Nodes[i].IsInApplication = true
			tree.Nodes[i].IsSecondInApplication = false
			// Inject nested application.
			newNode = Node{Tokens.Application, "Apply", oldDepthAtApplication, false, false, 0}
			tree.InjectNodeAtIndex(newNode, i)
			// Add the new node.
			tree.currentDepth = oldDepthAtApplication + 1
			newNode = Node{token, lexeme, tree.currentDepth, true, true, bracketsCounter}
			tree.Nodes = append(tree.Nodes, newNode)
			return
		} // if
	} // for --- search non nested application
}

// ApplyToPrevNode
// The previous node is not yet in an application, we insert an application and apply the
// new node to the previous node.
// token: of the new term.
// lexeme: of the new term.
// bracketsCounter: denotes the amount of opened brackets.
func (tree *ParseTree) ApplyToPrevNode(token int, lexeme string, bracketsCounter int) {
	var newNode Node
	newNode = Node{Tokens.Application, "Apply", tree.currentDepth, false, false, 0}
	tree.InjectNodeAtIndex(newNode, len(tree.Nodes)-1)
	// Set prev node as first in the application.
	tree.Nodes[len(tree.Nodes)-1].IsInApplication = true
	tree.Nodes[len(tree.Nodes)-1].IsSecondInApplication = false
	tree.Nodes[len(tree.Nodes)-1].Depth++
	tree.currentDepth++
	newNode = Node{token, lexeme, tree.currentDepth, true, true, bracketsCounter}
	tree.Nodes = append(tree.Nodes, newNode)
	return
}

// ClearTree
// Resets the tree for parsing the next line.
// Removes all Nodes and sets depth back to 0.
func (tree *ParseTree) ClearTree() {
	tree.currentDepth = 0
	tree.Nodes = nil
	tree.IndexDoubleDot = -1
}

// FindLambdas
// returns a slice containing all the indices of the lambdas in the node slice of our AST
func (tree *ParseTree) FindLambdas() []int {
	var indices []int
	for i := 0; i < len(tree.Nodes); i++ {
		if tree.Nodes[i].Token == Tokens.TokenLambda {
			indices = append(indices, i)
		}
	}
	return indices
} // FindLambdas

// InjectNodeAtIndex
// Adds the provided node to the tree at the given index, moves all other Nodes to the right of the slice
// and extends the slice by appending the last element.
// injectNode: node to inject into the tree.
// index: Denotes the position to inject to.
func (tree *ParseTree) InjectNodeAtIndex(injectNode Node, index int) {
	var NodesMoved []Node
	for i := index; i < len(tree.Nodes); i++ {
		NodesMoved = append(NodesMoved, tree.Nodes[i])
	}
	tree.Nodes[index] = injectNode
	for i := 0; i < len(NodesMoved)-1; i++ {
		tree.Nodes[index+1+i] = NodesMoved[i]
	}
	// append last element, since
	tree.Nodes = append(tree.Nodes, NodesMoved[len(NodesMoved)-1])
} // InjectNodeAtIndex

// RemoveRedundantBrackets
// Removes the difference in bracket-counters between the two Nodes.
// Nodes have to be the Left and Right term in an application for redundant brackets to be removed.
// Returns the remaining opening brackets on nodeLeft.
func (tree *ParseTree) RemoveRedundantBrackets(nodeLeft *Node, nodeRight *Node) int {
	left := nodeLeft.BracketCounter
	right := nodeRight.BracketCounter
	if left > 0 && right < 0 {
		if left > (-right) {
			// More open than closed brackets.
			nodeLeft.BracketCounter += nodeRight.BracketCounter
			nodeRight.BracketCounter = 0

		} else if left == (-right) {
			// #Open brackets == #closed brackets.
			nodeLeft.BracketCounter = 0
			nodeRight.BracketCounter = 0
		} else {
			// More closed than open brackets.
			nodeLeft.BracketCounter = 0
			nodeRight.BracketCounter += left
		}
	}
	return nodeLeft.BracketCounter
} // RemoveRedundantBrackets

// IncrementDepthFromIndex
// Increments the depth of all Nodes up and from the provided index.
func (tree *ParseTree) IncrementDepthFromIndex(index int) {
	for i := index; i < len(tree.Nodes); i++ {
		tree.Nodes[i].Depth++
	}
} // IncrementDepthFromIndex

// PrintTree
// Prints the entire AST to the console.
func (tree ParseTree) PrintTree() {
	if len(tree.Nodes) == 0 {
		fmt.Printf("\n Tree is empty but ")
		fmt.Println("works fine.")
	} else {
		var currentNode Node
		for i := 0; i < len(tree.Nodes); i++ {

			currentNode = tree.Nodes[i]
			if i > 2 && currentNode.Depth == 1 {
				break
			}
			var spaces = strings.Repeat("---", currentNode.Depth)

			switch currentNode.Token {
			case Tokens.TokenLambda:
				fmt.Printf("%s %s %s%s%s\n", spaces, currentNode.TokenToString(),
					currentNode.Lexeme, tree.Nodes[i+1].Lexeme, "^")
				i++
				break
			case Tokens.Application:
				fmt.Printf("%s %s:\n", spaces, currentNode.Lexeme)
				break
			case Tokens.TokenDoubleDot:
				fmt.Printf("%s %s:\n", spaces, currentNode.Lexeme)
				break
			default:
				fmt.Printf("%s %s %s\n", spaces, currentNode.TokenToString(), currentNode.Lexeme)
			} // switch
		} // for --- Nodes in tree

		fmt.Println("Type:")
		for i := tree.IndexDoubleDot; i < len(tree.Nodes); i++ {
			var spaces = strings.Repeat("---", tree.Nodes[i].Depth-1)
			fmt.Printf("%s %s %s\n", spaces, tree.Nodes[i].TokenToString(), tree.Nodes[i].Lexeme)
		}
	} // if --- tree not empty
} // PrintTree

func InjectSlice(index int, base *[]Node, injection []Node) {
	var adjust = (*base)[0].Depth
	for i := 0; i < len((*base)); i++ {
		(*base)[i].Depth -= adjust
	}
	adjust = injection[0].Depth
	// adjust depth
	for i := 0; i < len(injection); i++ {
		// For compensating the application of which this expr is a child.
		injection[i].Depth += (*base)[index].Depth - adjust
	}
	var NodesMoved []Node
	for i := index + 1; i < len((*base)); i++ {
		NodesMoved = append(NodesMoved, (*base)[i])
	}
	var prevNodes []Node
	for i := 0; i < index; i++ {
		prevNodes = append(prevNodes, (*base)[i])
	}
	for j := 0; j < len(injection); j++ {
		prevNodes = append(prevNodes, injection[j])
	}
	for j := 0; j < len(NodesMoved); j++ {
		prevNodes = append(prevNodes, NodesMoved[j])
	}
	*base = prevNodes
}
func (tree ParseTree) SubTreeToStandardOutput(index int) string {

	// Output string in standard format.
	var Output string
	// Whether the last added lexeme was a UVar.
	var lastIsUvar = false
	// Whether the last added lexeme was a LVariable.
	var lastIsVariable = false

	// Whether we are currently adding lexemes that are nested in a lambda.
	var isInLambda = false

	// Whether the extra opening bracket in a lambda needs to be adjusted.
	var ajdustOpenLambda = false
	// Helper-String used to remove redundant brackets.
	var adjustBrackets string

	// Difference in depth between the new Node and the previous one.
	var deltaDepth int
	// String supplied with the brackets of the Node according to delta depth.
	var brackets string
	// Depth at previous node.
	var previousDepth int
	judgementDepth := -1
	if tree.Nodes[index].Token == Tokens.TokenLambda {
		index++
	}
	previousDepth = tree.Nodes[index].Depth

	if tree.Nodes[index].Token == Tokens.TokenDoubleDot {
		judgementDepth = tree.Nodes[index].Depth
	}
	var subtreeDepth = tree.Nodes[index].Depth
	for i := index + 1; i < len(tree.Nodes) && tree.Nodes[i].Depth > subtreeDepth; i++ {

		deltaDepth = tree.Nodes[i].Depth - previousDepth
		for j := deltaDepth; j < 0; j++ {
			brackets += ")"
		}
		for j := deltaDepth; j > 0; j-- {
			brackets += "("
			j--
		}
		if (tree.Nodes[i].Token == Tokens.TokenUVar ||
			tree.Nodes[i].Token == Tokens.TokenFunction) &&
			(tree.Nodes[i].Depth == judgementDepth+1) {
			brackets += ":"
		}
		switch tree.Nodes[i].Token {
		case Tokens.TokenLambda:
			if isInLambda {
				adjustBrackets = ""
				for i := 1; i < len(brackets); i++ {
					adjustBrackets += string(brackets[i])
				}
				brackets = adjustBrackets
			}
			isInLambda = true
			ajdustOpenLambda = true
			Output += brackets
			Output += "\\"

			Output += tree.Nodes[i+1].Lexeme
			Output += "^"
			lastIsUvar = false
			lastIsVariable = false
			i++
			break
		case Tokens.TokenVariable:
			if lastIsVariable && !(tree.Nodes[i-1].IsSecondInApplication) {
				// This TokenVariable is applied to the last TokenVariable, we need a space.
				Output += " "
				lastIsVariable = false
			} else {
				lastIsVariable = true
			}
			previousDepth = tree.Nodes[i].Depth
			lastIsUvar = false
			if isInLambda {
				adjustBrackets = ""
				for i := 1; i < len(brackets); i++ {
					adjustBrackets += string(brackets[i])
				}
				brackets = adjustBrackets
				if !(tree.Nodes[i-2].Token == Tokens.TokenUVar || tree.Nodes[i-2].Token == Tokens.TokenFunction) {
					// Remove prev UVar to adjust brackes around it
					// Since this is a single Uvar
					Output = Output[:(len(Output) - len(tree.Nodes[i-1].Lexeme))]
					Output += "("
					Output += tree.Nodes[i-1].Lexeme
					Output += ")"
				}
			}
			isInLambda = false
			Output += brackets
			Output += tree.Nodes[i].Lexeme

			break
		case Tokens.TokenUVar:
			if lastIsUvar {
				if deltaDepth < 0 {
					Output += brackets
					Output += "->"

				} else {
					Output += "->"
					Output += brackets
				}

				Output += tree.Nodes[i].Lexeme
				lastIsUvar = true
				lastIsVariable = false
				previousDepth = tree.Nodes[i].Depth
				break
			}
			if isInLambda && ajdustOpenLambda {
				adjustBrackets = ""
				for i := 1; i < len(brackets); i++ {
					adjustBrackets += string(brackets[i])
				}
				brackets = adjustBrackets
				ajdustOpenLambda = false
			}
			Output += brackets
			Output += tree.Nodes[i].Lexeme
			lastIsUvar = true
			previousDepth = tree.Nodes[i].Depth
			break
		case Tokens.TokenFunction:
			if deltaDepth > 0 && lastIsUvar {
				if tree.Nodes[i-1].Token == Tokens.TokenUVar {
					// Remove prev UVar
					Output = Output[:(len(Output) - len(tree.Nodes[i-1].Lexeme))]
					// Add Brackets
					Output += brackets
					Output += tree.Nodes[i-1].Lexeme
				} else {
					// Remove prev UVar
					Output = Output[:(len(Output) - len(tree.Nodes[i-2].Lexeme))]
					// Add Brackets
					Output += brackets
					Output += tree.Nodes[i-2].Lexeme
				}
			} else {
				Output += brackets
			}
			previousDepth = tree.Nodes[i].Depth
			break
		case Tokens.Application:
			if isInLambda {
				adjustBrackets = ""
				for i := 1; i < len(brackets); i++ {
					adjustBrackets += string(brackets[i])
				}
				brackets = adjustBrackets
				if !(tree.Nodes[i-2].Token == Tokens.TokenUVar || tree.Nodes[i-2].Token == Tokens.TokenFunction) {
					// Remove prev UVar
					Output = Output[:(len(Output) - len(tree.Nodes[i-1].Lexeme))]
					Output += "("
					Output += tree.Nodes[i-1].Lexeme
					Output += ")"
				}
			}
			isInLambda = false
			Output += brackets
			previousDepth = tree.Nodes[i].Depth
		default:
			Output += brackets
			previousDepth = tree.Nodes[i].Depth
			break
		} // switch
		brackets = ""

	} // for
	// Remove redundant depth from tree.
	deltaDepth = subtreeDepth - previousDepth
	for deltaDepth < 0 {
		brackets += ")"
		deltaDepth++
	}
	Output += brackets
	return Output
}
