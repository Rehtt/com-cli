package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/Rehtt/gocui"
	"go.bug.st/serial"
)

func (a *App) initGui() (err error) {
	a.gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return
	}
	a.gui.InputEsc = true
	a.gui.SelFgColor = gocui.ColorGreen
	a.gui.Highlight = true
	a.gui.Mouse = true

	a.layout()
	if err = a.keybindings(); err != nil {
		return
	}

	a.registerSettings(
		NewPort(),
		NewBaudRate(),
		NewStopBits(),
		NewDataBits(),
		NewParity(),
		NewDisplayMode(),
		NewInputMode(),
	)

	return
}

func (a *App) layout() {
	a.gui.SetManagerFunc(func(g *gocui.Gui) error {
		maxX, maxY := a.gui.Size()
		// 设置区
		if v, err := a.gui.SetView("settings", 0, 0, (maxX*15/100)-1, (maxY * 70 / 100)); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Highlight = true
			v.SelBgColor = gocui.ColorGreen
			v.SelFgColor = gocui.ColorBlack
			v.Title = "Settings"
			a.viewArr = append(a.viewArr, v.Name())

			a.settingsView = v
			g.SetCurrentView("settings")
		}
		if v, err := a.gui.SetView("switch", 0, (maxY*70/100)+1, (maxX*15/100)-1, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			a.viewArr = append(a.viewArr, v.Name())
			a.switchView = v
		}
		if v, err := g.SetView("display", maxX*15/100, 0, maxX-1, (maxY * 70 / 100)); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Display"
			v.Frame = true
			v.Autoscroll = true
			v.Wrap = true
			a.viewArr = append(a.viewArr, v.Name())
			a.displayView = v
		}
		if v, err := g.SetView("input", maxX*15/100, maxY*70/100+1, maxX-1, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Input"
			v.Editable = true
			v.Wrap = true
			a.viewArr = append(a.viewArr, v.Name())
			a.inputView = v
		}

		a.Refresh()
		return nil
	})
}

func (a *App) registerSettings(s ...settings) {
	a.settings = append(a.settings, s...)
	if a.settingMap == nil {
		a.settingMap = map[string]settings{}
	}
	for _, v := range s {
		a.settingMap[v.Key()] = v
	}
}

func (a *App) keybindings() error {
	a.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return a.Close()
	})
	a.gui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, a.nextView)
	a.gui.SetKeybinding("error", gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return a.cloasError()
	})
	a.gui.SetKeybinding("error", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return a.cloasError()
	})
	a.gui.SetKeybinding("settings", gocui.KeyArrowUp, gocui.ModNone, upDown(-1))
	a.gui.SetKeybinding("settings", gocui.KeyArrowDown, gocui.ModNone, upDown(1))
	a.gui.SetKeybinding("settings", gocui.KeyEnter, gocui.ModNone, a.settingOptions)
	a.gui.SetKeybinding("settingOptions", gocui.KeyArrowDown, gocui.ModNone, upDown(1))
	a.gui.SetKeybinding("settingOptions", gocui.KeyArrowUp, gocui.ModNone, upDown(-1))
	a.gui.SetKeybinding("settingOptions", gocui.KeyEnter, gocui.ModNone, a.settingOptionEnter)
	a.gui.SetKeybinding("switch", gocui.KeyEnter, gocui.ModNone, a.switchEnter)
	a.gui.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, a.inputEnter)

	return nil
}

func (a *App) inputEnter(g *gocui.Gui, v *gocui.View) error {
	if a.port == nil || !a.port.IsRun() {
		a.ErrorMsg("未运行")
		return nil
	}
	_, cy := v.Cursor()
	for i := 0; i <= cy; i++ {
		s, err := v.Line(i)
		if err == nil {
			if err = a.port.WriteString(s); err != nil {
				a.ErrorMsg(err.Error())
				break
			}
		}
	}
	v.Clear()
	v.SetCursor(0, 0)
	return nil
}

func (a *App) switchEnter(g *gocui.Gui, v *gocui.View) error {
	if a.port != nil && a.port.IsRun() {
		a.port.Close()
		a.port = nil
	} else {
		var err error
		a.port, err = openSerial(a.settingMap["port"].Get(), &serial.Mode{
			BaudRate: a.settingMap["baud_rate"].Value().(int),
			Parity:   a.settingMap["parity"].Value().(serial.Parity),
			DataBits: a.settingMap["data_bits"].Value().(int),
			StopBits: a.settingMap["stop_bits"].Value().(serial.StopBits),
		})
		if err != nil {
			a.ErrorMsg(err.Error())
			return nil
		}
		a.port.SetDisplayMode(a.settingMap["display_mode"].Value().(int))
		a.port.SetInputMode(a.settingMap["input_mode"].Value().(int))
		go a.port.HandleRead(func(data []byte) {
			g.Update(func(g *gocui.Gui) error {
				data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
				a.displayView.Write(data)
				return nil
			})
		})
	}
	a.Refresh()
	return nil
}

