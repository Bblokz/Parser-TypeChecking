### Parser & Type Checker Commandline Application

This project contains a lexical analyser, parser and type checker for the following
Backus-Naur grammar:

```diff
- {judgement} ::= {expr} ':' {type}
- {expr} ::= {lvar} | '(' {expr} ')' | 'λ' {lvar} '^' {type} {expr} | {type} {expr}
- {type} ::= {uvar} | '(' {type} ')' | {type} '->' {type}
```
where {lvar} stands for any variable name that starts with a lowercase letter,
and {uvar} stands for any variable name that starts with an uppercase letter. A
variable name is alphanumerical: it consists of the letters a-z, A-Z, or the digits
0-9. The grammar is whitespace insensitive, but a whitespace is recognized to separate application of two variables.
The program supports international variable names.

#### What makes this Parser unique
Rather than the common practise of constructing the associated parse-tree of the analysed lines at hand afther the entire expression
was parsed, 
this project aims to construct the tree dynamically at parse-runtime. This means that 
passed on tokens are placed in their hierarchy directly with the information that is available at the time. When parsing
continues it might be needed to move tokens to new positions, depending on their context. For this purpose an
inject node method is used which carefully redistributes tokens among the parse-tree hierarchy in line with the grammar rules.

#### Setup
1) clone the repo and make sure your GOPATH environment variable can access the go.mod file.
2) Build main.go and run the executable

#### Runtime
To test the application with various expressions add lambda calculus expressions in the 
data.txt file. The lexical analyser is line-break sensitive, a new line will result in
a new expression with its own parse-tree.



