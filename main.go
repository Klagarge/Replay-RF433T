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
var USE_ROLLING_CODE = false

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
			if USE_ROLLING_CODE {
				rand.Seed(int64(seed))
				seed = rand.Uint64()
			}
			binary.BigEndian.PutUint64(b, seed)
			rf433t.Write(b)
			time.Sleep(1000 * time.Millisecond) // 3ms minimum delay
		}
	}()

	go func() {
		var seed = INITIAL_SEED
		var codes = make([]uint64, SIZE_OF_BUFFER_CODE)
		codes[0] = seed
		for i := 1; i < SIZE_OF_BUFFER_CODE; i++ {
			rand.Seed(int64(seed))
			codes[i] = rand.Uint64()
			seed = codes[i]
		}
		for {
			n, b, err := rf433t.Read()
			if err != nil {
				log.Printf("ERR: %s", err)
			} else {
				var txt = "Read %d bytes: "
				//fmt.Printf("Read %d bytes: ", n)
				txt = fmt.Sprintf(txt, n)
				for i := 0; i < n; i++ {
					//log.Printf("%02X ", b[i])
					txt = fmt.Sprintf("%s%02X ", txt, b[i])
				}
				var code = binary.BigEndian.Uint64(b[:8])
				for i := range codes {
					if code == codes[i] {
						//fmt.Print(" - Code match --> open door !")
						txt = fmt.Sprintf("%s - Code match --> open door !", txt)
						if USE_ROLLING_CODE {
							seed = code
							for j := 0; j < SIZE_OF_BUFFER_CODE; j++ {
								rand.Seed(int64(seed))
								codes[j] = rand.Uint64()
								seed = codes[j]
							}
						}
						break
					}
				}
				//fmt.Println()
				log.Printf(txt)
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
