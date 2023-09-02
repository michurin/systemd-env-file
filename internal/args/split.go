package args

import "strings"

func Split(s string) []string {
	if len(s) == 0 {
		return nil
	}
	return strings.Split(s, ":")
}
