package ast

import "main/stuff/find/core"

type Node interface {
	Evaluate(event core.FileEvent) core.Decision
}
