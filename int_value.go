package environment

import "strconv"

type intValue struct {
	ptr *int
}

func newIntValue(v int) *intValue {
	return &intValue{
		ptr: &v,
	}
}

func (val *intValue) Set(v string) error {

	iv, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return err
	}

	*val.ptr = int(iv)
	return nil
}
