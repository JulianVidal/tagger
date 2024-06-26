package main

import (
	"fmt"
	"os"

	"github.com/JulianVidal/tagger/app"
	"github.com/JulianVidal/tagger/internal/engine"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	engine.Init()
	m := app.New()
	// m.SetTags(engine.Tags()...)
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
