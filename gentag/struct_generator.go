package gentag

import (
	"fmt"
	"math/rand"
	"strings"
)

type StructGenerator struct {
	m     map[string]*Struct
	nameM map[string]bool
}

func NewStructGenerator() *StructGenerator {
	return &StructGenerator{
		m:     make(map[string]*Struct),
		nameM: make(map[string]bool),
	}
}

func (gen *StructGenerator) Generate(node *Node) *Struct {
	if node.Type != TypeStruct {
		return nil
	}

	st, ok := gen.m[node.Fingerprint()]
	if !ok {
		st = &Struct{
			Name: gen.genName(node.Tag),
		}
		for _, child := range node.Children {
			field := &Field{
				Tag: child.Tag,
			}
			field.TypeName = gen.nodeTypeName(child)
			if field.TypeName != "" {
				st.Field = append(st.Field, field)
			}
		}
		gen.m[node.Fingerprint()] = st
	}

	return st
}

func (gen *StructGenerator) genName(str string) string {
	name := tag2Name(str)
	if name == "" {
		name = "Struct"
	}

	nameTmp := name
	for {
		if gen.nameM[nameTmp] {
			nameTmp = name + fmt.Sprintf("%02d", rand.Intn(100))
			continue
		}
		gen.nameM[nameTmp] = true
		break
	}
	return nameTmp
}

func (gen *StructGenerator) nodeTypeName(node *Node) string {
	name := ""
	switch node.Type {
	case TypeArray:
		if len(node.Children) == 0 {
			name = "[]interface{}"
		} else {
			sub := gen.nodeTypeName(node.Children[0])
			name = "[]" + sub
		}
	case TypeStruct:
		sub := gen.Generate(node)
		name = "*" + sub.Name
	default:
		name = basicType(node.Type)
	}
	return name
}

func basicType(typ NodeType) string {
	switch typ {
	case TypeNull:
		return "interface{}"
	case TypeBool:
		return "bool"
	case TypeFloat64:
		return "float64"
	case TypeString:
		return "string"
	}
	return ""
}

func tag2Name(str string) string {
	if str == "" {
		return ""
	}
	items := strings.Split(str, "_")
	result := ""
	for _, item := range items {
		if item == "" {
			continue
		}
		s := strings.ToUpper(item[0:1]) + item[1:]
		result += s
	}
	return result
}
