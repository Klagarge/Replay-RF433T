package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
	"wmbusGo/RF433T"
	"wmbusGo/serialDevice"
)

var NSTARTSTOP = 20

func main() {
	namePort, err := serialDevice.Search()
	if err != nil {
		log.Fatal(err)
	}

	var d = serialDevice.New(namePort)
	var rf433t = RF433T.New(d)
	rf433t.Connect()

	time.Sleep(1 * time.Second)

	go func() {
		for {
			d.Write([]byte{0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
				0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
				0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
				0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
				0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
				0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
				0x48, 0x45, 0x4C, 0x4C, 0x4F,
				0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04,
				0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04})
			//fmt.Println("Sent")
			time.Sleep(1000 * time.Millisecond) // 3ms minimum delay
		}
	}()

	go func() {
		var start = 0
		var stop = 0
		var trame = false
		var buf []byte
		var cnt = 0
		for {
			n, b := d.Read()
			/*
				fmt.Printf("Received %d bytes: ", n)
				for i := 0; i < n; i++ {
					fmt.Printf("%02x ", b[i])
				}
				fmt.Println("")
			*/

			for i := 0; i < n; i++ {
				symbol := b[i]
				if start >= NSTARTSTOP {
					if symbol == 0x04 {
						stop = stop + 1
						if stop >= NSTARTSTOP {
							trame = true
						}
					} else {
						buf = append(buf, symbol)
						cnt = cnt + 1
					}
				}
				if symbol == 0x02 {
					start = start + 1
					stop = 0
					cnt = 0
					buf = make([]byte, 0)
				}
			}
			if trame {
				fmt.Printf("Trame: ")
				for i := 0; i < len(buf); i++ {
					fmt.Printf("%02x ", buf[i])
				}
				fmt.Println("")
				trame = false
				start = 0
				stop = 0
				cnt = 0
				buf = make([]byte, 0)
			}

		}
	}()

	go leaveOnQ()

	// Keep the main goroutine running
	select {}
}

func leaveOnQ() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press 'q' to quit")
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("ERR: %s", err)
		}
		if input == "q\n" {
			fmt.Println("Exiting program...")
			os.Exit(0)
		}
	}
}
