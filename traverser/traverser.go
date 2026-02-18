package traverser

import (
	"context"
	"main/stuff/find/core"
	"os"
	"path/filepath"
)

type Traverser struct {
	root string
	ctx  context.Context
	size int
}

func New(ctx context.Context, size int, root string) *Traverser {
	return &Traverser{
		root: root,
		size: size,
		ctx:  ctx,
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

func (t *Traverser) walk(dir string, events chan<- core.FileEvent, depth int) {

	lstat, err := os.Stat(dir)
	// file errors - return, stop moving deeper
	if err != nil {
		select {
		case <-t.ctx.Done():
		case events <- core.NewFileEventError(dir, depth, err):
		}
		return
	}

	// adding current directory
	select {
	case <-t.ctx.Done():
	case events <- core.NewFileEvent(dir, lstat, depth, lstat.Mode().Type()):
	}

	entries, err := os.ReadDir(dir)
	// file errors - return, stop moving deeper
	if err != nil {
		select {
		case <-t.ctx.Done():
		case events <- core.NewFileEventError(dir, depth, err):
		}
		return
	}

	for _, entry := range entries {
		childPath := filepath.Join(dir, entry.Name())

		if !entry.IsDir() {

			info, err := entry.Info()
			if err != nil {
				select {
				case <-t.ctx.Done():
				case events <- core.NewFileEventError(childPath, depth, err):
				}
				continue
			}

			select {
			case <-t.ctx.Done():
			case events <- core.NewFileEvent(childPath, info, depth, info.Mode().Type()):
			}

			continue
		}

		select {
		case <-t.ctx.Done():
			return
		default:
			t.walk(childPath, events, depth+1)
		}
	}
}
