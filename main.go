package main

import (
	// "bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	// "io"
	// "strconv"
	"encoding/binary"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
	"crypto/sha256"

)

// Size in bytes of the various parts of a message header
const (
	HeaderMagicSize    = 4
	HeaderCommandSize  = 12
	HeaderLengthSize   = 4
	HeaderChecksumSize = 4
	HeaderSize         = HeaderMagicSize + HeaderCommandSize + HeaderLengthSize + HeaderChecksumSize
)

const VERSION = "version"
const VERACK = "verack"

const MIN = 1
const MAX = 100

func random() int {
	return rand.Intn(MAX-MIN) + MIN
}

func main() {
	if "development" != "development" {
		log.SetFormatter(&log.JSONFormatter{})
	}

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
	log.WithFields(log.Fields{
		"port": PORT,
	}).Info("Server started")
	first256 := (sha256.New()).Sum(nil)
	secondhash := sha256.New()
	secondhash.Write(first256)
	log.Info(hex.EncodeToString(secondhash.Sum(nil)))

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}

func createVersion() []byte {
	decoded, err := hex.DecodeString(strings.Replace("16 1c 14 12 76 65 72 73 69 6f 6e 00 00 00 00 00 64 00 00 00 35 8d 49 32 62 ea 00 00 01 00 00 00 00 00 00 00 11 b2 d0 50 00 00 00 00 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 00 3b 2e b3 5d 8c e6 17 65 0f 2f 53 61 74 6f 73 68 69 3a 30 2e 37 2e 32 2f c0 3e 03 00", " ", "", -1))
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error()}).Error("Error decoding hex", err.Error())
		return nil
	}
	return decoded
}

func readMessageAndPrintIt(c net.Conn) {
	headerbuff := make([]byte, HeaderSize)

	_, err := c.Read(headerbuff)
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error()}).Error("Error reading header")
		return
	}
	HeaderLengthStarts := HeaderMagicSize + HeaderCommandSize
	size := headerbuff[HeaderLengthStarts : HeaderLengthStarts+HeaderLengthSize]
	sizeint := binary.LittleEndian.Uint32(size)
	messagebuff := make([]byte, sizeint)

	// read the full message, or return an error
	_, err = c.Read(messagebuff)
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error()}).Error("Error reading body:", err.Error())
		return
	}

	log.WithFields(log.Fields{"header": hex.EncodeToString(headerbuff), "body": hex.EncodeToString(messagebuff)}).Debug("Message Received")

}

func createVerack() []byte {

	networkMagic := []byte{ 0x16, 0x1c, 0x14, 0x12 }
	verackType := []byte{ 0x76, 0x65, 0x72, 0x61, 0x63, 0x6B, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00 }
	payloadSize := []byte{ 0x00, 0x00, 0x00, 0x00 }
	checkSum := []byte{ 0x5D, 0xF6, 0xE0, 0xE2 }

	message := append(networkMagic,append(verackType,append(payloadSize, checkSum...)...)...)

	return message
}

func handleConnection(c net.Conn) {
	log.WithFields(log.Fields{"remote_address": c.RemoteAddr().String()}).Debug("Serving inbound connection")

	readMessageAndPrintIt(c)
	c.Write(createVersion())

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

	c.Write(createVerack())

	readMessageAndPrintIt(c)

	log.Info("I have correctly read their message")


	time.Sleep(2000 * time.Millisecond)
	// c.Write([]byte{0x16, 0x1c, 0x14, 0x14, 0x76, 0x65, 0x72, 0x61, 0x63, 0x6B, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5D, 0xF6, 0xE0, 0xE2})

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
