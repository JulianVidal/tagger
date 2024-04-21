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

type Object struct {
	name    string
	format  string
	parents []*Node
}

func (obj Object) Print() {
	fmt.Printf("\tName:%s\n\tFormat:%s\n", obj.name, obj.format)
}

type Node struct {
	name     string
	children []*Node
	parents  []*Node
	objects  []*Object
}

func (node *Node) print() {
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
	table map[string]*Node
}

func NewEngine() *Engine {
	root := Node{
		name: "Root",
	}
	table := make(map[string]*Node)
	table[root.name] = &root

	return &Engine{
		table: table,
	}
}

// TODO: Deal with parent not existing
func (eng *Engine) AddTag(name string, parents []string) {
	node := Node{
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

func (eng *Engine) delTag(name string) error {
	node, exist := eng.table[name]
	if !exist {
		return errors.New("Tag not found in table")
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
	}

	eng.table["Root"].objects = append(eng.table["Root"].objects, node.objects...)

	eng.table[name] = nil

	return nil
}

// NOTE: Should objects be in the table map in Engine?
// TODO: Deal with a tag not existing
func (eng *Engine) AddObj(name string, format string, tags []string) {
	obj := &Object{
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

func (eng *Engine) delObj(obj *Object) {
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

func (eng *Engine) getAllObjectsFromTag(tag *Node) map[string]*Object {
	results := make(map[string]*Object)
	for _, obj := range tag.objects {
		results[obj.name] = obj
	}

	for _, sub_tag := range tag.children {
		mapUnion(results, eng.getAllObjectsFromTag(sub_tag))
	}

	return results
}

// TODO: Deal with a tag not existing
func (eng *Engine) Query(tags []string) []*Object {
	results := make(map[string]*Object)
	for _, tag_name := range tags {
		tag, exist := eng.table[tag_name]
		if !exist {
			panic("Tag not found when querying")
		}
		mapUnion(results, eng.getAllObjectsFromTag(tag))
	}

	var resultList []*Object
	for _, obj := range results {
		resultList = append(resultList, obj)
	}

	return resultList
}
