package main

import (
	"testing"
)

func TestPrepare(t *testing.T) {
	tt := []struct {
		contents string
		expectName string
		expectArgs []string
	}{
		{
			contents: "ls",
			expectName: "ls",
			expectArgs: nil,
		},
		{
			contents: " ls ",
			expectName: "ls",
			expectArgs: nil,
		},
		{
			contents: " ls -la",
			expectName: "ls",
			expectArgs: []string{"-la"},
		},
		{
			contents: `
ls \
-la
`,
            expectName: "ls",
			expectArgs: []string{"-la"},
		},
	}

	for _, tc := range tt {
		name, args := prepare(tc.contents)
		if name != tc.expectName || !stringSliceEqual(args, tc.expectArgs) {
			t.Errorf("\ngot:\n%s, %v \nexpected:\n%s, %v", name, args, tc.expectName, tc.expectArgs)
		}
	}
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}