package gentag

import (
	"fmt"
	"strings"
)

type NodeType int

const (
	TypeNull NodeType = iota
	TypeBool
	TypeFloat64
	TypeString
	TypeArray
	TypeStruct
)

type Node struct {
	Tag         string
	Type        NodeType
	Children    NodeList
	fingerprint string
}

func (node *Node) Fingerprint() string {
	if node.fingerprint == "" {
		node.fingerprint = fingerprint(node)
	}
	return node.fingerprint

}

type NodeList []*Node

func (l NodeList) Len() int {
	return len(l)
}

func (l NodeList) Less(i, j int) bool {
	return l[i].Tag < l[j].Tag
}

func (l NodeList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type Field struct {
	Tag      string
	TypeName string
}

func fingerprint(node *Node) string {
	switch node.Type {
	case TypeArray:
		child := "nil"
		if len(node.Children) > 0 {
			child = node.Children[0].Fingerprint()
		}
		return fmt.Sprintf("%s.%d:%s", node.Tag, node.Type, child)
	case TypeStruct:
		children := make([]string, 0)
		for _, child := range node.Children {
			finger := child.Fingerprint()
			children = append(children, finger)
		}
		return strings.Join(children, ">")
	default:
		return fmt.Sprintf("%s.%d", node.Tag, node.Type)
	}
}
