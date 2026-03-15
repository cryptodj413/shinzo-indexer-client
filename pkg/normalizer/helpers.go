package normalizer

// hasHexPrefix checks if a string has a valid hex prefix (0x or 0X)
func hasHexPrefix(s string) bool {
	return len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}
