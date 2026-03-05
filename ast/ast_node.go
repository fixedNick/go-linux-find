package ast

import (
	"find/core"
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
	switch n.Op {
	case TOKEN_LOGICAL_AND:
		left := n.Left.Eval(event)
		if !left.Match {
			return left
		}
		right := n.Right.Eval(event)
		return core.Decision{
			Match:   left.Match && right.Match,
			Actions: append(left.Actions, right.Actions...),
			Control: core.MergeControl(left.Control, right.Control),
		}
	case TOKEN_LOGICAL_OR:
		left := n.Left.Eval(event)
		if left.Match {
			return left
		}
		right := n.Right.Eval(event)
		return right
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
		res := n.Node.Eval(event)
		return core.Decision{
			Match:   !res.Match,
			Actions: res.Actions,
			Control: res.Control,
		}
	}
	panic("unexpected Operation in UnaryNode")
}
