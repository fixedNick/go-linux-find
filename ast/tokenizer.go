package ast

import (
	"find/core"
	"fmt"
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

var tokenNames = map[TokenType]string{
	TOKEN_VALUE:       "TOKEN_VALUE",
	TOKEN_LOGICAL_AND: "TOKEN_LOGICAL_AND",
	TOKEN_LOGICAL_OR:  "TOKEN_LOGICAL_OR",
	TOKEN_LOGICAL_NOT: "TOKEN_LOGICAL_NOT",
	TOKEN_LPAREN:      "TOKEN_LPAREN",
	TOKEN_RPAREN:      "TOKEN_RPAREN",
	TOKEN_PREDICATE:   "TOKEN_PREDICATE",
}

func (tt TokenType) String() string {
	if name, ok := tokenNames[tt]; ok {
		return name
	}
	return fmt.Sprintf("TokenType(%d)", int(tt))
}

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

func (t Tokenizer) isPredicate(p string) bool {
	_, ok := core.Predicates[p]
	return ok
}
