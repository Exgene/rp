package main

import (
	"fmt"
	"github.com/exgene/rp"
)

func main() {
	parser := rp.NewRegexEngine("[a-z]+")
	fmt.Print(parser.DoesMatch("a"))
	fmt.Print(parser.DoesMatch("b"))
	fmt.Print(parser.DoesMatch(""))
}
