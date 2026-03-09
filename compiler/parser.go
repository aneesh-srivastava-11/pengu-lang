package compiler

import (
	"fmt"
)

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) current() Token {
	if p.pos >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.pos]
}

func (p *Parser) expect(expectedType TokenType) error {
	tok := p.current()
	if tok.Type != expectedType {
		return fmt.Errorf("line %d: expected %v, got %v (%s)", tok.Line, expectedType, tok.Type, tok.Literal)
	}
	p.pos++
	return nil
}

func (p *Parser) match(expectedType TokenType) bool {
	if p.current().Type == expectedType {
		p.pos++
		return true
	}
	return false
}

func (p *Parser) Parse() (*Service, error) {
	service := &Service{}

	// Parse Version
	if err := p.expect(TokenVersion); err != nil {
		return nil, fmt.Errorf("line %d: file must start with version declaration", p.current().Line)
	}
	if p.current().Type != TokenNumber {
		return nil, fmt.Errorf("line %d: expected version number, got %v", p.current().Line, p.current().Literal)
	}
	service.Version = p.current().Literal
	service.Line = p.current().Line
	p.pos++

	// Parse Service
	if err := p.expect(TokenService); err != nil {
		return nil, fmt.Errorf("line %d: expected service declaration", p.current().Line)
	}
	if p.current().Type != TokenIdent {
		return nil, fmt.Errorf("line %d: expected service name, got %v", p.current().Line, p.current().Literal)
	}
	service.Name = p.current().Literal
	p.pos++

	// Parse Routes
	for p.current().Type != TokenEOF {
		if p.current().Type == TokenRoute {
			route, err := p.parseRoute()
			if err != nil {
				return nil, err
			}
			service.Routes = append(service.Routes, route)
		} else {
			return nil, fmt.Errorf("line %d: unexpected token %s", p.current().Line, p.current().Literal)
		}
	}

	return service, nil
}

func (p *Parser) parseRoute() (Route, error) {
	routeLine := p.current().Line
	routeIndent := p.current().Indent
	p.pos++ // consume 'route'

	// Method
	if p.current().Type != TokenIdent {
		return Route{}, fmt.Errorf("line %d: expected HTTP method (GET, POST, etc.)", p.current().Line)
	}
	method := p.current().Literal
	// Basic validation of valid methods
	validMethods := map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true, "PATCH": true}
	if !validMethods[method] {
		return Route{}, fmt.Errorf("line %d: invalid HTTP method %s", p.current().Line, method)
	}
	p.pos++

	// Path
	if p.current().Type != TokenString {
		return Route{}, fmt.Errorf("line %d: expected route path as string", p.current().Line)
	}
	path := p.current().Literal
	p.pos++

	route := Route{
		Method: method,
		Path:   path,
		Line:   routeLine,
	}

	// Parse Actions inside route
	// Actions must have greater indentation than the 'route' keyword
	for p.current().Type != TokenEOF && p.current().Type != TokenRoute {
		if p.current().Indent <= routeIndent {
			return Route{}, fmt.Errorf("line %d: actions inside routes must be indented", p.current().Line)
		}
		
		action, err := p.parseAction()
		if err != nil {
			return Route{}, err
		}
		route.Actions = append(route.Actions, action)
	}

	return route, nil
}

func (p *Parser) parseAction() (Action, error) {
	tok := p.current()
	action := Action{Line: tok.Line}

	switch tok.Type {
	case TokenLog:
		p.pos++
		action.Type = "log"
		if p.current().Type != TokenString {
			return Action{}, fmt.Errorf("line %d: log requires a message string", p.current().Line)
		}
		action.Args = append(action.Args, p.current().Literal)
		p.pos++
	case TokenRespond:
		p.pos++
		action.Type = "respond"
		if p.current().Type != TokenNumber {
			return Action{}, fmt.Errorf("line %d: respond requires status code and message", p.current().Line)
		}
		action.Args = append(action.Args, p.current().Literal)
		p.pos++
		
		if p.current().Type != TokenString {
			return Action{}, fmt.Errorf("line %d: respond requires status code and message", p.current().Line)
		}
		action.Args = append(action.Args, p.current().Literal)
		p.pos++
	default:
		return Action{}, fmt.Errorf("line %d: unknown action '%s'", tok.Line, tok.Literal)
	}

	return action, nil
}
