package grbl

import (
	"bufio"
	"bytes"
	"cncgo/backend/internal/api/models"
	"cncgo/backend/pkg/connection"
	"cncgo/backend/pkg/machine"
	"cncgo/backend/pkg/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	RX_BUFFER_SIZE = 128
	MOVE           = "move"
)

type Grbl struct {
	connection     connection.Service
	grblBufferSize []int
	state          models.StateMachine
	process        bool
	priorityGcode  []byte
}

func New(conn connection.Service) (machine.Service, error) {
	if conn == nil {
		return nil, errors.New("error to set machine, connection is nil")
	}
	grbl := &Grbl{connection: conn}
	err := grbl.init()
	if err != nil {
		return nil, err
	}
	return grbl, nil
}

func (grbl *Grbl) init() error {
	err := grbl.connection.Write([]byte("\r\n\r\n"))
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	return grbl.connection.Flush()
}

func (grbl *Grbl) GetName() string {
	return "GRBL"
}

func (grbl *Grbl) SendCommand(data models.Command) ([]byte, error) {
	var gCode []byte
	if data.Command == MOVE {
		gCode = fmt.Appendf([]byte{}, "%s%s\n", data.Axe, data.Distance)
	} else {
		if commandGrbl[data.Command] == "" {
			gCode = fmt.Appendf([]byte{}, "%s\n", data.Command)
		} else {
			gCode = fmt.Appendf([]byte{}, "%s\n", commandGrbl[data.Command])
		}
	}
	if grbl.process {
		grbl.priorityGcode = gCode
		return []byte("CNC running, set priority command"), nil
	}
	err := grbl.connection.Write(gCode)
	if err != nil {
		return nil, err
	}

	return grbl.connection.Read(), nil
}

func (grbl *Grbl) SendFile(path string) error {
	fmt.Println("sendFile: " + path)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	now := time.Now()
	grbl.process = true
	scan := bufio.NewScanner(file)
	lineCount := 0
	for scan.Scan() {
		if err = scan.Err(); err != nil {
			return err
		}
		if len(grbl.priorityGcode) > 0 {
			err = grbl.priorityCommand(grbl.priorityGcode)
			if err != nil {
				return err
			}
			grbl.priorityGcode = []byte{}
		}

		lineCount++
		if lineCount >= 10 {
			err = grbl.priorityCommand([]byte("?\n"))
			if err != nil {
				return err
			}
			lineCount = 0
		}

		line := scan.Bytes()
		line = bytes.TrimSpace(line)
		grbl.grblBufferSize = append(grbl.grblBufferSize, len(line)+1)
		grbl.handleReader()
		line = append(line, 0x0a)
		err = grbl.connection.Write(line)
		if err != nil {
			return err
		}

	}
	grbl.process = false
	fmt.Println("temps de traitement: " + time.Since(now).String())
	return nil
}

func (grbl *Grbl) priorityCommand(command []byte) error {
	grbl.grblBufferSize = append(grbl.grblBufferSize, len(command))
	grbl.handleReader()
	err := grbl.connection.Write(command)
	if err != nil {
		return err
	}
	return nil
}

func (grbl *Grbl) handleReader() {
	countOut := 0

	for utils.Sum(grbl.grblBufferSize) >= RX_BUFFER_SIZE-1 {
		//reset buffer if needed
		if countOut >= 5 {
			for utils.Sum(grbl.grblBufferSize) > grbl.state.BfW {
				grbl.grblBufferSize = grbl.grblBufferSize[1:]
			}
			countOut = 0
		}
		fmt.Println(utils.Sum(grbl.grblBufferSize))
		fmt.Println("Buffer write: ", grbl.state.BfW, " size last line: ", grbl.grblBufferSize[len(grbl.grblBufferSize)-1])
		outGrbl := grbl.connection.Read()
		nbOk := bytes.Count(outGrbl, []byte("ok"))
		log.Printf("outGrbl : %s && nbOk = %d", string(outGrbl), nbOk)
		log.Println(grbl.grblBufferSize)
		if bytes.Contains(outGrbl, []byte("<")) {
			grbl.HandleState(outGrbl)
			//nbOk += bytes.Count(outGrbl, []byte("<"))
		}
		if nbOk > 0 {
			grbl.grblBufferSize = grbl.grblBufferSize[nbOk:]
			log.Println(grbl.grblBufferSize)

		}
		if bytes.Contains(outGrbl, []byte("error")) {
			log.Println(string(outGrbl))
		}

		//check if outGrbl is empty and set counter
		outGrbl = bytes.TrimSpace(outGrbl)
		if len(outGrbl) < 1 {
			countOut++
		} else {
			countOut = 0
		}

	}
}

func (grbl *Grbl) HandleState(line []byte) {

	start := bytes.IndexByte(line, '<')
	end := bytes.IndexByte(line, '>')

	if start != -1 && end != -1 && end > start {
		valeurEntreChevrons := line[start+1 : end]
		datas := bytes.Split(valeurEntreChevrons, []byte("|"))
		for i := 0; i < len(datas); i++ {
			res := bytes.Split(datas[i], []byte(":"))
			if len(res) > 0 {
				datas[i] = res[1]
			}
			res = bytes.Split(datas[i], []byte(","))
			switch i {
			case 0:
				grbl.state.State = string(datas[i])
			case 1:
				grbl.state.MPos = string(datas[i])
			case 2:
				if len(res) == 2 {
					r, _ := strconv.Atoi(string(res[0]))
					w, _ := strconv.Atoi(string(res[1]))
					grbl.state.BfW = w
					grbl.state.BfR = r
				} else {
					grbl.state.BfW = 128
					grbl.state.BfR = 100
				}
			case 3:
				grbl.state.FS = string(datas[i])
			case 4:
				grbl.state.WC0 = string(datas[i])
			}
		}

	}
}

func (grbl *Grbl) GetState() *models.StateMachine {
	return &grbl.state
}
