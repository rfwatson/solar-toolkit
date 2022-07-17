package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"git.netflux.io/rob/solar-toolkit/inverter"
)

func main() {
	var inverterAddr string

	flag.StringVar(&inverterAddr, "inverter-addr", "", "IP+port of solar inverter")
	flag.Parse()

	if inverterAddr == "" {
		flag.Usage()
		os.Exit(1)
	}

	conn, err := net.Dial("udp", inverterAddr)
	if err != nil {
		log.Fatalf("error dialing: %s", err)
	}
	defer conn.Close()

	var inv inverter.ET

	runtimeData, err := inv.RuntimeData(context.Background(), conn)
	if err != nil {
		log.Fatalf("error fetching runtime data: %s", err)
	}

	json, err := json.Marshal(runtimeData)
	if err != nil {
		log.Fatalf("error encoding runtime data: %s", err)
	}

	fmt.Fprint(os.Stdout, string(json))
}
