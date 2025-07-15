package gohst

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/jdodson3106/gohst/internal/auth"
	"github.com/jdodson3106/gohst/internal/ssh"
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
	t := auth.Session{}
	t.ListProfiles()
}

func Run() {
	client := ssh.NewClientWithKeyHost("34.0.128.156", ssh.RSAKeyHostDefault())
	client.Connect()
	defer client.Close()

	out, err := client.RunCommand("whoami")
	if err != nil {
		log.Fatalf("error running command: %s", err)
	}
	fmt.Println(string(out))
}

func readLine(prompt string) ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)

	in, err := reader.ReadSlice('\n')
	if len(in) == 0 {
		return nil, err
	}
	return in[:len(in)-1], err
}
