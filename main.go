package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
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

	v, err := unmarshalJson(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unmarshal json: %v\n", err)
		os.Exit(3)
	}

	node := parseNode("", v)
	root := node
	for root.Type == TypeArray {
		if len(root.Children) == 0 {
			os.Exit(0)
		}
		root = root.Children[0]
	}

	gen := NewStructGenerator()
	gen.Generate(root)

	text := render(gen)
	fmt.Printf("%s", text)

}

func unmarshalJson(data []byte) (interface{}, error) {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func parseNode(tag string, v interface{}) *Node {
	node := &Node{
		Tag: tag,
	}
	switch c := v.(type) {
	case bool:
		node.Type = TypeBool
	case int:
		node.Type = TypeInt
	case float64:
		node.Type = TypeFloat64
	case string:
		node.Type = TypeString
	case map[string]interface{}:
		node.Type = TypeStruct
		node.Children = NodeList{}
		for k, v := range c {
			child := parseNode(k, v)
			if child == nil {
				continue
			}
			node.Children = append(node.Children, child)
		}
		sort.Sort(node.Children)
	case []interface{}:
		node.Type = TypeArray
		if len(c) > 0 {
			child := parseNode("", c[0])
			node.Children = append(node.Children, child)
		}
	}
	node.fingerprint = fingerprint(node)
	return node
}

func render(gen *StructGenerator) string {
	result := ""
	omitemptyStr := ""
	if omitempty {
		omitemptyStr = ",omitempty"
	}

	structList := StructList{}
	for _, st := range gen.m {
		structList = append(structList, st)
	}
	sort.Sort(structList)

	for _, st := range structList {
		result += fmt.Sprintf("type %s struct {\n", st.Name)
		for _, field := range st.Field {
			result += fmt.Sprintf("\t%s %s `%s:\"%s%s\"`\n", tag2Name(field.Tag), field.TypeName, tagtype, field.Tag, omitemptyStr)
		}
		result += "}\n\n"
	}
	return result
}
