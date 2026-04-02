package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.String("port", defaultPort(), "port to listen on")
	devicesPath := flag.String("devices", "data.csv", "path to the devices CSV file")
	flag.Parse()

	deviceIDs, err := loadDeviceIDs(*devicesPath)
	if err != nil {
		log.Fatalf("load devices csv: %v", err)
	}
	store := NewStore(deviceIDs)

	handler := NewHandler(store)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/devices/{device_id}/heartbeat", handler.PostHeartbeat)
	mux.HandleFunc("POST /api/v1/devices/{device_id}/stats", handler.PostStats)
	mux.HandleFunc("GET /api/v1/devices/{device_id}/stats", handler.GetStats)

	addr := ":" + *port
	log.Printf("fleet monitor listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func defaultPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}

	return "6733"
}
