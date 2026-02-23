package ast

import (
	"fmt"
	"main/stuff/find/core"
	"path/filepath"
	"regexp"
	"strconv"
)

type PredicateHanlder func(value Value, event core.FileEvent) bool

var expectTypes = map[string]ValueType{
	"-name":  RegexType,
	"-depth": IntType,
	"-type":  StringType,
}
var predicates = map[string]PredicateHanlder{
	"-name":  namePredicate,
	"-depth": depthPredicate,
	"-type":  typePredicate,
}

type ValueType int

const (
	StringType ValueType = iota
	IntType
	RegexType
	BoolType
)

func (vt ValueType) ParseValue(val string) (interface{}, error) {
	if vt == StringType {
		return val, nil
	}
	if vt == IntType {
		return strconv.Atoi(val)
	}
	if vt == RegexType {
		return regexp.Compile(val)
	}
	if vt == BoolType {
		return strconv.ParseBool(val)
	}
	return nil, fmt.Errorf("[ValueType.ParseValue] unexpected value type")
}

type Value struct {
	Type   ValueType
	Raw    string
	Parsed interface{}
}

type PredicateNode struct {
	Name  string
	Value Value
}

func namePredicate(value Value, event core.FileEvent) bool {
	return value.Parsed.(*regexp.Regexp).Match([]byte(filepath.Base(event.Path())))
}
func depthPredicate(value Value, event core.FileEvent) bool {
	return value.Parsed == event.Depth()
}

type TypeHandler func(event core.FileEvent) bool

var typeHandlers = map[string]TypeHandler{
	"f": func(event core.FileEvent) bool { return event.FileType().IsRegular() },
	"d": func(event core.FileEvent) bool { return event.FileType().IsDir() },
}

func typePredicate(value Value, event core.FileEvent) bool {
	handler, ok := typeHandlers[value.Parsed.(string)]
	if !ok {
		panic("unsupported type of file")
	}
	return handler(event)
}
