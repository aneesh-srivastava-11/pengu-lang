package compiler

import (
	"strings"
	"unicode"
)

type TokenType string

const (
	TokenEOF     TokenType = "EOF"
	TokenIdent   TokenType = "IDENT"
	TokenString  TokenType = "STRING"
	TokenNumber  TokenType = "NUMBER"

	TokenVersion TokenType = "VERSION"
	TokenService TokenType = "SERVICE"
	TokenRoute   TokenType = "ROUTE"
	TokenLog     TokenType = "LOG"
	TokenRespond TokenType = "RESPOND"
)

var keywords = map[string]TokenType{
	"version": TokenVersion,
	"service": TokenService,
	"route":   TokenRoute,
	"log":     TokenLog,
	"respond": TokenRespond,
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Indent  int
}

func Tokenize(input string) []Token {
	var tokens []Token
	lines := strings.Split(input, "\n")

	for i, line := range lines {
		lineNum := i + 1

		indent := 0
		for _, ch := range line {
			if ch == ' ' {
				indent++
			} else if ch == '\t' {
				indent += 4
			} else {
				break
			}
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}

		chars := []rune(trimmed)
		pos := 0

		for pos < len(chars) {
			ch := chars[pos]

			if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
				pos++
				continue
			}

			if ch == '"' {
				pos++
				start := pos
				for pos < len(chars) && chars[pos] != '"' {
					pos++
				}
				tokens = append(tokens, Token{Type: TokenString, Literal: string(chars[start:pos]), Line: lineNum, Indent: indent})
				if pos < len(chars) {
					pos++
				}
			} else if unicode.IsLetter(ch) || ch == '_' {
				start := pos
				for pos < len(chars) && (unicode.IsLetter(chars[pos]) || unicode.IsDigit(chars[pos]) || chars[pos] == '_') {
					pos++
				}
				literal := string(chars[start:pos])
				tokType, ok := keywords[literal]
				if !ok {
					tokType = TokenIdent
				}
				tokens = append(tokens, Token{Type: tokType, Literal: literal, Line: lineNum, Indent: indent})
			} else if unicode.IsDigit(ch) {
				start := pos
				for pos < len(chars) && unicode.IsDigit(chars[pos]) {
					pos++
				}
				tokens = append(tokens, Token{Type: TokenNumber, Literal: string(chars[start:pos]), Line: lineNum, Indent: indent})
			} else {
				// Skip unknown single characters like punctuation if they don't match our specific needs
				pos++
			}
		}
	}

	tokens = append(tokens, Token{Type: TokenEOF, Literal: "", Line: len(lines) + 1, Indent: 0})
	return tokens
}
