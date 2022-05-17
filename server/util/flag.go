package util

type Flag uint16

func (r Flag) Has(flag Flag) bool {
	return r&flag != 0
}

func (r *Flag) Append(flag Flag) {
	*r |= flag
}

func (r *Flag) Remove(flag Flag) {
	*r &= ^flag
}
