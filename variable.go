package environment

// A Variable represents the state of a variable.
type Variable struct {
	name  string
	value Value
	def   interface{}
	usage string
}
