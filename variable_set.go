package environment

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// A VariableSet represents a set of defined environment variable. The zero value of a FlagSet
// has no name and has ContinueOnError error handling.
//
// Flag names must be unique within a FlagSet. An attempt to define a flag whose
// name is already in use will cause a panic.
type VariableSet struct {
	Usage       func() string
	name        string
	variables   map[string]*Variable
	errHandling ErrorHandling
	output      io.Writer
}

// NewVariableSet returns a new, empty variable set with the specified name and
// error handling property. If the name is not empty, it will be printed
// in the default usage message and in error messages.
func NewVariableSet(name string, handling ErrorHandling) *VariableSet {
	set := &VariableSet{
		Usage:       nil,
		name:        name,
		errHandling: handling,
	}
	set.Usage = set.defaultUsage

	return set
}

// String defines a string value with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func (set *VariableSet) String(name string, def string, usage string) *string {
	sv := newStringValue(def)
	set.value(sv, def, name, usage)
	return sv.ptr
}

// Int defines a int value with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func (set *VariableSet) Int(name string, def int, usage string) *int {
	iv := newIntValue(def)
	set.value(iv, def, name, usage)
	return iv.ptr
}

// Parse parses the environment variables from argument. Must be called
// after all variables are defined and before variable are accessed by the program.
func (set *VariableSet) Parse(environments []string) error {

	keys := set.sortedKeys()
	for _, key := range keys {
		variable := set.variables[key]

		osVal, found := set.lookup(environments, key)
		if !found {
			continue
		}

		err := variable.value.Set(osVal)
		if err == nil {
			continue
		}

		e := fmt.Errorf("parse error %s: %w", key, err)
		switch set.errHandling {

		case ContinueOnError:
			return e

		case ExitOnError:
			w := set.output
			if w == nil {
				w = os.Stderr
			}
			set.fprintf(w, "%v\n", e)
			os.Exit(2)

		default:
			fallthrough
		case PanicOnError:
			panic(e)
		}
	}

	return nil
}

func (set *VariableSet) value(v Value, def interface{}, name string, usage string) {

	if set.variables == nil {
		set.variables = make(map[string]*Variable)
	}

	_, ok := set.variables[name]
	if ok {
		msg := fmt.Sprintf("%s is already exist", name)
		panic(msg)
	}

	variable := &Variable{
		name:  name,
		value: v,
		def:   def,
		usage: usage,
	}
	set.variables[name] = variable
}

func (set *VariableSet) defaultUsage() string {
	var sb strings.Builder

	if len(set.name) > 0 {
		set.fprintf(&sb, "Usage of %s\n", set.name)
	}

	keys := set.sortedKeys()
	for _, key := range keys {
		variable := set.variables[key]

		if variable.def != nil {
			set.fprintf(&sb, "  %s: default(%v)\n", variable.name, variable.def)
		} else {
			set.fprintf(&sb, "  %s:\n", variable.name)
		}

		for _, line := range strings.Split(variable.usage, "\n") {
			set.fprintf(&sb, "    %s\n", line)
		}
	}

	return sb.String()
}

func (_ *VariableSet) fprintf(w io.Writer, format string, args ...interface{}) {
	_, _ = fmt.Fprintf(w, format, args...)
}

func (_ *VariableSet) lookup(environments []string, key string) (string, bool) {

	for _, pair := range environments {

		if !strings.HasPrefix(pair, key) {
			continue
		}

		for i := 0; i < len(pair); i++ {
			if pair[i] == '=' {
				return pair[i+1:], true
			}
		}

		return "", false
	}

	return "", false
}

func (set *VariableSet) sortedKeys() []string {
	var keys []string
	for key := range set.variables {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