func (a *App) nextView(g *gocui.Gui, v *gocui.View) error {
	var find bool
	for _, name := range a.viewArr {
		if name == v.Name() {
			find = true
		}
	}
	if !find {
		return nil
	}
	nextIndex := (a.active + 1) % len(a.viewArr)
	name := a.viewArr[nextIndex]
	g.SetCurrentView(name)
	g.SetViewOnTop(name)

	switch name {
	case "input", "display":
		g.Cursor = true
	default:
		g.Cursor = false
	}
	// if nextIndex == 2 || nextIndex == 1 {
	// 	g.Cursor = true
	// } else {
	// 	g.Cursor = false
	// }

	a.active = nextIndex
	return nil
}

func (a *App) settingOptionEnter(g *gocui.Gui, v *gocui.View) error {
	if a.currentSettingOption == nil {
		return errors.New("settingOptionEnter: currentSettingOption is nil")
	}
	_, i := v.Cursor()
	list, _, err := a.currentSettingOption.Trigger()
	if err != nil {
		return err
	}
	if len(list) <= i {
		return errors.New("settingOptionEnter: index out of range")
	}
	if err = a.currentSettingOption.Set(list[i]); err != nil {
		return err
	}
	a.Refresh()
	g.DeleteView(v.Name())
	g.SetCurrentView(a.lastView.Name())
	return nil
}

func (a *App) settingOptions(g *gocui.Gui, v *gocui.View) error {
	if a.port != nil && a.port.IsRun() {
		a.ErrorMsg("运行中，请停止再进行设置")
		return nil
	}
	maxX, maxY := g.Size()
	_, y := v.Cursor()
	s := a.settings[y]
	list, show, err := s.Trigger()
	if err != nil {
		a.ErrorMsg(err.Error())
		return nil
	}
	if show {
		// 选项窗口
		view, err := g.SetView("settingOptions", maxX/2-(maxX/4), maxY/2-(maxY/4), maxX/2+(maxX/4), maxY/2+(maxY/4))
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
		view.Highlight = true
		view.SelBgColor = gocui.ColorGreen
		view.SelFgColor = gocui.ColorBlack

		// 定位
		oldValue := s.Get()
		for i, value := range list {
			if oldValue == value {
				view.SetCursor(0, i)
				break
			}
		}

		a.currentSettingOption = s
		fmt.Fprint(view, strings.Join(list, "\n"))
		_, err = g.SetCurrentView("settingOptions")

		a.lastView = v
	}
	return err
}

func upDown(d int) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if v != nil {
			cx, cy := v.Cursor()
			ox, oy := v.Origin()
			var err error
			if cy+d+oy > -1 && cy+d+oy < v.LineLen() {
				err = v.SetCursor(cx, cy+d)
			}
			if err != nil {
				if oy+d+cy > -1 && oy+d+cy < v.LineLen() {
					err = v.SetOrigin(ox, oy+d)
				}
			}
			return err
			// if err := v.SetCursor(cx, cy+d); err != nil {
			// 	ox, oy := v.Origin()
			// 	if err := v.SetOrigin(ox, oy+d); err != nil {
			// 		return err
			// 	}
			// }
		}
		return nil
	}
}

func (a *App) ErrorMsg(errStr string) error {
	a.lastView = a.gui.CurrentView()
	maxX, maxY := a.gui.Size()
	v, err := a.gui.SetView("error", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	fmt.Fprintln(v, errStr)
	a.gui.SelFgColor = gocui.ColorRed
	a.gui.SetCurrentView("error")

	return nil
}

func (a *App) cloasError() error {
	if err := a.gui.DeleteView("error"); err != nil {
		return err
	}
	if a.lastView != nil {
		if _, err := a.gui.SetCurrentView(a.lastView.Name()); err != nil {
			return err
		}
	}
	a.gui.SelFgColor = gocui.ColorGreen
	return nil
}

func (a *App) Refresh() {
	a.gui.Update(func(g *gocui.Gui) error {
		a.displayView.Title = fmt.Sprintf("Display-(%s)", a.settingMap["input_mode"].Get())
		a.inputView.Title = fmt.Sprintf("Input-(%s)", a.settingMap["display_mode"].Get())
		{
			a.settingsView.Clear()
			buf := make([]string, 0, len(a.settings))
			for _, v := range a.settings {
				buf = append(buf, fmt.Sprintf("%s: %s", v.Name(), v.Get()))
			}
			fmt.Fprint(a.settingsView, strings.Join(buf, "\n"))
		}
		{
			a.switchView.Clear()
			viewWidth, viewHeight := a.switchView.Size()
			x := (viewWidth / 2) - 4
			y := (viewHeight / 2) + 1
			fmt.Fprint(a.switchView, strings.Repeat("\n", y))
			fmt.Fprint(a.switchView, strings.Repeat(" ", x))
			if a.port != nil && a.port.IsRun() {
				a.switchView.BgColor = gocui.ColorGreen
				fmt.Fprint(a.switchView, "运行中")
			} else {
				a.switchView.BgColor = gocui.ColorRed
				fmt.Fprint(a.switchView, "未运行")
			}
		}
		return nil
	})
}
