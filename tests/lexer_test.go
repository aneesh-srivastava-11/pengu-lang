package tests

import (
	"pengu-lang/compiler"
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []compiler.Token
	}{
		{
			name:  "Basic service",
			input: "version 1\nservice my_app",
			expected: []compiler.Token{
				{Type: compiler.TokenVersion, Literal: "version", Line: 1, Indent: 0},
				{Type: compiler.TokenNumber, Literal: "1", Line: 1, Indent: 0},
				{Type: compiler.TokenService, Literal: "service", Line: 2, Indent: 0},
				{Type: compiler.TokenIdent, Literal: "my_app", Line: 2, Indent: 0},
				{Type: compiler.TokenEOF, Literal: "", Line: 3, Indent: 0},
			},
		},
		{
			name: "Indented route actions",
			input: `route GET "/test"
    log "hello"`,
			expected: []compiler.Token{
				{Type: compiler.TokenRoute, Literal: "route", Line: 1, Indent: 0},
				{Type: compiler.TokenIdent, Literal: "GET", Line: 1, Indent: 0},
				{Type: compiler.TokenString, Literal: "/test", Line: 1, Indent: 0},
				{Type: compiler.TokenLog, Literal: "log", Line: 2, Indent: 4},
				{Type: compiler.TokenString, Literal: "hello", Line: 2, Indent: 4},
				{Type: compiler.TokenEOF, Literal: "", Line: 3, Indent: 0},
			},
		},
		{
			name: "Line comments ignored",
			input: `// this is a comment
version 1`,
			expected: []compiler.Token{
				{Type: compiler.TokenVersion, Literal: "version", Line: 2, Indent: 0},
				{Type: compiler.TokenNumber, Literal: "1", Line: 2, Indent: 0},
				{Type: compiler.TokenEOF, Literal: "", Line: 3, Indent: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compiler.Tokenize(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Tokenize() =\n%v\nwant\n%v", got, tt.expected)
			}
		})
	}
}
