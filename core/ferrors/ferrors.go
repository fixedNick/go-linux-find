package ferrors

import (
	"fmt"
)

type ParseError struct {
	Pos     int
	Lexeme  string
	Message string
}

type EvalError struct {
	Message string
}

type SemanticError struct {
	Predicate string
	Value     string
	Message   string
}

func (err SemanticError) Error() string {
	return fmt.Sprintf("[Predicate: %s / Value: %s] SemanticError: %s", err.Predicate, err.Value, err.Message)
}

func (err ParseError) Error() string {
	return fmt.Sprintf("[Lexeme: %s At: %d] ParseError: %s", err.Lexeme, err.Pos, err.Message)
}

func (err EvalError) Error() string {
	return err.Message
}
