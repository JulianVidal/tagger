package engine

// TODO: Remove "Root" tag, it is not needed
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

// TODO: Add a way to edit the children of tags
// TODO: Add a way to edit the parents of tags
// TODO: Add a way to edit the parents of objects

var table map[string]*TagNode

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

func InitEngine() {
	table = make(map[string]*TagNode)
}

func AddTag(name string, parents []string) error {
	if _, exist := table[name]; exist {
		return fmt.Errorf("Tag '%s' already exists", name)
	}

	node := TagNode{
		name: name,
	}

	for _, parent_name := range parents {
		if _, exist := table[parent_name]; !exist {
			return fmt.Errorf("Parent tag '%s' not found", parent_name)
		}
	}

	for _, parent_name := range parents {
		parent, _ := table[parent_name]
		parent.children = append(parent.children, &node)
		node.parents = append(node.parents, parent)
	}

	table[name] = &node

	return nil
}

func DelTag(name string) error {
	tag, exist := table[name]
	if !exist {
		return fmt.Errorf("Tag not found in engine: %s", name)
	}

	for _, child := range tag.children {
		if len(child.parents) == 1 {
			return fmt.Errorf("Tag '%s' would orphan child '%s'", name, child.name)
		}
	}

	for _, object := range tag.objects {
		if len(object.parents) == 1 {
			return fmt.Errorf("Tag '%s' would orphan object '%s'", name, object.name)
		}
	}

	for _, child := range tag.children {
		child.parents, _ = delItemFromSlice(child.parents, tag)
	}

	for _, parent := range tag.parents {
		parent.children, _ = delItemFromSlice(parent.children, tag)
	}

	table[name] = nil

	return nil
}

// NOTE: Should objects be in the table map in Engine?
func AddObj(name string, format string, tags []string) error {
	obj := &ObjNode{
		name:   name,
		format: format,
	}

	if len(tags) == 0 {
		return fmt.Errorf("No tags were provided")
	}

	for _, tag_name := range tags {
		if _, exist := table[tag_name]; !exist {
			return fmt.Errorf("Tag '%s' doesn't exist in engine", tag_name)
		}
	}

	for _, tag_name := range tags {
		tag, _ := table[tag_name]
		obj.parents = append(obj.parents, tag)
		tag.objects = append(tag.objects, obj)
	}

	return nil
}

// TODO: Deal with object not being acknowledged by its parent
func DelObj(obj *ObjNode) error {
	for _, parent := range obj.parents {
		parent.objects, _ = delItemFromSlice(parent.objects, obj)
	}
	return nil
}

func Print() {
	for _, v := range table {
		v.print()
	}
}

func getAllObjectsFromTag(tag *TagNode) map[string]*ObjNode {
	results := make(map[string]*ObjNode)
	for _, obj := range tag.objects {
		results[obj.name] = obj
	}

	for _, sub_tag := range tag.children {
		mapUnion(results, getAllObjectsFromTag(sub_tag))
	}

	return results
}

// TODO: Deal with a tag not existing
func Query(tags []string) ([]*ObjNode, error) {
	results := make(map[string]*ObjNode)

	for _, tag_name := range tags {
		if _, exist := table[tag_name]; !exist {
			return nil, fmt.Errorf("Tag '%s' not found", tag_name)
		}
	}

	for _, tag_name := range tags {
		tag, _ := table[tag_name]
		mapUnion(results, getAllObjectsFromTag(tag))
	}

	var resultList []*ObjNode
	for _, obj := range results {
		resultList = append(resultList, obj)
	}

	return resultList, nil
}
