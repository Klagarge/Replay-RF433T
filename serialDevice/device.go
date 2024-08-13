package serialDevice

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

type Device struct {
	portName string
	mode     *serial.Mode
	port     *serial.Port
	mu       sync.Mutex
}

func Search() (string, error) {
	portList, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return "", errors.New("Can't get ports list")
	}

	if len(portList) == 0 {
		return "", errors.New("No serial ports found!")
	}
	var portName string
	for _, port := range portList {
		if port.IsUSB {
			portName = port.Name
			break
		}
	}

	if portName == "" {
		return "", errors.New("No USB serial port found!")
	}

	return portName, nil
}

func New(portName string) *Device {
	return &Device{
		portName: portName,
	}
}

func (d *Device) SetSpeed(speed int) {
	d.mode = &serial.Mode{
		BaudRate: speed,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
		DataBits: 8,
	}
}

func (d *Device) Connect() {
	p, err := serial.Open(d.portName, d.mode)
	if err != nil {
		log.Fatalf("ERR on opening: %s", err)
	}
	d.port = &p
	fmt.Printf("Connected to %s\n", d.portName)

	err = p.SetReadTimeout(serial.NoTimeout)
	if err != nil {
		log.Fatalf("ERR on setting read timeout: %s", err)
	}
	d.port = &p
}

func (d *Device) Close() {
	d.mu.Lock()
	defer d.mu.Unlock()
	p := *d.port
	err := p.Close()
	if err != nil {
		log.Fatalf("ERR on closing: %s", err)
	}
	fmt.Println("Port closed")
}

func (d *Device) Write(b []byte) {
	d.mu.Lock()
	defer d.mu.Unlock()
	p := *d.port
	_, err := p.Write(b)
	if err != nil {
		log.Fatalf("ERR on writing: %s", err)
	}
}

func (d *Device) Read() (int, []byte) {
	d.mu.Lock()
	d.mu.Unlock()
	p := *d.port
	buffer := make([]byte, 128)
	n, err := p.Read(buffer)
	if err != nil {
		log.Fatalf("ERR on reading: %s", err)
	}
	if n == 0 {
		return 0, nil
	}
	return n, buffer
}

func (d *Device) IsPortNil() bool {
	if d.port == nil {
		return true
	}
	return false
}
