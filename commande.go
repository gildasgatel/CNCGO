package main

import (
	"fmt"
	"image/color"
	"log"
	"regexp"
	"strconv"
	"time"
)

var (
	pannel      int
	speedPlus   bool
	speedMinus  bool
	speedReset  bool
	changeSpeed bool
	STATUS_CNC  string
)

func stateCNC() {
	writeOnPort("?")
	infoState.SetText("")
	fmt.Println(infoMachine)
	if infoMachine != nil {
		for _, v := range infoMachine {
			infoState.SetText(infoState.Text() + " " + v)
		}
		switch infoMachine[0] {
		case "<Idle":
			stateRec.FillColor = color.RGBA{66, 191, 38, 255}
			STATUS_CNC = "Idle"
		case "<Run":
			stateRec.FillColor = color.RGBA{239, 119, 25, 255}
			STATUS_CNC = "Run"
		case "<Hold:0":
			stateRec.FillColor = color.RGBA{184, 2, 2, 255}
			STATUS_CNC = "Hold"
		}
		stateRec.Refresh()
	}
	if len(infoMachine) > 2 {
		rg := regexp.MustCompile(`\d.?.?`)
		s := rg.FindAllString(infoMachine[2], -1)
		if len(s) == 2 {
			//	screenBuff.Text = "Buffer :" + s[0]
			screenPlanner.Text = "Pannel :" + s[1]
			//	screenBuff.Refresh()
			screenPlanner.Refresh()
		}

	}
	infoState.Refresh()
}
func play() {

	for fileScanner.Scan() {
		for STATUS_CNC == "Hold" {
			time.Sleep(time.Second)
			log.Println("Wait....")
			if STATUS_CNC == "Idle" || STATUS_CNC == "Run" {
				stateCNC()
				break
			}
		}
		fmt.Printf("STATUS_CNC: %v\n", STATUS_CNC)
		gcode := fileScanner.Text()
		rg := regexp.MustCompile("F.*")
		s := rg.FindIndex([]byte(gcode))
		if len(s) > 0 {
			s[0] = s[0] + 1
			p, _ := strconv.ParseFloat(gcode[s[0]:s[1]], 32)
			if speedPlus {
				pannel += 100
				speedPlus = false
			}
			if speedMinus {
				pannel -= 100
				speedMinus = false
			}
			if speedReset || !changeSpeed {
				pannel = int(p)
				speedReset = false
			}
			screenSpeed.Text = "Speed : " + strconv.Itoa(pannel)
			screenBuff.Text = " RX :" + strconv.Itoa(totalBuffer) + " lenght RX :" + strconv.Itoa(len(bufferRx))
			screenBuff.Refresh()
			screenSpeed.Refresh()
			gcode = gcode[:s[0]] + strconv.Itoa(pannel)

		}
		writeOnPort(gcode)
		/*	if len(bufferRx) < 4 {
			//stateCNC()
		}*/
	}
	fileGCode.Close()

}
