package parser

import "errors"

var (
	ErrUnexpectedToken = errors.New("Unexpected token found")
	ErrUnexpectedEdge  = errors.New("The edge type is not defined")
)
