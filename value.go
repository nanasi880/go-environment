package environment

// Value is the interface to the dynamic value stored in a value.
type Value interface {
	Set(v string) error
}
