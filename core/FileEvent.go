package core

import (
	"io/fs"
	"os"
)

type FileEvent struct {
	path     string
	fileInfo os.FileInfo
	depth    int
	ftype    fs.FileMode
	err      error
}

func NewFileEvent(path string, fileInfo os.FileInfo, depth int, t fs.FileMode) FileEvent {
	return FileEvent{
		path:     path,
		fileInfo: fileInfo,
		depth:    depth,
		ftype:    t,
	}
}
func NewFileEventError(path string, depth int, err error) FileEvent {
	return FileEvent{
		path:  path,
		depth: depth,
		err:   err,
	}
}

// getters
func (f FileEvent) Path() string          { return f.path }
func (f FileEvent) FileInfo() os.FileInfo { return f.fileInfo }
func (f FileEvent) Depth() int            { return f.depth }
func (f FileEvent) FileType() fs.FileMode { return f.ftype }
func (f FileEvent) Err() error            { return f.err }
