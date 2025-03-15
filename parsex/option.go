package parsex

import "strings"

type Option interface {
	Id() string
	Match() []string
	Describe() string
}

func getMatches(opt Option) string {
	out := make([]string, 0, len(opt.Match()))

	for _, match := range opt.Match() {
		if len(match) < 2 {
			out = append(out, "-"+match)
		} else {
			out = append(out, "--"+match)
		}
	}

	return strings.Join(out, ", ")
}
