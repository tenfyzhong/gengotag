package gentag

type Struct struct {
	Name  string
	Field []*Field
}

type StructList []*Struct

func (l StructList) Len() int {
	return len(l)
}

func (l StructList) Less(i, j int) bool {
	return l[i].Name < l[j].Name
}

func (l StructList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
