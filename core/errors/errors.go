package errors

import (
	"fmt"
	"main/stuff/find/ast"
)

type ParseError struct {
	Pos     int
	Token   ast.Token
	Message string
}

type EvalError struct {
	Node    ast.AstNode
	Message string
}

func (err ParseError) Error() string {
	return fmt.Sprintf("%s at %d with %s", err.Message, err.Pos, err.Token.Lexeme)
}

func (err EvalError) Error() string {
	return err.Message
}
