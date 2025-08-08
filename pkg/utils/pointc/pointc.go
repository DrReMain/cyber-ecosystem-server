package pointc

func P[T comparable](in T) *T {
	return &in
}

func PStatus32t8(value *uint32) *uint8 {
	if value == nil {
		return nil
	}
	result := new(uint8)
	*result = uint8(*value)
	return result
}

func PStatus8t32(value *uint8) *uint32 {
	if value == nil {
		return nil
	}
	result := new(uint32)
	*result = uint32(*value)
	return result
}
