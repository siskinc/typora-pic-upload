package strx

// IsAlpha is alphabet
func IsAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// IsNumber is number
func IsNumber(c byte) bool {
	return c >= '0' && c <= '9'
}

// IsNumberOrAlpha
func IsNumberOrAlpha(c byte) bool {
	return IsAlpha(c) || IsNumber(c)
}
