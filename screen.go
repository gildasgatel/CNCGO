package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	label         *widget.TextGrid
	proBar        *widget.ProgressBar
	scrollText    *container.Scroll
	stateRec      *canvas.Rectangle
	infoState     *widget.TextGrid
	screenBuff    *canvas.Text
	screenPlanner *canvas.Text
	screenSpeed   *canvas.Text
	w             fyne.Window
	a             fyne.App
	dat           string
)

func display() {

	a = app.New()
	w = a.NewWindow("CNCGo")
	w.Resize(fyne.NewSize(400, 400))

	connectionScreen()

}

func connectionScreen() {

	slc := widget.NewSelectEntry(autoSelecDevice())
	if autoSelecDevice() != nil {
		slc.Text = autoSelecDevice()[0]
	}

	btn := widget.NewButton("go", func() {
		device = slc.Text
		connection()
		mainScreen()
	})
	btnUpdate := widget.NewButton("scan", func() {
		slc.SetOptions(autoSelecDevice())
		slc.Text = autoSelecDevice()[0]
		slc.Refresh()

	})
	contCenter := container.NewVBox(layout.NewSpacer(), slc, btn, btnUpdate, layout.NewSpacer())
	w.SetContent(contCenter)
	w.ShowAndRun()

}

func mainScreen() {

	input := widget.NewEntry()
	winDialog := a.NewWindow("Open File")
	winDialog.Resize(fyne.NewSize(400, 200))

	proBar = widget.NewProgressBar()
	proBar.Min = 0
	proBar.Max = 128
	label = widget.NewTextGrid()
	label.SetText("Connection...\n")
	label.SetText(label.Text() + "CNCGo 1.0  Port: " + device + "\n")
	infoState = widget.NewTextGrid()
	scrollText = container.NewScroll(label)
	stateRec = canvas.NewRectangle(color.Black)
	stateRec.SetMinSize(fyne.NewSize(50, 50))
	stateRec.Show()
	screenPlanner = canvas.NewText("Pannel :", color.White)
	screenBuff = canvas.NewText("Buffer :", color.White)
	screenSpeed = canvas.NewText("Speed :", color.White)

	btn := widget.NewButton("go", func() {
		writeOnPort(input.Text)
	})
	btnSetting := widget.NewButton("Settings", func() {
		writeOnPort("$$")
	})

	btnX1 := widget.NewButton("X+", func() {
		sendGcode(xPlus())
	})
	btnX2 := widget.NewButton("X-", func() {
		sendGcode(xMinus())
	})
	btnY1 := widget.NewButton("Y+", func() {
		sendGcode(yPlus())
	})
	btnY2 := widget.NewButton("Y-", func() {
		sendGcode(yMinus())
	})
	btnZ1 := widget.NewButton("Z+", func() {
		sendGcode(zPlus())
	})
	btnZ2 := widget.NewButton("Z-", func() {
		sendGcode(zMinus())
	})
	btnP2 := widget.NewButton("Start", func() {
		go play()
		time.Sleep(time.Second * 2)
		//stateCNC()
	})
	btnP1 := widget.NewButton("Set X0 Y0 Z0", func() {

		writeOnPort("G10 P1 L20 X0 Y0 Z0")
	})
	btnFile := widget.NewButton("Open File", func() {
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, e error) {
			if e != nil {
				dialog.ShowError(e, w)
				return
			}
			if uc.URI().Extension() == ".nc" {
				dat = uc.URI().String()
				testGcode()
			}
		}, w)
	})
	btnEtat := widget.NewButton("etat", func() {
		stateCNC()
	})
	btnFeedP := widget.NewButton("Speed +10", func() {
		speedPlus = true
		changeSpeed = true
	})
	btnFeedM := widget.NewButton("Speed -10", func() {
		speedMinus = true
		changeSpeed = true
	})
	btnFeedR := widget.NewButton("Speed Reset", func() {
		speedReset = true
		changeSpeed = true
	})
	btnPlay := widget.NewButton("Reprendre", func() {
		writeOnPort("~")
		stateCNC()
	})
	btnPause := widget.NewButton("Pause", func() {
		writeOnPort("!")
		stateCNC()
	})

	btnOrigine := widget.NewButton("Go X0 Y0 Z0", func() {
		sendGcode(origine())
	})
	startGRBL()
	scrollText.SetMinSize(fyne.NewSize(200, 200))
	commande := container.NewHBox(btnX1, btnX2, btnY1, btnY2, btnZ1, btnZ2, btnP1, btnOrigine, btnP2, btnFile)
	infoS := container.NewHBox(infoState)
	cont := container.NewVBox(stateRec, input, btn, btnSetting, scrollText, proBar, infoS, commande)
	barreLateral := container.NewVBox(btnEtat, btnFeedP, btnFeedM, btnFeedR, screenBuff, screenPlanner, screenSpeed, btnPlay, btnPause)
	contMain := container.NewHBox(cont, barreLateral)

	w.SetContent(contMain)

}
