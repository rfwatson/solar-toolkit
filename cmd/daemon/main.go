package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"git.netflux.io/rob/solar-toolkit/inverter"
)

const (
	httpUserAgent = "solar-toolkit (git.netflux.io)"
	httpTimeout   = time.Second * 5
)

func main() {
	var (
		inverterAddr    string
		gatewayEndpoint string
		gatewayUsername string
		gatewayPassword string
		pollInterval    time.Duration
	)

	flag.StringVar(&inverterAddr, "inverter-addr", "", "IP+port of solar inverter")
	flag.StringVar(&gatewayEndpoint, "endpoint", "", "URL to post metrics to")
	flag.StringVar(&gatewayUsername, "username", "", "HTTP basic auth username")
	flag.StringVar(&gatewayPassword, "password", "", "HTTP basic auth password")
	flag.DurationVar(&pollInterval, "pollInterval", time.Minute, "Poll interval, example: 60s")
	flag.Parse()

	if gatewayEndpoint == "" {
		flag.Usage()
		os.Exit(1)
	}

	conn, err := net.Dial("udp", inverterAddr)
	if err != nil {
		log.Fatalf("error dialing: %s", err)
	}
	defer conn.Close()

	var inv inverter.ET

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	client := http.Client{Timeout: httpTimeout}

	for ; true; <-ticker.C {
		runtimeData, err := inv.RuntimeData(context.Background(), conn)
		if err != nil {
			log.Printf("error fetching runtime data: %s", err)
			continue
		}

		reqBody, err := json.Marshal(runtimeData)
		if err != nil {
			log.Printf("error encoding runtime data: %s", err)
			continue
		}

		req, err := http.NewRequest(http.MethodPost, gatewayEndpoint, bytes.NewReader(reqBody))
		if err != nil {
			log.Printf("error building request: %s", err)
			continue
		}

		req.Header.Set("content-type", "application/json")
		req.Header.Set("user-agent", httpUserAgent)
		if gatewayUsername != "" && gatewayPassword != "" {
			req.SetBasicAuth(gatewayUsername, gatewayPassword)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("error sending request: %s", err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("unexpected HTTP response code: %d", resp.StatusCode)
			continue
		}
		log.Printf("OK: %s", runtimeData.PVPower.String())
	}
}
