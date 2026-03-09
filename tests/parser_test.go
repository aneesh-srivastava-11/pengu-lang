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
		{
			name: "V2 complete service",
			input: `version 2
service advanced
middleware logging
middleware auth jwt
health enable
metrics enable

route POST "/data"
    parse json MyReqStruct
    auth jwt
    respond 201 "ok"`,
			expected: &compiler.Service{
				Version:        "2",
				Name:           "advanced",
				Line:           1,
				Middleware:     []string{"logging", "auth jwt"},
				HealthEnabled:  true,
				MetricsEnabled: true,
				Routes: []compiler.Route{
					{
						Method: "POST",
						Path:   "/data",
						Line:   8,
						Actions: []compiler.Action{
							{Type: "parse_json", Args: []string{"MyReqStruct"}, Line: 9},
							{Type: "auth", Args: []string{"jwt"}, Line: 10},
							{Type: "respond", Args: []string{"201", "ok"}, Line: 11},
						},
					},
				},
			},
			expectErr: false,
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
