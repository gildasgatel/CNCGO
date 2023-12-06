package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	fileScanner *bufio.Scanner
	fileGCode   *os.File
)

func testGcode() {
	var err error
	if strings.HasPrefix(dat, "file://") {
		dat = dat[7:]
	} else {
		fmt.Println("NO DATA")
		return
	}
	fileGCode, err = os.Open(dat)
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	fileScanner = bufio.NewScanner(fileGCode)
	fmt.Println("Open data")

}

func sendGcode(g []string) {

	for _, v := range g {
		writeOnPort(v)
	}
}

func origine() []string {
	c := []string{"G0X0Y0Z0"}
	return c
}

func xPlus() []string {
	c := []string{"G91", "G0X10", "G90"}
	return c
}
func xMinus() []string {
	c := []string{"G91", "G0X-10", "G90"}
	return c
}
func yPlus() []string {
	c := []string{"G91", "G0Y10", "G90"}
	return c
}
func yMinus() []string {
	c := []string{"G91", "G0Y-10", "G90"}
	return c
}
func zPlus() []string {
	c := []string{"G91", "G0Z10", "G90"}
	return c
}
func zMinus() []string {
	c := []string{"G91", "G0Z-10", "G90"}
	return c
}
