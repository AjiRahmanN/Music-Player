package repositories

import "testing"

func TestEscapeLike(t *testing.T) {
	cases := map[string]string{
		"abc":    "abc",
		"100%":   `100\%`,
		"a_b":    `a\_b`,
		`a\b`:    `a\\b`,
		`%_\`:    `\%\_\\`,
		"":       "",
		"lagu A": "lagu A",
	}
	for in, want := range cases {
		if got := EscapeLike(in); got != want {
			t.Errorf("EscapeLike(%q) = %q, want %q", in, got, want)
		}
	}
}
