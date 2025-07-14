package gohst

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jdodson3106/gohst/internal/auth"
	"github.com/jroimartin/gocui"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"golang.org/x/term"
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

	// TODO: move to flags
	hostAddr, err := readLine("Host IP: ")
	if err != nil {
		panic(err)
	}

	// TODO: move to flags
	keyPath, err := readLine("Enter path to private key: ")
	if err != nil {
		panic(err)
	}

	f, err := os.ReadFile(string(keyPath))
	if err != nil {
		panic(err)
	}

	requiresPW := false
	signer, err := ssh.ParsePrivateKey(f)
	if err != nil {
		fmt.Printf("%+v\n", err)
		if err.Error() == "ssh: this private key is passphrase protected" {
			requiresPW = true
		} else {
			panic(err)
		}
	}

	if requiresPW {
		fmt.Print("Enter passphrase: ")
		pw, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Printf("error reading input: %s\n", err)
			panic(err)
		}
		signer, err = ssh.ParsePrivateKeyWithPassphrase(f, pw)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("\nConnecting to host %s\n", hostAddr)

	// TODO: save config in file
	cb, err := knownhosts.New("/Users/jdd/.ssh/known_hosts")
	if err != nil {
		panic(err)
	}

	cc := ssh.ClientConfig{
		User:            "jdd",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: cb,
		Timeout:         2 * time.Second,
	}

	hostAddr = append(hostAddr, ":22"...)
	client, err := ssh.Dial("tcp", string(hostAddr), &cc)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf

	if err := session.Run("whoami"); err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
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
