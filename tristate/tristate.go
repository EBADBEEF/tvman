package tristate

type Value int

func (f *Value) Set(val bool) {
	if val {
		*f = 2
	} else {
		*f = 1
	}
}

func (f *Value) Unset() {
	*f = 0
}

func (f *Value) IsOn() bool {
	return *f == 2
}
func (f *Value) IsOff() bool {
	return *f == 1
}
func (f *Value) IsUnset() bool {
	return *f == 0
}
