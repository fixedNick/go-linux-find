package ast

import (
	"find/core"
	"find/core/ferrors"
	"fmt"
)

type TokenStreamer interface {
	Peek() Token
	Next() Token
	Pos() int
	Expect(TokenType) (Token, error)
	EOF() bool
}

type Parser struct {
	stream    TokenStreamer
	validator *ASTValidator
}

func NewParser(stream TokenStreamer, validator *ASTValidator) *Parser {
	return &Parser{
		stream: stream,
	}
}

func (p *Parser) Parse() (AstNode, []error) {
	root, err := p.parseExpression()
	if err != nil {
		return nil, []error{err}
	}

	verr := p.validator.Validate(root)
	return root, verr
}
func (p *Parser) parseExpression() (AstNode, error) {
	return p.parseOr()
}
func (p *Parser) parseOr() (AstNode, error) {

	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for !p.stream.EOF() && p.stream.Peek().TokenType == TOKEN_LOGICAL_OR {
		p.stream.Next()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = BinaryNode{
			Op:    TOKEN_LOGICAL_OR,
			Left:  left,
			Right: right,
		}
	}

	return left, nil

}

// find . -type f -name g -o -type d -name b
// and_expr    := unary_expr ( AND unary_expr | implicit_and )*
func (p *Parser) parseAnd() (AstNode, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for !p.stream.EOF() {
		t := p.stream.Peek()

		if t.TokenType == TOKEN_LOGICAL_AND {
			p.stream.Next()
			right, err := p.parseUnary()
			if err != nil {
				return nil, err
			}
			left = BinaryNode{
				Op:    TOKEN_LOGICAL_AND,
				Left:  left,
				Right: right,
			}
			continue
		}

		if t.TokenType == TOKEN_PREDICATE || t.TokenType == TOKEN_LPAREN || t.TokenType == TOKEN_LOGICAL_NOT {
			right, err := p.parseUnary()
			if err != nil {
				return nil, err
			}
			left = BinaryNode{
				Op:    TOKEN_LOGICAL_AND,
				Left:  left,
				Right: right,
			}
			continue
		}

		break
	}
	return left, nil
}

// unary_expr  := NOT unary_expr | primary
func (p *Parser) parseUnary() (AstNode, error) {
	// possible: predicate, group, not
	if p.stream.Peek().TokenType == TOKEN_LOGICAL_NOT {
		p.stream.Next()
		un, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		return UnaryNode{
			Op:   TOKEN_LOGICAL_NOT,
			Node: un,
		}, nil
	}
	return p.parsePrimary()
}

// primary     := predicate | group
func (p *Parser) parsePrimary() (AstNode, error) {
	if p.stream.Peek().TokenType == TOKEN_LPAREN {
		return p.parseGroup()
	}
	return p.parsePredicate()
}

// predicate value
func (p *Parser) parsePredicate() (AstNode, error) {
	pred, err := p.stream.Expect(TOKEN_PREDICATE)
	if err != nil {
		return nil, ferrors.ParseError{
			Pos:     p.stream.Pos(),
			Lexeme:  p.stream.Peek().Lexeme,
			Message: err.Error(),
		}
	}
	val, err := p.stream.Expect(TOKEN_VALUE)
	if err != nil {
		return nil, ferrors.ParseError{
			Pos:     p.stream.Pos(),
			Lexeme:  p.stream.Peek().Lexeme,
			Message: err.Error(),
		}
	}

	predicate, ok := core.Predicates[pred.Lexeme]
	if !ok {
		return nil, ferrors.ParseError{
			Pos:     p.stream.Pos(),
			Lexeme:  p.stream.Peek().Lexeme,
			Message: fmt.Sprintf("unknown predicate type of lexeme: %s", pred.Lexeme),
		}
	}

	v, errs := predicate.ParseValue(val.Value)
	if len(errs) > 0 {
		return nil, ferrors.SemanticError{
			Predicate: pred.Lexeme,
			Value:     val.Lexeme,
			Message:   errs[len(errs)-1].Error(),
		}
	}

	return core.PredicateNode{
		Name:  pred.Lexeme,
		Value: v,
	}, nil
}

// group       := '(' expression ')'
func (p *Parser) parseGroup() (AstNode, error) {
	_, err := p.stream.Expect(TOKEN_LPAREN)
	if err != nil {
		return nil, ferrors.ParseError{
			Pos:     p.stream.Pos(),
			Lexeme:  p.stream.Peek().Lexeme,
			Message: err.Error(),
		}
	}

	node, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.stream.Expect(TOKEN_RPAREN)
	if err != nil {
		return nil, ferrors.ParseError{
			Pos:     p.stream.Pos(),
			Lexeme:  p.stream.Peek().Lexeme,
			Message: err.Error(),
		}
	}
	return node, nil
}
