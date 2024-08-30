package RF433T

import (
	"RF433Go/serialDevice"
	"errors"
	"sync"
)

var NSTARTSTOP = 5

type RF433T struct {
	device *serialDevice.Device
	muPort sync.Mutex
}

func New(device *serialDevice.Device) *RF433T {
	return &RF433T{
		device: device,
	}
}

func (r *RF433T) Connect() {
	if r.device.IsPortNil() {
		// for RF433T max 9600
		// for Flipper Zero max 4800 (sample max 4792)
		r.device.SetSpeed(2400)
		r.device.Connect()
	}
}

func (r *RF433T) Disconnect() {
	r.device.Close()
}

func (r *RF433T) Write(data []byte) {
	buffer := make([]byte, 0)
	for i := 0; i < NSTARTSTOP*2; i++ {
		buffer = append(buffer, byte(0xAA))
	}
	for i := 0; i < NSTARTSTOP; i++ {
		buffer = append(buffer, byte(0x02))
	}
	for i := 0; i < len(data); i++ {
		buffer = append(buffer, data[i])
	}
	for i := 0; i < NSTARTSTOP; i++ {
		buffer = append(buffer, byte(0x4))
	}
	r.device.Write(buffer)
}

func (r *RF433T) Read() (int, []byte, error) {
	var start = 0
	var stop = 0
	var buf []byte
	var cnt = 0
	for {
		n, b := r.device.Read()
		if n == 0 {
			return 0, nil, errors.New("No data")
		}

		for i := 0; i < n; i++ {
			symbol := b[i]
			if start >= NSTARTSTOP {
				if symbol == 0x04 {
					stop = stop + 1
					if stop >= NSTARTSTOP {
						start = 0
						stop = 0
						return len(buf), buf, nil
					}
				} else {
					stop = 0
					buf = append(buf, symbol)
					cnt = cnt + 1
				}
			} else {
				if symbol == 0x02 {
					start = start + 1
					stop = 0
					cnt = 0
					buf = make([]byte, 0)
				} else {
					start = 0
				}
			}
		}

	}
}
