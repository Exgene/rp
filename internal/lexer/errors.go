package lexer

import "errors"

var (
	ErrEmptyRepeatTarget                = errors.New("Missing target for Repeat Sequence.")
	ErrUndefinedRepeatType              = errors.New("Undefined repeat type.")
	ErrRepeatValueBetweenFlowerBrackets = errors.New("Invalid value between the flower bracket syntax.")
	ErrUnexpectedToken                  = errors.New("Unexpected token found.")
	ErrExpectedValueShapeMismatch       = errors.New("Value interface for any does not satisfy the lexer type.")
	ErrMissingClosingSquareBracket      = errors.New("No matching closing bracket found.")
	ErrMissingClosingParen              = errors.New("No matching closing parantheses found.")
	ErrMissingClosingFlowerBracket      = errors.New("No matching closing flower bracket found.")
	ErrInvalidToken                     = errors.New("This token should not be present here.")
)
