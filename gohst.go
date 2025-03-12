/*
Copyright © 2025 Justin Dodson <EMAIL ADDRESS>
*/

package gohst

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func Start() {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Fatalln(err)
	}
	defer g.Close()

	if err := InitUI(g); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	//cmd.Execute()
}

func quit(g *gocui.Gui, view *gocui.View) error {
	return gocui.ErrQuit
}

func Playground() {
	fmt.Println("running playground")
	t := session{}
	t.listProfiles()
}
