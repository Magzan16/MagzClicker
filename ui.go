package main

func buildUI(hinst uintptr) {
	controls = nil

	regularFont = makeFont(16, 400, "Segoe UI")
	smallFont = makeFont(14, 400, "Segoe UI")
	sectionFont = makeFont(20, 700, "Segoe UI")
	titleFont = makeFont(29, 700, "Segoe UI")
	hotkeyFont = makeFont(20, 700, "Segoe UI")
	buttonFont = makeFont(24, 700, "Segoe UI")

	whiteBrush, _, _ = pCreateSolidBrush.Call(uintptr(rgb(255, 255, 255)))
	lightBlueBrush, _, _ = pCreateSolidBrush.Call(uintptr(rgb(237, 247, 255)))
	blackBrush, _, _ = pCreateSolidBrush.Call(uintptr(rgb(0, 0, 0)))

	logo := createCtl("STATIC", "", SS_ICON|SS_REALSIZECONTROL, 32, 30, 170, 170, idLogo)
	bigIcon, _, _ := pLoadImage.Call(hinst, 1, IMAGE_ICON, 160, 160, 0)
	if bigIcon != 0 {
		pSendMessage.Call(logo, STM_SETIMAGE, IMAGE_ICON, bigIcon)
	}
	magz := label("Magz", 30, 197, 94, 42, idTitleMagz)
	clicker := label("Clicker", 112, 197, 130, 42, idTitleClicker)
	applyFont(magz, titleFont)
	applyFont(clicker, titleFont)
	textColors[clicker] = rgb(210, 37, 37)
	textColors[magz] = rgb(28, 28, 28)

	gInterval := group("Click Interval", 252, 28, 530, 194)
	applyFont(gInterval, sectionFont)

	label("Hours", 278, 75, 90, 26, 0)
	label("Minutes", 398, 75, 90, 26, 0)
	label("Seconds", 518, 75, 90, 26, 0)
	label("Milliseconds", 638, 75, 120, 26, 0)
	numericField("0", 274, 108, 100, 40, idHours, 0, 999)
	numericField("0", 394, 108, 100, 40, idMinutes, 0, 59)
	numericField("0", 514, 108, 100, 40, idSeconds, 0, 59)
	numericField("500", 634, 108, 120, 40, idMillis, 0, 999)
	hint := label("Example: 500 ms = 2 clicks per second", 276, 165, 430, 26, 0)
	applyFont(hint, smallFont)
	textColors[hint] = rgb(110, 110, 110)

	gHot := group("Hotkey", 30, 250, 355, 150)
	applyFont(gHot, sectionFont)
	h1 := label("Global start / stop:", 55, 292, 260, 24, 0)
	applyFont(h1, smallFont)
	textColors[h1] = rgb(95, 95, 95)
	hk := centerLabel("CTRL + ALT + S", 55, 326, 305, 48, idHotkey)
	applyFont(hk, hotkeyFont)
	textColors[hk] = rgb(45, 45, 45)

	gMouse := group("Mouse Action", 410, 250, 372, 150)
	applyFont(gMouse, sectionFont)
	lMouse := label("Mouse", 438, 298, 90, 28, 0)
	lAction := label("Action", 438, 345, 90, 28, 0)
	applyFont(lMouse, regularFont)
	applyFont(lAction, regularFont)
	combo([]string{"Left Button", "Right Button", "Middle Button"}, 530, 292, 225, idButton, 0)
	combo([]string{"Single Click", "Double Click"}, 530, 339, 225, idClickType, 0)

	banner := createCtl("STATIC", "", WS_BORDER, 30, 420, 752, 80, 0)
	bgBrushes[banner] = lightBlueBrush
	modeTitle := label("↖   Mode: Follow current mouse position", 58, 438, 650, 28, idModeTitle)
	modeSub := label("Clicks wherever the mouse cursor is currently located. Orange pulse = click.", 98, 468, 650, 24, idModeSub)
	applyFont(modeTitle, sectionFont)
	applyFont(modeSub, smallFont)
	textColors[modeTitle] = rgb(24, 78, 150)
	textColors[modeSub] = rgb(70, 70, 70)
	bgBrushes[modeTitle] = lightBlueBrush
	bgBrushes[modeSub] = lightBlueBrush

	gCount := group("Click Count", 30, 515, 752, 92)
	applyFont(gCount, sectionFont)
	countLabel := label("Number of clicks", 55, 550, 130, 28, 0)
	applyFont(countLabel, regularFont)
	limitEdit := edit("", 187, 544, 105, 38, idClickLimit)
	applyFont(limitEdit, regularFont)
	estimate := label("Estimated run time: Until stopped", 315, 550, 420, 28, idEstimate)
	applyFont(estimate, regularFont)
	textColors[estimate] = rgb(70, 70, 70)
	limitHint := label("Leave blank to run until you stop it.", 187, 580, 320, 20, 0)
	applyFont(limitHint, smallFont)
	textColors[limitHint] = rgb(110, 110, 110)

	startStop := button("", 30, 625, 230, 58, idStartStop, true)
	applyFont(startStop, buttonFont)

	status := label("●  Status: Stopped", 285, 641, 175, 28, idStatus)
	count := label("Clicks: 0", 485, 641, 120, 28, idPerformed)
	applyFont(status, regularFont)
	applyFont(count, regularFont)
	textColors[status] = rgb(105, 105, 105)
	textColors[count] = rgb(90, 90, 90)

	exitBtn := button("Exit", 662, 630, 120, 48, idExit, false)
	applyFont(exitBtn, regularFont)

	for _, h := range controls {
		applyFont(h, regularFont)
	}
	applyFont(gInterval, sectionFont)
	applyFont(gHot, sectionFont)
	applyFont(gMouse, sectionFont)
	applyFont(gCount, sectionFont)
	applyFont(magz, titleFont)
	applyFont(clicker, titleFont)
	applyFont(hk, hotkeyFont)
	applyFont(modeTitle, sectionFont)
	applyFont(modeSub, smallFont)
	applyFont(h1, smallFont)
	applyFont(hint, smallFont)
	applyFont(limitHint, smallFont)
	applyFont(startStop, buttonFont)

	updateEstimatedRunTime()
	updateUIState()
}
