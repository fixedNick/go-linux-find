package core

import (
	"fmt"
	"os"
)

type Action interface {
	Execute(event FileEvent)
}

type DeleteAction struct{}

func (a DeleteAction) Execute(event FileEvent) {
	var err error
	if event.fileInfo.IsDir() {
		err = os.RemoveAll(event.Path())
	} else {
		err = os.Remove(event.Path())
	}

	if err != nil {
		fmt.Println("Exec error:", err.Error())
	}
}
