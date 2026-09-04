package syserrors

import "fmt"

func FormatCode(err error) string {
	v, ok := Get(err, "code")
	if !ok {
		return ""
	}

	s, ok := v.(string)
	if !ok {
		return ""
	}

	return s
}

// FormatStack returns formatted stack(s).
func FormatStack(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%+v", err)
}
