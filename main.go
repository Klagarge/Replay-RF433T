package main

import (
	"RF433Go/RF433T"
	"RF433Go/serialDevice"
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"
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
		var i uint64 = 250
		b := make([]byte, 8)
		for {
			binary.BigEndian.PutUint64(b, i)
			rf433t.Write(b)
			i++
			time.Sleep(500 * time.Millisecond) // 3ms minimum delay
		}
	}()

	go func() {
		for {
			n, b, err := rf433t.Read()
			if err != nil {
				log.Printf("ERR: %s", err)
			} else {
				fmt.Printf("Read %d bytes: ", n)
				for i := 0; i < n; i++ {
					fmt.Printf("%02X ", b[i])
				}
				fmt.Println()
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
