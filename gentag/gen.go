package gentag

import (
	"encoding/json"
	"fmt"
	"sort"
)

func Gen(data []byte, tagtype string, omitempty bool) (string, error) {
	v, err := unmarshalJson(data)
	if err != nil {
		return "", err
	}

	node := parseNode("", v)
	root := node
	for root.Type == TypeArray {
		if len(root.Children) == 0 {
			return "", nil
		}
		root = root.Children[0]
	}

	gen := NewStructGenerator()
	gen.Generate(root)

	text := render(gen, tagtype, omitempty)
	return text, nil
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
	case float64:
		node.Type = TypeFloat64
	case string:
		node.Type = TypeString
	case map[string]interface{}:
		node.Type = TypeStruct
		node.Children = NodeList{}
		for k, v := range c {
			child := parseNode(k, v)
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

func render(gen *StructGenerator, tagtype string, omitempty bool) string {
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
