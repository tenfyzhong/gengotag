package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tenfyzhong/gengotag/gentag"
)

var (
	tagtype   = ""
	file      = ""
	omitempty = false
)

func init() {
	flag.StringVar(&tagtype, "type", "json", "type of tag")
	flag.StringVar(&file, "file", "", "the file to read")
	flag.BoolVar(&omitempty, "omitempty", false, "omitempty tag")
	flag.Parse()
	if file == "" {
		fmt.Fprintf(os.Stderr, "filename is empty\n")
		os.Exit(1)
	}
}

func main() {
	text, err := gentag.Gen(file, tagtype, omitempty)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	fmt.Printf("%s", text)
}
