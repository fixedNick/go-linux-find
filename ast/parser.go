package ast

type TokenStreamer interface {
	Peek() Token
	Next() Token
	Expect(TokenType) (Token, error)
	EOF() bool
}

type Parser struct {
	stream TokenStreamer
}

func NewParser(stream TokenStreamer) *Parser {
	return &Parser{
		stream: stream,
	}
}

func (p *Parser) Parse() AstNode {
	return p.parseExpression()
}
func (p *Parser) parseExpression() AstNode {
	return p.parseOr()
}
func (p *Parser) parseOr() AstNode {

	left := p.parseAnd()
	for !p.stream.EOF() && p.stream.Peek().TokenType == TOKEN_LOGICAL_OR {
		p.stream.Next()
		right := p.parseAnd()
		left = BinaryNode{
			Op:    TOKEN_LOGICAL_OR,
			Left:  left,
			Right: right,
		}
	}

	return left

}

// find . -type f -name g -o -type d -name b
// and_expr    := unary_expr ( AND unary_expr | implicit_and )*
func (p *Parser) parseAnd() AstNode {
	left := p.parseUnary()
	for !p.stream.EOF() {
		t := p.stream.Peek()

		if t.TokenType == TOKEN_LOGICAL_AND {
			p.stream.Next()
			right := p.parseUnary()
			left = BinaryNode{
				Op:    TOKEN_LOGICAL_AND,
				Left:  left,
				Right: right,
			}
			continue
		}

		if t.TokenType == TOKEN_PREDICATE || t.TokenType == TOKEN_LPAREN || t.TokenType == TOKEN_LOGICAL_NOT {
			right := p.parseUnary()
			left = BinaryNode{
				Op:    TOKEN_LOGICAL_AND,
				Left:  left,
				Right: right,
			}
			continue
		}

		break
	}
	return left
}

// unary_expr  := NOT unary_expr | primary
func (p *Parser) parseUnary() AstNode {
	// possible: predicate, group, not
	if p.stream.Peek().TokenType == TOKEN_LOGICAL_NOT {
		p.stream.Next()
		return UnaryNode{
			Op:   TOKEN_LOGICAL_NOT,
			Node: p.parseUnary(),
		}
	}
	return p.parsePrimary()
}

// primary     := predicate | group
func (p *Parser) parsePrimary() AstNode {
	if p.stream.Peek().TokenType == TOKEN_LPAREN {
		return p.parseGroup()
	}
	return p.parsePredicate()
}

// predicate value
func (p *Parser) parsePredicate() AstNode {
	pred, err := p.stream.Expect(TOKEN_PREDICATE)
	if err != nil {
		panic(err)
	}
	val, err := p.stream.Expect(TOKEN_VALUE)
	if err != nil {
		panic(err)
	}
	return PredicateNode{
		Name:  pred.Lexeme,
		Value: val.Value,
	}
}

// group       := '(' expression ')'
func (p *Parser) parseGroup() AstNode {
	p.stream.Expect(TOKEN_LPAREN)
	node := p.parseExpression()
	p.stream.Expect(TOKEN_RPAREN)
	return node
}
