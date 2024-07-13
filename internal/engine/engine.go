package engine

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var tagMap map[string]*Tag
var objectMap map[string]*Object

func Init(target_dir string, engine_path string) {
	data, err := os.ReadFile(engine_path)
	if err != nil {
		New()
		return
	}

	err = FromJson(data)
	if err != nil {
		panic(err)
	}

	// Walks target directory, only follows 1 symlink
	err = filepath.Walk(target_dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path == target_dir || info.IsDir() {
				return nil
			}

			// If it is a symlink it reads it
			if info.Mode()&fs.ModeSymlink != 0 {
				sym_path, err := os.Readlink(path)
				if err != nil {
					return err
				}

				sym_info, err := os.Stat(sym_path)
				if err != nil {
					panic(err)
				}

				if sym_info.IsDir() {
					err = filepath.Walk(sym_path,
						func(in_sym string, info os.FileInfo, err error) error {
							if err != nil {
								return err
							}
							if in_sym == path || info.IsDir() {
								return nil
							}
							p := filepath.Join(filepath.Base(sym_path), strings.TrimPrefix(in_sym, sym_path+string(os.PathSeparator)))
							NewObject(p) // Ignore error if object already exists
							return nil
						})
				} else {
					p := strings.TrimPrefix(path, target_dir+string(os.PathSeparator))
					NewObject(p) // Ignore error if object already exists
				}

			} else {
				p := strings.TrimPrefix(path, target_dir+string(os.PathSeparator))
				NewObject(p) // Ignore error if object already exists
			}
			return nil
		})
	if err != nil {
		panic(err)
	}
}

func New() {
	tagMap = make(map[string]*Tag)
	objectMap = make(map[string]*Object)
}

func Print() {
	fmt.Println("Printing engine:")
	for _, v := range tagMap {
		v.Print()
	}
}

func String() string {
	str := ""
	str += fmt.Sprintln("Printing engine:")
	for _, v := range tagMap {
		str += v.String()
	}
	for _, v := range objectMap {
		str += v.String() + "\n"
	}
	return str
}

func FindTag(name string) (*Tag, bool) {
	tag, exist := tagMap[name]
	return tag, exist
}

func FindObject(name string) (*Object, bool) {
	object, exist := objectMap[name]
	return object, exist
}

func Tags() []string {
	return mapKeys(tagMap)
}

func Objects() []string {
	return mapKeys(objectMap)
}

func getAllObjectsFromTag(tag *Tag) map[string]*Object {
	results := make(map[string]*Object)
	for _, object := range tag.objects {
		results[object.Name()] = object
	}

	for _, childTag := range tag.children {
		mapUnion(results, getAllObjectsFromTag(childTag))
	}

	return results
}

func Query(tags ...*Tag) ([]Object, error) {
	if len(tags) == 0 {
		var resultList []Object
		for _, object := range objectMap {
			resultList = append(resultList, *object)
		}
		return resultList, nil
	}

	results := make(map[string]*Object)

	for _, tag := range tags {
		if _, exist := tagMap[tag.name]; !exist {
			return nil, fmt.Errorf("Tag '%s' not found", tag)
		}
	}

	for _, tag := range tags {
		mapUnion(results, getAllObjectsFromTag(tag))
	}

	var resultList []Object
	for _, object := range results {
		resultList = append(resultList, *object)
	}

	return resultList, nil
}

func TagExists(tag string) bool {
	_, ok := tagMap[tag]
	return ok
}

func ObjectExists(object string) bool {
	_, ok := objectMap[object]
	return ok
}
