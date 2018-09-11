package main

import (
	"strings"
	"testing"
)

func TestLowerFunction(t *testing.T) {
	testStrings := map[string]string{
		"CamelCase":   "camel_case",
		"camelCase":   "camel_case",
		"__camelCase": "__camel_case",
		"Camel Case":  "camel_case",
		"ipAddr":      "ip_addr",
		"ipAddr1":     "ip_addr1",
		"IPAddr":      "ip_addr",
	}

	for k, v := range testStrings {
		r := toLower(k)
		if r != v {
			t.Error("test", k, "expected", v, "got", r)
		}
	}
}

func BenchmarkToLower(b *testing.B) {
	for i := 0; i < b.N; i++ {
		toLower("Bernardo O'Higgins")
	}
}

func BenchmarkToLowerRef(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.ToLower("Bernardo O'Higgins")
	}
}

func toLower(v string) string {
	s := make([]byte, 0)

	for i := 0; i < len(v); i++ {
		// 32 = space
		// from 49 to 57 digits
		// from 65 to 90 uppercase A..Z
		// from 97 to 122 lowercase a..z
		if v[i] < 49 || v[i] > 122 || v[i] == 32 {
			continue
		}

		// catch uppercase chars
		if v[i] >= 65 && v[i] <= 90 {
			// if current char isn't are first, prepend it with `_`
			// also look behind to next char, if it small prepend `_`
			// if it uppercase just lower it
			// it's needed for IPAddress to convert ip_address
			if i >= 1 && v[i+1] >= 97 {
				s = append(s, 95)
			}
			s = append(s, v[i]+32)
			continue
		}
		// just copy lowercase chars
		s = append(s, v[i])
	}

	return string(s)
}
