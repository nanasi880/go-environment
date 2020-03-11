package environment

type stringValue struct {
	ptr *string
}

func newStringValue(v string) *stringValue {
	return &stringValue{
		ptr: &v,
	}
}

func (val *stringValue) Set(v string) error {
	*val.ptr = v
	return nil
}
