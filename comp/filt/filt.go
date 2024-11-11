package filt

import "strings"

type Prefix string

func (prefix Prefix) Filter(s []string) []string {
	out := make([]string, 0)
	for _, v := range s {
		if strings.HasPrefix(v, string(prefix)) {
			out = append(out, v)
		}
	}
	return out
}

type MaxLen int

func (maxLen MaxLen) Filter(s []string) []string {
	out := make([]string, 0)
	for _, v := range s {
		if len(v) <= int(maxLen) {
			out = append(out, v)
		}
	}
	return out
}

type MinLen int

func (minLen MinLen) Filter(s []string) []string {
	out := make([]string, 0)
	for _, v := range s {
		if len(v) >= int(minLen) {
			out = append(out, v)
		}
	}
	return out
}
