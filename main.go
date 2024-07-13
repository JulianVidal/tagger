package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/JulianVidal/tagger/app"
	"github.com/JulianVidal/tagger/internal/engine"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Gets the target directory
	// Uses $XDG_DATA_HOME if it exists otherwise $HOME/.local/share
	target_dir, exists := os.LookupEnv("XDG_DATA_HOME")
	if !exists {
		target_dir, exists = os.LookupEnv("HOME")
		if !exists {
			panic("HOME environment variable doesn't exists")
		}
		target_dir = filepath.Join(target_dir, ".local/share")
	}
	target_dir = filepath.Join(target_dir, "tagger")

	if _, err := os.Stat(target_dir); os.IsNotExist(err) {
		err := os.Mkdir(target_dir, os.ModeDir)
		if err != nil {
			panic(err)
		}
	}

	engine.Init(target_dir, "engine.json")

	m := app.New()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	// bt_tag, err := engine.NewTag("BT")
	// if err != nil {
	// 	panic(err)
	// }
	//
	// sat_tag, err := engine.NewTag("Sat")
	// if err != nil {
	// 	panic(err)
	// }
	// err = sat_tag.AddTags(bt_tag)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// sat_obj, err := engine.NewObject("/home/julian/temp/satshelf.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// err = sat_obj.AddTags(sat_tag)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// start_obj, err := engine.NewObject("/home/julian/temp/starter.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// err = start_obj.AddTags(bt_tag)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// engine.Print()
	//
	// bt_tag, exists := engine.FindTag("BT")
	// if !exists {
	// 	panic("Tag doesn't exist")
	// }
	// objs, err := engine.Query(bt_tag)
	// if err != nil {
	// 	panic(err)
	// }
	// for _, obj := range objs {
	// 	obj.Print()
	// }
	//
	// data, err := engine.ToJson()
	// if err != nil {
	// 	panic(err)
	// }
	// err = os.WriteFile("engine.json", data, 0644)
	// if err != nil {
	// 	panic(err)
	// }
}
