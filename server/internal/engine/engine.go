package engine

// NOTE: Consider changing all array inputs to variadic parameters
import (
	"errors"
	"fmt"
)

func delItemFromSlice[S ~[]I, I comparable](s S, a I) (S, error) {
	for i, b := range s {
		if a == b {
			s[i] = s[len(s)-1]
			s = s[:len(s)-1]
			return s, nil
		}
	}
	return nil, errors.New("Item to delete not found in slice")
}

// TODO: How should it deal if the same key has two different values, which should it pick?
func mapUnion[M map[K]V, K comparable, V comparable](a M, b M) {
	for k, v := range b {
		if _, exist := a[k]; exist {
			if a[k] != v {
				panic("Map have the same key but different value.")
			}
		}
		a[k] = v
	}
}

type ObjNode struct {
	name    string
	format  string
	parents []*TagNode
}

func (obj ObjNode) Print() {
	fmt.Printf("\tName:%s\n\tFormat:%s\n", obj.name, obj.format)
}

type TagNode struct {
	name     string
	children []*TagNode
	parents  []*TagNode
	objects  []*ObjNode
}

func (node *TagNode) print() {
	fmt.Printf("node: %s\n", node.name)
	for _, obj := range node.objects {
		obj.Print()
	}
	for _, child := range node.children {
		child.print()
	}
}

// TODO: Add a way to edit the children of tags
// TODO: Add a way to edit the parents of tags
// TODO: Add a way to edit the parents of objects

type Engine struct {
	table map[string]*TagNode
}

func NewEngine() *Engine {
	root := TagNode{
		name: "Root",
	}
	table := make(map[string]*TagNode)
	table[root.name] = &root

	return &Engine{
		table: table,
	}
}

// TODO: Deal with parent not existing
// TODO: Deal with tag already existing
func (eng *Engine) AddTag(name string, parents []string) {
	if _, exist := eng.table[name]; exist {
		panic("Tag already exists")
	}

	node := TagNode{
		name: name,
	}

	if len(parents) == 0 {
		parents = append(parents, "Root")
	}

	for _, parent_name := range parents {
		parent, exist := eng.table[parent_name]
		if !exist {
			panic("Parent not found when adding tag")
		}
		parent.children = append(parent.children, &node)
		node.parents = append(node.parents, parent)
	}

	eng.table[name] = &node
}

// TODO: Deal with child not being acknowledged by its parent
// TODO: Deal with parent not being acknowledged by its children
// TODO: Deal with orphan due to deleted parent
func (eng *Engine) delTag(name string) error {
	node, exist := eng.table[name]
	if !exist {
		panic("Tag not found in table")
	}

	for _, parent := range node.parents {
		var err error
		parent.children, err = delItemFromSlice(parent.children, node)
		if err != nil {
			panic(err)
		}
	}

	for _, child := range node.children {
		var err error
		child.parents, err = delItemFromSlice(child.parents, node)
		if err != nil {
			panic(err)
		}
		if len(child.parents) == 0 {
			panic("Orphaned child")
		}
	}

	eng.table["Root"].objects = append(eng.table["Root"].objects, node.objects...)

	eng.table[name] = nil

	return nil
}

// NOTE: Should objects be in the table map in Engine?
// TODO: Deal with a tag not existing
func (eng *Engine) AddObj(name string, format string, tags []string) {
	obj := &ObjNode{
		name:   name,
		format: format,
	}

	if len(tags) == 0 {
		tags = append(tags, "Root")
	}

	for _, tag_name := range tags {
		tag, exist := eng.table[tag_name]
		if !exist {
			panic("Tag doesn't exist when adding object")
		}
		obj.parents = append(obj.parents, tag)
		tag.objects = append(tag.objects, obj)
	}
}

// TODO: Deal with object not being acknowledged by its parent
func (eng *Engine) delObj(obj *ObjNode) {
	for _, parent := range obj.parents {
		var err error
		parent.objects, err = delItemFromSlice(parent.objects, obj)
		if err != nil {
			panic(err)
		}
	}
}

func (eng *Engine) Print() {
	eng.table["Root"].print()
}

func (eng *Engine) getAllObjectsFromTag(tag *TagNode) map[string]*ObjNode {
	results := make(map[string]*ObjNode)
	for _, obj := range tag.objects {
		results[obj.name] = obj
	}

	for _, sub_tag := range tag.children {
		mapUnion(results, eng.getAllObjectsFromTag(sub_tag))
	}

	return results
}

// TODO: Deal with a tag not existing
func (eng *Engine) Query(tags []string) []*ObjNode {
	results := make(map[string]*ObjNode)
	for _, tag_name := range tags {
		tag, exist := eng.table[tag_name]
		if !exist {
			panic("Tag not found when querying")
		}
		mapUnion(results, eng.getAllObjectsFromTag(tag))
	}

	var resultList []*ObjNode
	for _, obj := range results {
		resultList = append(resultList, obj)
	}

	return resultList
}
