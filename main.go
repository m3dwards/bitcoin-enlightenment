package main

import (
	// "bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	// "io"
	// "strconv"
	// "strings"
	"time"
	//"encoding/hex"
	"encoding/binary"
)

// Size in bytes of the various parts of a message header
const (
	HeaderMagicSize = 4
	HeaderCommandSize = 12
	HeaderLengthSize = 4
	HeaderChecksumSize = 4
	HeaderSize = HeaderMagicSize + HeaderCommandSize + HeaderLengthSize + HeaderChecksumSize
)

const VERSION = "version"
const VERACK = "verack"

const MIN = 1
const MAX = 100

func random() int {
	return rand.Intn(MAX-MIN) + MIN
}

func main() {
        arguments := os.Args
        if len(arguments) == 1 {
                fmt.Println("Please provide a port number!")
                return
        }

        PORT := ":" + arguments[1]
        l, err := net.Listen("tcp4", PORT)
        if err != nil {
                fmt.Println(err)
                return
        }
        defer l.Close()
        rand.Seed(time.Now().Unix())

        for {
                c, err := l.Accept()
                if err != nil {
                        fmt.Println(err)
                        return
                }
                go handleConnection(c)
        }
}

func handleConnection(c net.Conn) {
        fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	headerbuf := make([]byte, HeaderSize)

	_, err := c.Read(headerbuf)
	if err != nil {
		fmt.Println("Error reading header:", err.Error())
		return
	}
	HeaderLengthStarts := HeaderMagicSize + HeaderCommandSize
        size := headerbuf[HeaderLengthStarts : HeaderLengthStarts + HeaderLengthSize]
	sizeint := binary.LittleEndian.Uint32(size)
	fmt.Println(size)
	fmt.Println(sizeint)
	messagebuff := make([]byte, sizeint)
	// read the full message, or return an error
	_, err = c.Read(messagebuff)
	if err != nil {
		fmt.Println("Error reading body:", err.Error())
		return
	}

	fmt.Printf("received header %x\n", headerbuf)
	fmt.Printf("received %x\n", messagebuff)

	// buf := make([]byte, 1024)
	// // Read the incoming connection into the buffer.
	// _, err := c.Read(buf)
	// if err != nil {
	// 	fmt.Println("Error reading:", err.Error())
	// }
	// fmt.Printf(hex.EncodeToString(buf))
	// Send a response back to person contacting us.
	// c.Write([]byte("Message received."))
	// Close the connection when you're done with it.

	// netData, err := bufio.NewReader(c).ReadString('\n')
	// if err != nil {
	//         fmt.Println(err)
	//         return
	// }

	// temp := strings.TrimSpace(string(netData))
	// if temp == "STOP" {
	//         break
	// }

	// result := strconv.Itoa(random()) + "\n"
	// c.Write([]byte(string(result)))
        c.Close()
}
