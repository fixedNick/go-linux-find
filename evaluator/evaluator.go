package evaluator

import (
	"context"
	"fmt"
	"main/stuff/find/ast"
	"main/stuff/find/core"
)

type Evaluator struct {
	ctx context.Context
}

func New(ctx context.Context) *Evaluator {
	return &Evaluator{
		ctx: ctx,
	}
}

func (e *Evaluator) Evaluate(node ast.AstNode, events <-chan core.FileEvent) {
	fmt.Println("Matches: ")
	for event := range events {
		d := node.Eval(event)
		if d.Match {
			fmt.Println(event.Path())
		}
	}
}
