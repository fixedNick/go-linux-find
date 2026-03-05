package main

import (
	"context"
	"find/ast"
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

	// tokenize received command
	// parse from tokenStream to AST
	// ! create traverser and put AST into it
	// ! receive in&out channels from traverser.Run

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
	traverser := traverser.New(ctx, 100, basePath, root)
	events := traverser.Run()

	for event := range events {
		if event.Err() != nil {
			continue
		}
		for _, a := range root.Eval(event).Actions {
			a.Execute(event)
		}
		fmt.Printf("[DIR: %v | DEPTH: %d] %s\n", event.FileInfo().IsDir(), event.Depth(), event.Path())
	}

}
