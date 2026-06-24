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
	fmt.Print(parser.MatchString("a"))
	fmt.Print(parser.MatchString("b"))
	fmt.Print(parser.MatchString(""))
}
