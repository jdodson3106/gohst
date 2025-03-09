package gohst

import (
	"context"
	"fmt"
	"log"
	"strings"

	//"os

	cui "github.com/jroimartin/gocui"
)

const (
	StatusTop = iota
	StatusBottom

	// renderings
	WidgetPadding int = 1
	InputHeight   int = 2
	LineHeight    int = 2

	// view names
	MainGui        = "main"
	StatusBarView  = "status"
	SearchBarView  = "searchBar"
	AwsLoginView   = "awsLogin"
	ParametersView = "parameters"
	SecretsView    = "secrets"

	// titles
	SearchBarTitle     = "Search For Parameter Key"
	ParameterListTitle = "Parameters"
	SecretsListTitle   = "Secrets"
)

type frame struct {
	start point
	end   point
}

func (f frame) withPadding() frame {
	return frame{
		start: point{f.start.x + WidgetPadding, f.start.y + WidgetPadding},
		end:   point{f.end.x - WidgetPadding, f.end.y - WidgetPadding},
	}
}

func (f frame) height() int {
	return f.end.y - f.start.y
}

func (f frame) width() int {
	return f.end.x - f.start.x
}

func (f *frame) calcAvailableFrame(used frame) {
}

type point struct {
	x int
	y int
}

type UI struct {
	gui     *cui.Gui
	name    string
	views   []*cui.View
	session *session
}

func InitUI(g *cui.Gui) error {
	ui := &UI{
		gui:  g,
		name: MainGui,
		session: &session{
			context:    context.WithValue(context.Background(), "profile", "QAProfile"),
			isLoggedIn: true,
		},
	}
	g.Cursor = true
	g.SetManager(ui)
	return nil
}

func (u *UI) Layout(g *cui.Gui) error {
	maxX, maxY := g.Size()
	mainFrame := frame{
		start: point{0, 0},
		end:   point{maxX, maxY},
	}.withPadding()

	if err := u.buildStatusBar(StatusTop, mainFrame); err != nil {
		if err != cui.ErrUnknownView {
			log.Fatalf("error status bar: %s\n", err)
		}
	}

	if err := u.buildSearchView(maxX, maxY); err != nil {
		if err != cui.ErrUnknownView {
			log.Fatalf("error building search view: %s\n", err)
		}
		if _, err := g.SetCurrentView(SearchBarView); err != nil {
			return err
		}
	}

	if err := u.buildMainView(maxX, maxY); err != nil {
		if err != cui.ErrUnknownView {
			log.Fatalf("error building folder view: %s\n", err)
		}
	}

	return nil
}

func (u *UI) buildStatusBar(location int, f frame) error {
	var loggedInStatus string
	if u.session.isLoggedIn {
		loggedInStatus = "\033[0;32mActive"
	} else {
		loggedInStatus = "\033[0;31mInactive"
	}

	currentProfile, ok := u.session.profile()
	if !ok {
		currentProfile = "\033[0;31mNo Profile Found"
	} else {
		currentProfile = "\033[0;32m" + currentProfile
	}

	start := point{x: WidgetPadding, y: WidgetPadding}
	end := point{x: f.end.x - WidgetPadding, y: LineHeight + WidgetPadding}

	if v, err := u.buildView(start, end, StatusBarView); err != nil {
		if err != cui.ErrUnknownView {
			return err
		}

		v.Frame = false

		loggedInString := fmt.Sprintf("Session: %s", loggedInStatus)
		profileString := fmt.Sprint("\033[0;37mProfile: " + currentProfile)
		curX := f.end.x - len(profileString)
		spacing := (curX - len(loggedInString) - 1) - (WidgetPadding * 2)

		sb := strings.Builder{}
		sb.WriteString(loggedInString)
		for i := 0; i < spacing; i++ {
			sb.WriteString(" ")
		}
		sb.WriteString(profileString)
		v.Write([]byte(sb.String()))
	}
	return nil
}

func (u *UI) buildSearchView(maxX, maxY int) error {
	startY := WidgetPadding + LineHeight + 1

	start := point{x: WidgetPadding, y: startY}
	end := point{x: maxX - WidgetPadding, y: startY + InputHeight}

	//log.Fatalf("start: %+v, end: %+v\n", start, end)

	v, err := u.buildView(start, end, SearchBarView)
	if err != nil {
		return err
	}
	v.Title = SearchBarTitle
	v.Editable = true
	return nil
}

func (u *UI) buildMainView(maxX, maxY int) error {
	startY := WidgetPadding + LineHeight + 1 + InputHeight + 1
	ps := point{x: WidgetPadding, y: startY}
	pe := point{x: maxX/2 - 1, y: maxY - WidgetPadding}

	if err := u.buildParameterListView(ps, pe); err != nil {
		return err
	}

	ss := point{x: maxX/2 + 1, y: startY}
	se := point{x: maxX - WidgetPadding, y: maxY - WidgetPadding}
	if err := u.buildSecretsListView(ss, se); err != nil {
		return err
	}

	return nil
}

func (u *UI) buildParameterListView(start, end point) error {
	v, err := u.buildView(start, end, ParametersView)
	if err != nil && err != cui.ErrUnknownView {
		log.Fatalln("error building parameters view: ", err)
	}

	v.Title = ParameterListTitle
	return err
}

func (u *UI) buildSecretsListView(start, end point) error {
	v, err := u.buildView(start, end, SecretsView)
	if err != nil && err != cui.ErrUnknownView {
		log.Fatalln("error building secrets view: ", err)
	}

	v.Title = SecretsListTitle
	return err
}

func (u *UI) buildView(start, end point, name string) (*cui.View, error) {
	if v, err := u.gui.SetView(name, start.x, start.y, end.x, end.y); err != nil {
		if err != cui.ErrUnknownView {
			return nil, err
		}
		u.views = append(u.views, v)
		return v, err
	}
	return u.lookupViewByName(name)
}

func (u *UI) lookupViewByName(name string) (*cui.View, error) {
	for _, v := range u.views {
		view := *v
		if view.Name() == name {
			return v, nil
		}
	}
	return nil, cui.ErrUnknownView
}
