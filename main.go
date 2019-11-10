package main

import (
	// "bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	// "io"
	// "strconv"
	"strings"
	"time"
	"encoding/hex"
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

	headerbuff := make([]byte, HeaderSize)

	_, err := c.Read(headerbuff)
	if err != nil {
		fmt.Println("Error reading header:", err.Error())
		return
	}
	HeaderLengthStarts := HeaderMagicSize + HeaderCommandSize
        size := headerbuff[HeaderLengthStarts : HeaderLengthStarts + HeaderLengthSize]
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

	fmt.Printf("received header %x\n", headerbuff)
	fmt.Printf("received %x\n", messagebuff)

	decoded, err := hex.DecodeString(strings.Replace("16 1c 14 12 76 65 72 73 69 6f 6e 00 00 00 00 00 64 00 00 00 35 8d 49 32 62 ea 00 00 01 00 00 00 00 00 00 00 11 b2 d0 50 00 00 00 00 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 00 3b 2e b3 5d 8c e6 17 65 0f 2f 53 61 74 6f 73 68 69 3a 30 2e 37 2e 32 2f c0 3e 03 00", " ", "", -1))
	if err != nil {
		fmt.Println("Error decoding hex:", err.Error())
		return
	}

	c.Write(decoded)

	// time.Sleep(200 * time.Millisecond)

	// _, err = c.Read(headerbuff)
	// if err != nil {
	// 	fmt.Println("Error reading header:", err.Error())
	// 	return
	// }

        // size = headerbuff[HeaderLengthStarts : HeaderLengthStarts + HeaderLengthSize]
	// sizeint = binary.LittleEndian.Uint32(size)

	// _, err = c.Read(messagebuff)
	// if err != nil {
	// 	fmt.Println("Error reading body:", err.Error())
	// 	return
	// }

	// fmt.Printf("received header %x\n", headerbuff)
	// fmt.Printf("received %x\n", messagebuff)

	time.Sleep(200 * time.Millisecond)
	// c.Write([]byte{0x16, 0x1c, 0x14, 0x14, 0x76, 0x65, 0x72, 0x61, 0x63, 0x6B, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5D, 0xF6, 0xE0, 0xE2})

	time.Sleep(2000 * time.Millisecond)

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
