package main

import (
	"RF433Go/RF433T"
	"RF433Go/serialDevice"
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var INITIAL_SEED uint64 = 250
var SIZE_OF_BUFFER_CODE = 256

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
		var seed = INITIAL_SEED
		b := make([]byte, 8)
		for {
			rand.Seed(int64(seed))
			seed = rand.Uint64()
			binary.BigEndian.PutUint64(b, seed)
			rf433t.Write(b)
			time.Sleep(1000 * time.Millisecond) // 3ms minimum delay
		}
	}()

	go func() {
		var seed = INITIAL_SEED
		for {
			n, b, err := rf433t.Read()
			if err != nil {
				log.Printf("ERR: %s", err)
			} else {
				fmt.Printf("Read %d bytes: ", n)
				for i := 0; i < n; i++ {
					fmt.Printf("%02X ", b[i])
				}
				for i := 0; i < SIZE_OF_BUFFER_CODE; i++ {
					rand.Seed(int64(seed))
					expected := rand.Uint64()
					received := binary.BigEndian.Uint64(b[:8])
					if expected == received {
						fmt.Print(" - Codes match --> open door !")
						break
					}
					seed = expected
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
