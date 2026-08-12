package analyzer

import "testing"

func TestShouldSkipDir(t *testing.T) {

	cases := map[string]bool{
		".git":          true,
		".gradle":       true,
		".idea":         true,
		"build":         true,
		"node_modules":  true,
		"target":        true,
		".infra-doctor": true,
		"src":           false,
		".github":       false,
		"main":          false,
	}

	for name, want := range cases {
		if got := shouldSkipDir(name); got != want {
			t.Errorf("shouldSkipDir(%q) = %v, want %v", name, got, want)
		}
	}
}
