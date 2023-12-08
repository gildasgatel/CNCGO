package usb

import (
	"cncgo/api/models"
	"cncgo/pkg/connection"
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

const (
	TIMEOUT        = 500
	RX_BUFFER_SIZE = 128
)

type UsbSender struct {
	Port   *serial.Port
	dataCh chan []byte
}

func New(config *models.Config) (connection.Service, error) {

	port, err := serial.OpenPort(&serial.Config{
		Name: config.PortName,
		Baud: config.BaudRate,
	})
	if err != nil {
		return nil, err
	}
	dataCh := make(chan []byte)

	usb := &UsbSender{
		Port:   port,
		dataCh: dataCh,
	}

	usb.listen()
	return usb, nil
}

func (usb *UsbSender) GetName() string {
	return "USB"
}

func (usb *UsbSender) listen() {
	fmt.Println("Listend started")

	go func() {
		buffer := make([]byte, RX_BUFFER_SIZE)
		for {
			bytesRead, err := usb.Port.Read(buffer)
			if err != nil {
				log.Fatalf("Erreur lors de la lecture depuis le port série: %v", err)
			}
			receivedData := make([]byte, bytesRead) // Création d'un slice pour copier les données lues
			copy(receivedData, buffer[:bytesRead])  // Copie des données lues dans le slice
			usb.dataCh <- receivedData              // Envoie des données lues via le canal
		}
	}()
}

func (usb *UsbSender) Write(data []byte) error {
	fmt.Printf("USB write: %s\n", string(data))
	_, err := usb.Port.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (usb *UsbSender) Read() (resp []byte) {

	timeout := time.After(TIMEOUT * time.Millisecond)
	for {
		select {
		case processedData := <-usb.dataCh:
			//	fmt.Printf(" %s", processedData)
			resp = append(resp, processedData...)
		case <-timeout:
			//fmt.Println("Timeout stop reading.")
			return
		}
	}
}

func (usb *UsbSender) Flush() error {
	return usb.Port.Flush()
}

func (usb *UsbSender) Close() error {
	return usb.Port.Close()
}
