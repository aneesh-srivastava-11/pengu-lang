package tests

import (
	"pengu-lang/compiler"
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    *compiler.Service
		expectErr   bool
		errContains string
	}{
		{
			name: "Valid simple service",
			input: `version 1
service test
route GET "/hi"
    respond 200 "hello"`,
			expected: &compiler.Service{
				Version: "1",
				Name:    "test",
				Line:    1,
				Routes: []compiler.Route{
					{
						Method: "GET",
						Path:   "/hi",
						Line:   3,
						Actions: []compiler.Action{
							{Type: "respond", Args: []string{"200", "hello"}, Line: 4},
						},
					},
				},
			},
			expectErr: false,
		},
		{
			name:        "Missing version",
			input:       `service test`,
			expected:    nil,
			expectErr:   true,
			errContains: "file must start with version declaration",
		},
		{
			name: "Unindented action",
			input: `version 1
service test
route GET "/hi"
respond 200 "hello"`, // no indent
			expected:    nil,
			expectErr:   true,
			errContains: "actions inside routes must be indented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := compiler.Tokenize(tt.input)
			parser := compiler.NewParser(tokens)
			got, err := parser.Parse()

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Parse() =\n%+v\nwant\n%+v", got, tt.expected)
			}
		})
	}
}
