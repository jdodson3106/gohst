/*
Copyright © 2025 Justin Dodson <EMAIL ADDRESS>
*/

package main

import (
	"log"

	"github.com/jroimartin/gocui"
)

func main() {
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
