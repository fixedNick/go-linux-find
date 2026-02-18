package ast

import (
	"slices"
)

type Token struct {
	TokenType
	Lexeme string
	Value  string
}

type TokenType int

const (
	TOKEN_VALUE       TokenType = 1
	TOKEN_LOGICAL_AND TokenType = 2
	TOKEN_LOGICAL_OR  TokenType = 3
	TOKEN_LOGICAL_NOT TokenType = 4
	TOKEN_LPAREN      TokenType = 5
	TOKEN_RPAREN      TokenType = 6
	TOKEN_PREDICATE   TokenType = 7
)

type Tokenizer struct{}

func (t *Tokenizer) Tokenize(args []string) *TokenStream {
	result := make([]Token, 0, len(args))

	for _, arg := range args {

		if t.isPredicate(arg) {
			result = append(result, Token{
				TokenType: TOKEN_PREDICATE,
				Lexeme:    arg,
			})
			continue
		}

		if tokenType, ok := t.isLogical(arg); ok {
			result = append(result, Token{
				TokenType: tokenType,
				Lexeme:    arg,
			})
			continue
		}

		if arg == "(" {
			result = append(result, Token{
				TokenType: TOKEN_LPAREN,
			})
			continue
		}

		if arg == ")" {
			result = append(result, Token{
				TokenType: TOKEN_RPAREN,
			})
			continue
		}

		result = append(result, Token{
			TokenType: TOKEN_VALUE,
			Value:     arg,
		})
	}

	return NewTokenStream(result)
}

var logicalOps = map[string]TokenType{
	"-a": TOKEN_LOGICAL_AND,
	"-o": TOKEN_LOGICAL_OR,
	"!":  TOKEN_LOGICAL_NOT,
}

func (t Tokenizer) isLogical(p string) (TokenType, bool) {

	logical, ok := logicalOps[p]
	return logical, ok
}

var predicates = []string{"-name", "-depth", "-type"}

func (t Tokenizer) isPredicate(p string) bool {
	return slices.Contains(predicates, p)
}
