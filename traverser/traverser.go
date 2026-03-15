package traverser

import (
	"context"
	"find/ast"
	"find/core"
	"os"
	"path/filepath"
)

type Traverser struct {
	root string
	ctx  context.Context
	size int
	ast  ast.AstNode
	stop bool
}

func New(ctx context.Context, size int, root string, astRoot ast.AstNode) *Traverser {
	return &Traverser{
		root: root,
		size: size,
		ctx:  ctx,
		ast:  astRoot,
	}
}

// Async func whitch returns channel with found files
// `Run()` finished when returned channel is closed
func (t *Traverser) Run() <-chan core.FileEvent {
	events := make(chan core.FileEvent, t.size)
	go func() {
		t.walk(filepath.Clean(t.root), events, 0)
		close(events)
	}()
	return events
}

// TODO:
// separate channels of events and eventErrors
// add option to hide/view errors
func (t *Traverser) walk(dir string, events chan<- core.FileEvent, depth int) {
	if t.stop {
		return
	}

	lstat, err := os.Lstat(dir)
	// file errors - return, stop moving deeper
	if err != nil {
		select {
		case <-t.ctx.Done():
		case events <- core.NewFileEventError(dir, depth, err):
		}
		return
	}

	// current directory
	event := core.NewFileEvent(dir, lstat, depth, lstat.Mode().Type())
	decision := t.ast.Eval(event)

	if decision.Match {
		select {
		case <-t.ctx.Done():
		case events <- event:
		}
	}

	switch decision.Control {
	case core.ControlPrune:
		return
	case core.ControlQuit:
		if event.IsDir() {
			t.stop = true
			return
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		select {
		case <-t.ctx.Done():
		case events <- core.NewFileEventError(dir, depth, err):
		}
	}

	for _, entry := range entries {
		childPath := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			select {
			case <-t.ctx.Done():
				return
			default:
				t.walk(childPath, events, depth+1)
			}

			continue
		}

		info, err := entry.Info()
		if err != nil {
			select {
			case <-t.ctx.Done():
			case events <- core.NewFileEventError(childPath, depth, err):
			}
			continue
		}

		childEvent := core.NewFileEvent(childPath, info, depth+1, info.Mode().Type())
		decision := t.ast.Eval(childEvent)
		if decision.Match {
			select {
			case <-t.ctx.Done():
			case events <- childEvent:
			}
		}

		if decision.Control == core.ControlQuit {
			t.stop = true
			return
		}

	}
}
