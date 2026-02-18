package ast

import "fmt"

type TokenStream struct {
	buf []Token
	pos int
}

func NewTokenStream(buf []Token) *TokenStream {
	return &TokenStream{
		buf: buf,
		pos: 0,
	}
}

// Peek's token, check it type and return it if true, panic if false
func (ts *TokenStream) Expect(tt TokenType) (Token, error) {
	if ts.EOF() {
		return Token{}, fmt.Errorf("TokenStream EOF")
	}
	if ts.Peek().TokenType == tt {
		return ts.Next(), nil
	}
	return Token{}, fmt.Errorf("expected token type: %v. Given: %v", tt, ts.Peek().TokenType)
}

func (ts *TokenStream) Peek() Token {
	return ts.buf[ts.pos]
}

func (ts *TokenStream) Next() Token {
	t := ts.buf[ts.pos]
	ts.pos++
	return t
}

func (ts *TokenStream) EOF() bool {
	return ts.pos >= len(ts.buf)
}
