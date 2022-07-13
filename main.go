package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"time"

	"git.netflux.io/rob/solar-toolkit/inverter"
)

const commandTimeout = time.Second * 5

func main() {
	var ipAddr string
	flag.StringVar(&ipAddr, "ipaddr", "", "IP address/port")
	flag.Parse()

	if ipAddr == "" {
		flag.Usage()
		os.Exit(1)
	}

	arg := flag.Arg(0)
	if arg != "discover" && arg != "runtime" && arg != "info" {
		log.Fatal("missing command: [discover|runtime|info]")
	}

	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	conn, err := net.Dial("udp", ipAddr)
	if err != nil {
		log.Fatalf("error dialing: %s", err)
	}
	defer conn.Close()

	var (
		inverter inverter.ET
		output   any
	)

	switch arg {
	case "discover":
		log.Fatal("not yet implemented")
	case "info":
		output, err = inverter.DeviceInfo(ctx, conn)
		if err != nil {
			log.Fatalf("error getting device info: %s", err)
		}
	case "runtime":
		output, err = inverter.RuntimeData(ctx, conn)
		if err != nil {
			log.Fatalf("error getting runtime data: %s", err)
		}
	}

	json, err := json.Marshal(output)
	if err != nil {
		log.Fatalf("error encoding JSON: %s", err)
	}

	if _, err = os.Stdout.Write(json); err != nil {
		log.Fatalf("error writing to stdout: %s", err)
	}
}
