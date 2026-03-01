package ast

import (
	"main/stuff/find/core"
)

type AstNode interface {
	Eval(event core.FileEvent) core.Decision
}

type BinaryNode struct {
	Op    TokenType
	Left  AstNode
	Right AstNode
}

func (n BinaryNode) Eval(event core.FileEvent) core.Decision {
	l := n.Left.Eval(event).Match
	r := n.Right.Eval(event).Match

	switch n.Op {
	case TOKEN_LOGICAL_AND:
		return core.Decision{
			Match: l && r,
		}
	case TOKEN_LOGICAL_OR:
		return core.Decision{
			Match: l || r,
		}
	}
	panic("unexpected Operation in BinaryNode")
}

type UnaryNode struct {
	Op   TokenType
	Node AstNode
}

func (n UnaryNode) Eval(event core.FileEvent) core.Decision {
	switch n.Op {
	case TOKEN_LOGICAL_NOT:
		return core.Decision{
			Match: !n.Node.Eval(event).Match,
		}
	}
	panic("unexpected Operation in UnaryNode")
}

func (n PredicateNode) Eval(event core.FileEvent) core.Decision {
	p, ok := predicates[n.Name]
	if !ok {
		return core.Decision{
			Match: false,
		}
	}
	return core.Decision{
		Match: p.Handler(n.Value, event),
	}
}
