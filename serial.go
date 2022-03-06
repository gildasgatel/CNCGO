package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
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
	chanStream  chan string
)

func connection() {

	fmt.Println("Connection....")
	config := &serial.Config{
		Name:        device,
		Baud:        115200,
		ReadTimeout: time.Millisecond,
	}
	fmt.Println(config.Name)
	stream, _ = serial.OpenPort(config)
	log.Println("...connecter")

}

func listenedStream() {
	chanStream = make(chan string)
	fmt.Println("Listend started")
	buff := make([]byte, 128)
	result := ""
	go func() {
		for {
			for {
				n, err := stream.Read(buff)
				result += string(buff[:n])
				if err != nil && len(result) > 0 {
					fmt.Println(result)
					break
				}
			}
			chanStream <- result
			result = ""
		}
	}()

}

func readStream() {
	fmt.Println("Read started")
	rg := regexp.MustCompile(`ok`)
	go func() {
		for {
			result := <-chanStream
			if len(result) > 0 {
				log.Println(result)
			}
			if len(result) > 0 && result[0] == '<' {
				fmt.Println("Etat trouvé")
				infoMachine = strings.Split(string(result), "|")
			}
			r := rg.FindAllString(string(result), -1)
			if r != nil {
				bufferRx = bufferRx[len(r):]
				screenBuff.Text = "buffer : " + strconv.Itoa(totalBuffer)
				screenBuff.Refresh()
			}
		}
	}()
}

func writeOnPort(s string) {

	totalBuffer = 0
	for _, v := range bufferRx {
		totalBuffer += v
	}
	if totalBuffer+len(s+"\n") < 128 {
		bufferRx = append(bufferRx, len(s+"\n"))
		stream.Write([]byte(s + "\n"))
	} else {
		time.Sleep(time.Millisecond)
		writeOnPort(s)
	}

}
func startGRBL() {
	log.Println("Initializing grbl...")
	stream.Write([]byte("\r\n\r\n"))
	time.Sleep(time.Second * 2)
	stream.Flush()
	log.Println("Start reading...")
	listenedStream()
	readStream()
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
