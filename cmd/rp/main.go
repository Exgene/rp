package main

import (
	"fmt"
	"github.com/exgene/rp"
)

func main() {
	parser, err := rp.NewRegexEngine("[a-z]+")
	if err != nil {
		panic("Error: " + err.Error())
	}
	fmt.Print(parser.DoesMatch("a"))
	fmt.Print(parser.DoesMatch("b"))
	fmt.Print(parser.DoesMatch(""))
}
