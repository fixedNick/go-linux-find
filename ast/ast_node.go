package ast

import (
	"fmt"
	"main/stuff/find/core"
	"path/filepath"
	"regexp"
	"strconv"
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

type PredicateNode struct {
	Name  string
	Value string
}

func (n PredicateNode) Eval(event core.FileEvent) core.Decision {
	switch n.Name {
	case "-name":
		m, err := regexp.Match(n.Value, []byte(filepath.Base(event.Path())))
		if err != nil {
			panic(fmt.Sprintf("REGEX ERR: %v", err))
		}
		return core.Decision{
			Match: m,
		}
	case "-depth":
		d, err := strconv.Atoi(n.Value)
		if err != nil {
			panic("depth must be INT")
		}
		return core.Decision{
			Match: d == event.Depth(),
		}
	case "-type":
		switch n.Value {
		case "f":
			return core.Decision{
				Match: event.FileType().IsRegular(),
			}
		case "d":
			return core.Decision{
				Match: event.FileType().IsDir(),
			}
		default:
			panic("unsupported file type")
		}
	}
	return core.Decision{
		Match: true,
	}
}
