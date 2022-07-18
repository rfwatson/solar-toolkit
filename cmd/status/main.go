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
	var meterData bool
	var err error

	flag.StringVar(&inverterAddr, "inverter-addr", "", "IP+port of solar inverter")
	flag.BoolVar(&meterData, "meter-data", false, "print meter data, not sensors")
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
	var result []byte

	if meterData {
		meterData, err := inv.MeterData(context.Background(), conn)
		if err != nil {
			log.Fatalf("error fetching meter data: %s", err)
		}

		result, err = json.Marshal(meterData)
		if err != nil {
			log.Fatalf("error encoding meter data: %s", err)
		}
	} else {
		runtimeData, err := inv.RuntimeData(context.Background(), conn)
		if err != nil {
			log.Fatalf("error fetching runtime data: %s", err)
		}

		result, err = json.Marshal(runtimeData)
		if err != nil {
			log.Fatalf("error encoding runtime data: %s", err)
		}
	}

	fmt.Fprint(os.Stdout, string(result))
}
