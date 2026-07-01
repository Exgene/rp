package main

import (
	"fmt"
	"github.com/exgene/rp"
)

func main() {
	parser, err := rp.NewRegexEngine("[a-z]+@[a-z]+.[a-z]{2,}")
	parser.PrintNFA()
	if err != nil {
		panic("Error: " + err.Error())
	}
	fmt.Print(parser.MatchString("kausthubhjrao@gmail.com"))
	fmt.Print(parser.MatchString("kausthubhjrao@gmail.co"))
	fmt.Print(parser.MatchString("kausthubhjrao@gmail.o"))
	// fmt.Print(parser.MatchString("b"))
	// fmt.Print(parser.MatchString(""))
}
