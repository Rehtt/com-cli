package main

import (
	"github.com/Rehtt/gocui"
)

type App struct {
	gui  *gocui.Gui
	info *Info

	lastView     *gocui.View
	settingsView *gocui.View
	switchView   *gocui.View
	displayView  *gocui.View
	inputView    *gocui.View
	viewArr      []string
	active       int

	settings             []settings
	settingMap           map[string]settings
	currentSettingOption settings
	port                 *serialPort
}

func NewApp(info *Info) (*App, error) {
	a := &App{
		info: info,
	}
	if err := a.initGui(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run() error {
	return a.gui.MainLoop()
}

func (a *App) Close() error {
	a.gui.Close()
	return gocui.ErrQuit
}
