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

func (vt ValueType) ParseValue(raw string) (Value, error) {

	v := Value{
		Type: vt,
		Raw:  raw,
	}

	switch vt {
	case StringType:
		v.Str = &raw
		return v, nil
	case IntType:
		i, err := strconv.Atoi(raw)
		if err != nil {
			return v, fmt.Errorf("invalid int value: %s", raw)
		}
		v.Int = &i
		return v, nil
	case RegexType:
		r, err := regexp.Compile(raw)
		if err != nil {
			return v, fmt.Errorf("invalid regex: %s", raw)
		}
		v.Regex = r
		return v, nil
	case BoolType:
		s, err := strconv.ParseBool(raw)
		if err != nil {
			return v, fmt.Errorf("invalid bool: %s", raw)
		}
		v.Bool = &s
		return v, nil
	}

	return v, fmt.Errorf("[ValueType.ParseValue] unknown ValueType")
}

type Value struct {
	Type  ValueType
	Raw   string
	Int   *int
	Str   *string
	Regex *regexp.Regexp
	Bool  *bool
}

type PredicateNode struct {
	Name  string
	Value Value
}

func namePredicate(value Value, event core.FileEvent) bool {
	if value.Regex == nil {
		return false
	}
	return value.Regex.Match([]byte(filepath.Base(event.Path())))
}
func depthPredicate(value Value, event core.FileEvent) bool {
	if value.Int == nil {
		return false
	}
	return *value.Int == event.Depth()
}

type TypeHandler func(event core.FileEvent) bool

var typeHandlers = map[string]TypeHandler{
	"f": func(event core.FileEvent) bool { return event.FileType().IsRegular() },
	"d": func(event core.FileEvent) bool { return event.FileType().IsDir() },
}

func typePredicate(value Value, event core.FileEvent) bool {
	if value.Str == nil {
		return false
	}
	handler, ok := typeHandlers[*value.Str]
	if !ok {
		return false
	}
	return handler(event)
}
