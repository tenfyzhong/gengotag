package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read file %v\n", err)
		os.Exit(2)
	}

	text, err := gentag.Gen(data, tagtype, omitempty)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(3)
	}
	fmt.Printf("%s", text)
}
