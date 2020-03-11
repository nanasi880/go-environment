package environment

import (
	"os"
)

// Default is the default set of environment variable set, parsed from os.Environment().
var Default = NewVariableSet("", ExitOnError)

// Usage returns Default.Usage().
var Usage = func() string {
	return Default.Usage()
}

// String defines a string value with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func String(name string, def string, usage string) *string {
	return Default.String(name, def, usage)
}

// Int defines a int value with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func Int(name string, def int, usage string) *int {
	return Default.Int(name, def, usage)
}

// Parse parses the environment variables from os.Environ(). Must be called
// after all variables are defined and before variable are accessed by the program.
func Parse() {
	_ = Default.Parse(os.Environ())
}
