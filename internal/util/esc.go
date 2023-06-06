package util

// EscReturns changes any actual carriage returns or line returns into
// their backslashed equivalents.
func EscReturns(s string) string {
	var out string

	runes := []rune(s)
	for _, r := range runes {
		switch r {
		case '\n':
			out += "\\n"
		case '\r':
			out += "\\r"
		default:
			out += string(r)
		}
	}

	return out
}

// UnEscReturns changes any escaped carriage returns or line returns
// into their actual values.
func UnEscReturns(s string) string {
	var out string

	runes := []rune(s)
	n := len(runes)
	for i := 0; i < n; i++ {
		if runes[i] == '\\' && i != n-1 {
			if runes[i+1] == 'n' {
				out += "\n"
				i++
				continue
			}
			if runes[i+1] == 'r' {
				out += "\r"
				i++
				continue
			}
		}
		out += string(runes[i])
	}

	return out
}
