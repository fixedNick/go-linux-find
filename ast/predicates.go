package ast

import (
	"fmt"
	"main/stuff/find/core"
	"maps"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Predicate struct {
	Name         string
	AllowedTypes []ValueType
	Handler      func(v Value, event core.FileEvent) bool
	Validate     func(PredicateNode) error
}

type PredicateList map[string]Predicate

var predicates = PredicateList{
	"-name": Predicate{
		Name: "-name",
		AllowedTypes: []ValueType{
			RegexType, StringType,
		},
		Handler: namePredicate,
		Validate: func(p PredicateNode) error {
			if p.Value.Regex == nil && p.Value.Str == nil {
				return fmt.Errorf("%s node: regex is nil", p.Name)
			}
			return nil
		},
	},
	"-depth": Predicate{
		Name: "-depth",
		AllowedTypes: []ValueType{
			IntType,
		},
		Handler: depthPredicate,
		Validate: func(p PredicateNode) error {
			if p.Value.Int == nil {
				return fmt.Errorf("%s node: int is nil", p.Name)
			}
			return nil
		},
	},
	"-type": Predicate{
		Name: "-type",
		AllowedTypes: []ValueType{
			StringType,
		},
		Handler: typePredicate,
		Validate: func(p PredicateNode) error {
			if _, ok := typeHandlers[p.Value.Raw]; !ok {
				var sb strings.Builder
				sb.WriteString("[ ")
				keys := maps.Keys(typeHandlers)
				for key := range keys {
					sb.WriteString(fmt.Sprintf("'%s' ", key))
				}
				sb.WriteString("]")
				return fmt.Errorf("%s %s unavailable. Use: %s", p.Name, p.Value.Raw, sb.String())
			}
			return nil
		},
	},
	"-iname":    Predicate{},
	"-path":     Predicate{},
	"-ipath":    Predicate{},
	"-size":     Predicate{},
	"-empty":    Predicate{},
	"-mindepth": Predicate{},
	"-maxdepth": Predicate{},
	"-mtime":    Predicate{},
	"-atime":    Predicate{},
	"-ctime":    Predicate{},
	"-perm":     Predicate{},
	"-user":     Predicate{},
	"-group":    Predicate{},
}

type ValueType int

const (
	StringType ValueType = iota
	IntType
	RegexType
	BoolType
)

func (p Predicate) ParseValue(raw string) (Value, []error) {

	errs := []error{}

	for _, aType := range p.AllowedTypes {

		v := Value{
			Type: aType,
			Raw:  raw,
		}

		switch aType {
		case StringType:
			v.Str = &raw
			return v, nil
		case IntType:
			i, err := strconv.Atoi(raw)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid int value: %s", raw))
				continue
			}
			v.Int = &i
			return v, nil
		case RegexType:
			r, err := regexp.Compile(raw)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid regex: %s", raw))
				continue
			}
			v.Regex = r
			return v, nil
		case BoolType:
			s, err := strconv.ParseBool(raw)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid bool: %s", raw))
				continue
			}
			v.Bool = &s
			return v, nil
		default:
			errs = append(errs, fmt.Errorf("[Predicate.ParseValue] unknown AllowedType %d", aType))
		}

	}
	return Value{}, errs
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
