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
	err := os.Remove(event.Path())
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}
