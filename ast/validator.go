package ast

import (
	"find/core/ferrors"
	"fmt"
)

type ASTValidator struct{}

func (v *ASTValidator) Validate(root AstNode) []error {
	astErrors := []error{}
	predNodes := []PredicateNode{}

	// tree validator
	v.walk(root, &predNodes, &astErrors)

	// validate predicates conflicts
	v.conflicts(predNodes, &astErrors)

	return astErrors
}

func (v *ASTValidator) walk(root AstNode, predNodes *[]PredicateNode, astErrors *[]error) {
	switch r := root.(type) {
	case BinaryNode:
		if r.Left == nil {
			*astErrors = append(*astErrors, fmt.Errorf("left node of %s is nil", r.Op.String()))
		} else {
			v.walk(r.Left, predNodes, astErrors)
		}

		if r.Right == nil {
			*astErrors = append(*astErrors, fmt.Errorf("left node of %s is nil", r.Op.String()))
		} else {
			v.walk(r.Right, predNodes, astErrors)
		}
	case UnaryNode:
		if r.Node == nil {
			*astErrors = append(*astErrors, fmt.Errorf("unary node of %s is nil", r.Op.String()))
			return
		}
		v.walk(r.Node, predNodes, astErrors)
	case PredicateNode:
		v.validatePredicate(r, predNodes, astErrors)
	}
}

func (v *ASTValidator) validatePredicate(p PredicateNode, predNodes *[]PredicateNode, astErrors *[]error) {
	predicate, ok := predicates[p.Name]
	if !ok {
		*astErrors = append(*astErrors, ferrors.SemanticError{
			Predicate: p.Name,
			Value:     p.Value.Raw,
			Message:   "unknown predicate",
		})
		return
	}

	if err := predicate.Validate(p); err != nil {
		*astErrors = append(*astErrors, err)
		return
	}

	// all good, predicate is active
	*predNodes = append(*predNodes, p)
}

func (v *ASTValidator) conflicts(predNodes []PredicateNode, astErrors *[]error) {
	// panic("not implemented")
}
