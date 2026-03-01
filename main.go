package main

import (
	"context"
	"find/ast"
	"find/evaluator"
	"find/traverser"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: find <path> [predicates]")
		return
	}

	ctx := context.Background()

	basePath := os.Args[1]
	predicates := os.Args[2:]

	tokenizer := ast.Tokenizer{}
	token_stream := tokenizer.Tokenize(predicates)
	parser := ast.NewParser(token_stream, &ast.ASTValidator{})
	root, errors := parser.Parse()
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err.Error())
		}
		return
	}

	// create ExecutionPlan
	// pass to basePath
	traverser := traverser.New(ctx, 100, basePath)
	events := traverser.Run()

	evaluator := evaluator.New(ctx)
	evaluator.Evaluate(root, events)
}
