package core

import (
	"fmt"
	"maps"
	"regexp"
	"strconv"
	"strings"
)

type PredicateKind int

const (
	FilterPredicate PredicateKind = iota
	ActionPredicate
	ControlPredicate
)

type Predicate struct {
	Name         string
	AllowedTypes []ValueType
	Handler      func(v Value, event FileEvent) Decision
	Validate     func(PredicateNode) error
	Kind         PredicateKind
	Action       Action
	NoValue      bool
	Control      ControlSignal
}

type PredicateList map[string]Predicate

var Predicates = PredicateList{
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
		Kind: FilterPredicate,
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
		Kind: FilterPredicate,
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
		Kind: FilterPredicate,
	},
	"-delete": Predicate{
		Name:    "-delete",
		Kind:    ActionPredicate,
		NoValue: true,
		Action:  DeleteAction{},
	},
	"-quit": Predicate{
		Name:    "-quit",
		NoValue: true,
		Kind:    ControlPredicate,
		Control: ControlQuit,
	},
	"-mindepth": Predicate{
		Name: "-mindepth",
		AllowedTypes: []ValueType{
			IntType,
		},
		Handler: func(v Value, event FileEvent) Decision {
			if v.Int == nil {
				return Decision{Match: false}
			}
			return Decision{Match: event.Depth() >= *v.Int}
		},
		Validate: func(p PredicateNode) error {
			if p.Value.Int == nil {
				return fmt.Errorf("%s node: int is nil", p.Name)
			}
			return nil
		},
		Kind: FilterPredicate,
	},
	"-maxdepth": Predicate{
		Name: "-maxdepth",
		AllowedTypes: []ValueType{
			IntType,
		},
		Handler: func(v Value, event FileEvent) Decision {
			if v.Int == nil {
				return Decision{Match: false}
			}
			return Decision{Match: event.Depth() < *v.Int}
		},
		Validate: func(p PredicateNode) error {
			if p.Value.Int == nil {
				return fmt.Errorf("%s node: int is nil", p.Name)
			}
			return nil
		},
		Kind: FilterPredicate,
	},
	"-iname": Predicate{},
	"-path":  Predicate{},
	"-ipath": Predicate{},
	"-size":  Predicate{},
	"-empty": Predicate{},
	"-mtime": Predicate{},
	"-atime": Predicate{},
	"-ctime": Predicate{},
	"-perm":  Predicate{},
	"-user":  Predicate{},
	"-group": Predicate{},
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
		case StringType:
			v.Str = &raw
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

func namePredicate(value Value, event FileEvent) Decision {
	if value.Regex != nil {
		match := value.Regex.Match([]byte(event.Path()))
		if match {
			return Decision{Match: match}
		}
	}

	if value.Str != nil {
		return Decision{
			Match: event.Path() == *value.Str,
		}
	}
	return Decision{Match: false}
}
func depthPredicate(value Value, event FileEvent) Decision {
	if value.Int == nil {
		return Decision{Match: false}
	}
	return Decision{Match: *value.Int == event.Depth()}
}

type TypeHandler func(event FileEvent) bool

var typeHandlers = map[string]TypeHandler{
	"f": func(event FileEvent) bool { return event.FileType().IsRegular() },
	"d": func(event FileEvent) bool { return event.FileType().IsDir() },
}

func typePredicate(value Value, event FileEvent) Decision {
	if value.Str == nil {
		return Decision{
			Match: true,
		}
	}
	handler, ok := typeHandlers[*value.Str]
	if !ok {
		return Decision{
			Match: false,
		}
	}
	return Decision{
		Match: handler(event),
	}
}

func (n PredicateNode) Eval(event FileEvent) Decision {
	p, ok := Predicates[n.Name]
	if !ok {
		return Decision{
			Match: false,
		}
	}

	switch p.Kind {
	case FilterPredicate:
		decision := p.Handler(n.Value, event)
		return Decision{
			Match: decision.Match,
		}
	case ActionPredicate:
		return Decision{
			Match: true,
			Actions: []Action{
				p.Action,
			},
		}
	case ControlPredicate:
		return Decision{
			Match:   true,
			Control: p.Control,
		}
	}

	return Decision{}
}
