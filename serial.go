package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hedhyw/Go-Serial-Detector/pkg/v1/serialdet"
	"github.com/tarm/serial"
)

var (
	device      string
	stream      *serial.Port
	infoMachine []string
	bufferRx    []int
	totalBuffer int // Max 128
)

func connection() {
	var err error
	fmt.Println("Connection....")
	config := &serial.Config{
		Name: device,
		Baud: 115200,
		//	ReadTimeout: time.Second,
	}
	fmt.Println(config.Name)
	stream, err = serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("...connecter")
}
func readPort() { // read stream
	index := 0
	readBuff := make([]byte, 128)
	for {
		index++
		n, _ := stream.Read(readBuff)
		etat := string(readBuff[:n])
		log.Println(etat, "index : ", index)

		reg := regexp.MustCompile(`<(^\s)>`) // catch result of "?"
		r := reg.Find([]byte(etat))
		if r != nil {
			infoMachine = strings.Split(string(r), "|")
		}

		regOk := regexp.MustCompile(`ok`) // catch "ok" and ajust RX []
		rOk := regOk.FindAll([]byte(etat), -1)
		if rOk != nil {
			bufferRx = bufferRx[len(rOk):]
			log.Println(totalBuffer)
		} /* else {
			bufferRx = bufferRx[1:]
		}*/
	}
}

func writeOnPort(s string) {

	fmt.Println("Write...")
	totalBuffer = 0              // init RX count
	for _, v := range bufferRx { // Add sum of RX []
		totalBuffer += v
	}
	if totalBuffer+len(s+"\n") < 128 { // Rx MAX 128
		bufferRx = append(bufferRx, len(s+"\n")) // Add sum in bufferRx
		stream.Write([]byte(s + "\n"))
		log.Printf("s: %v\n", s)
		label.SetText(label.Text() + s + "\n")
		label.Refresh()
		scrollText.ScrollToBottom()
	} else {
		time.Sleep(time.Second / 2) // delay before recall
		writeOnPort(s)
	}

}
func startGRBL() {
	log.Println("Initializing grbl...")
	stream.Write([]byte("\r\n\r\n"))
	time.Sleep(time.Second * 2)
	stream.Flush()
	log.Println("Start reading...")
	go readPort()

}

func autoSelecDevice() (dev []string) {
	if list, err := serialdet.List(); err == nil {
		for _, p := range list {
			dev = append(dev, p.Path())
		}
	} else {
		log.Println("no device")
	}
	return
}
