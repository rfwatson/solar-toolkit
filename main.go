package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"git.netflux.io/rob/goodwe-go/command"
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

	infoCmd, err := command.NewAA55("010200", "0182")
	if err != nil {
		log.Fatalf("error building command: %s", err)
	}

	resp, err := command.Send(infoCmd, conn)
	if err != nil {
		log.Fatalf("error sending command: %s", err)
	}

	modelName := strings.TrimSpace(string(resp[12:22]))
	serialNum := string(resp[38:54])

	log.Printf("modelName = %q, serialNum = %q\n", modelName, serialNum)

	dataCmd := command.NewModbus(command.ModbusCommandTypeRead, 0x891c, 0x007d)
	resp, err = command.Send(dataCmd, conn)
	if err != nil {
		log.Fatalf("error sending command: %s", err)
	}

	log.Printf("rcvd modbus resp = %X", resp)
}
