package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const (
	readBufSizeBytes = 4_096
)

func main() {
	var ipAddr string
	flag.StringVar(&ipAddr, "ipaddr", "", "IP address/port")
	flag.Parse()

	if ipAddr == "" {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("ipAddr", ipAddr)

	conn, err := net.Dial("udp", ipAddr)
	if err != nil {
		log.Fatalf("error dialing: %s", err)
	}
	defer conn.Close()

	resp, err := sendAA55Command(conn, "010200", "0182")
	if err != nil {
		log.Fatalf("error sending command: %s", err)
	}

	modelName := strings.TrimSpace(string(resp[12:22]))
	serialNum := string(resp[38:54])

	log.Printf("modelName = %q, serialNum = %q\n", modelName, serialNum)
}

func sendAA55Command(conn io.ReadWriter, payload, expectedResponseType string) ([]byte, error) {
	bytes, err := hex.DecodeString("AA55C07F" + payload)
	if err != nil {
		return nil, fmt.Errorf("error decoding hex string: %s", err)
	}

	bytes = append(bytes, checksum(bytes)...)

	_, err = fmt.Fprint(conn, string(bytes))
	if err != nil {
		return nil, fmt.Errorf("error writing to socket: %s", err)
	}

	log.Printf("sent data to socket: %X", bytes)

	p := make([]byte, readBufSizeBytes)
	n, err := bufio.NewReader(conn).Read(p)
	if err != nil {
		return nil, fmt.Errorf("error reading from socket: %s", err)
	}

	return validateResponse(p[:n], expectedResponseType)
}

const (
	aa55ResponseLengthIndex  = 6
	aa55ResponseLengthOffset = 9
)

func validateResponse(p []byte, expectedResponseType string) ([]byte, error) {
	if len(p) < 8 {
		return nil, fmt.Errorf("response truncated")
	}

	expectedLen := int(p[aa55ResponseLengthIndex] + aa55ResponseLengthOffset)
	if len(p) != expectedLen {
		return nil, fmt.Errorf("unexpected response length %d (expected %d)", len(p), expectedLen)
	}

	responseType := hex.EncodeToString(p[4:6])
	if responseType != expectedResponseType {
		return nil, fmt.Errorf("unexpected response type `%s` (expected `%s`)", responseType, expectedResponseType)
	}

	var s uint16
	for _, b := range p[:len(p)-2] {
		s += uint16(b)
	}
	expSum := binary.BigEndian.Uint16(p[len(p)-2:])
	if s != expSum {
		return nil, fmt.Errorf("invalid response checksum %d (expected %d)", s, expSum)
	}

	return p, nil
}

func checksum(p []byte) []byte {
	var v uint16
	for _, byte := range p {
		v += uint16(byte)
	}

	c := make([]byte, 4)
	binary.BigEndian.PutUint16(c, v)
	return c
}
